package cli

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hinshun/vt10x"
)

func normalizeAgentAttachRenderer(v string) string {
	normalized := strings.ToLower(strings.TrimSpace(v))
	switch normalized {
	case "", attachRendererAuto:
		return attachRendererAuto
	case attachRendererRaw:
		return attachRendererRaw
	case attachRendererVT10X:
		return attachRendererVT10X
	default:
		return normalized
	}
}

func proxyAgentAttachIOWithRenderer(
	root string,
	sessionID string,
	conn *net.UnixConn,
	in io.Reader,
	out io.Writer,
	mode agentAttachMode,
	renderer string,
) error {
	switch normalizeAgentAttachRenderer(renderer) {
	case attachRendererVT10X:
		return proxyAgentAttachIOWithVT10X(root, sessionID, conn, in, out, mode)
	case attachRendererAuto:
		return proxyAgentAttachIOWithVT10X(root, sessionID, conn, in, out, mode)
	default:
		return fmt.Errorf("unknown attach renderer: %s", renderer)
	}
}

func proxyAgentAttachIOWithVT10X(
	root string,
	sessionID string,
	conn *net.UnixConn,
	in io.Reader,
	out io.Writer,
	mode agentAttachMode,
) error {
	if conn == nil {
		return fmt.Errorf("broker connection is nil")
	}
	if mode.fullscreen && isTerminalWriter(out) {
		writeAttachTerminalEnter(out)
	} else if mode.clearOnEnter && isTerminalWriter(out) {
		writeAttachTerminalClear(out)
	}
	if !mode.fullscreen && mode.writeBoundary && isTerminalWriter(out) {
		writeAttachSessionBoundary(out)
	}
	if mode.flushInput && isTerminalReader(in) {
		flushTerminalInputBuffer(in)
	}

	restore, err := maybeEnterRawMode(in, out)
	if err != nil {
		return err
	}
	if mode.flushInput && isTerminalReader(in) {
		defer flushTerminalInputBuffer(in)
	}
	if restore != nil {
		defer restore()
	}
	if mode.restoreShell && isTerminalWriter(out) {
		defer writeAttachTerminalRestore(out)
	}

	cols, rows := terminalSize(in, out)
	if cols <= 0 {
		cols = 80
	}
	if rows <= 0 {
		rows = 24
	}
	term := vt10x.New(vt10x.WithSize(cols, rows))

	stopResizeWatcher := startAttachResizeWatcher(root, sessionID, in, out)
	defer stopResizeWatcher()

	readErrCh := make(chan error, 1)
	chunkCh := make(chan []byte, 32)
	go func() {
		buf := make([]byte, 8192)
		for {
			n, rerr := conn.Read(buf)
			if n > 0 {
				payload := append([]byte(nil), buf[:n]...)
				chunkCh <- payload
			}
			if rerr != nil {
				close(chunkCh)
				readErrCh <- rerr
				return
			}
		}
	}()

	inputResCh := make(chan attachInputResult, 1)
	go func() {
		inputResCh <- forwardAttachInput(conn, in, mode.localDetach)
	}()

	var sigintCh chan os.Signal
	if mode.localDetach {
		sigintCh = make(chan os.Signal, 1)
		signal.Notify(sigintCh, os.Interrupt, syscall.SIGINT)
		defer signal.Stop(sigintCh)
	}

	for {
		select {
		case <-sigintCh:
			_ = conn.Close()
			return errAgentAttachDetached
		case inputRes := <-inputResCh:
			if inputRes.detached {
				_ = conn.Close()
				return errAgentAttachDetached
			}
			if isAgentAttachIOError(inputRes.err) {
				_ = conn.Close()
				return inputRes.err
			}
		case chunk, ok := <-chunkCh:
			if !ok {
				readErr := <-readErrCh
				if isAgentAttachIOError(readErr) {
					return readErr
				}
				return nil
			}
			if len(chunk) == 0 {
				continue
			}
			curCols, curRows := terminalSize(in, out)
			if curCols > 0 && curRows > 0 && (curCols != cols || curRows != rows) {
				cols = curCols
				rows = curRows
				term.Resize(cols, rows)
			}
			_, _ = term.Write(chunk)
			_, _ = io.WriteString(out, "\r\x1b[2J\x1b[H")
			_, _ = io.WriteString(out, term.String())
		}
	}
}
