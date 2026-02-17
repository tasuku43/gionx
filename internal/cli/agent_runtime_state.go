//go:build experimental

package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/tasuku43/kra/internal/infra/paths"
)

type agentRuntimeSessionRecord struct {
	SessionID      string `json:"session_id"`
	RootPath       string `json:"root_path"`
	WorkspaceID    string `json:"workspace_id"`
	ExecutionScope string `json:"execution_scope"`
	RepoKey        string `json:"repo_key"`
	Kind           string `json:"kind"`
	PID            int    `json:"pid"`
	StartedAt      int64  `json:"started_at"`
	UpdatedAt      int64  `json:"updated_at"`
	Seq            int64  `json:"seq"`
	RuntimeState   string `json:"runtime_state"`
	ExitCode       *int   `json:"exit_code"`
	storagePath    string `json:"-"`
}

const (
	agentRuntimeExitedKeepPerWorkspace = 3
	agentRuntimeExitedRetention        = 24 * time.Hour
)

func newAgentRuntimeSessionID(now time.Time, pid int) string {
	return fmt.Sprintf("s-%d-%d", now.UnixNano(), pid)
}

func saveAgentRuntimeSession(record agentRuntimeSessionRecord) error {
	kraHome, err := paths.KraHomeDir()
	if err != nil {
		return fmt.Errorf("resolve KRA_HOME: %w", err)
	}
	rootHash := hashRootPath(record.RootPath)
	dir := filepath.Join(kraHome, "state", "agents", rootHash)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create runtime state dir: %w", err)
	}

	path := filepath.Join(dir, record.SessionID+".json")
	b, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal runtime session: %w", err)
	}
	if err := writeFileAtomically(path, append(b, '\n'), 0o644); err != nil {
		return fmt.Errorf("write runtime session: %w", err)
	}
	return nil
}

func loadAgentRuntimeSessions(root string) ([]agentRuntimeSessionRecord, error) {
	kraHome, err := paths.KraHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve KRA_HOME: %w", err)
	}
	rootHash := hashRootPath(root)
	dir := filepath.Join(kraHome, "state", "agents", rootHash)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []agentRuntimeSessionRecord{}, nil
		}
		return nil, fmt.Errorf("read runtime state dir: %w", err)
	}
	rows := make([]agentRuntimeSessionRecord, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var r agentRuntimeSessionRecord
		if err := json.Unmarshal(b, &r); err != nil {
			continue
		}
		r.SessionID = strings.TrimSpace(r.SessionID)
		r.WorkspaceID = strings.TrimSpace(r.WorkspaceID)
		r.ExecutionScope = strings.TrimSpace(strings.ToLower(r.ExecutionScope))
		r.RepoKey = strings.TrimSpace(r.RepoKey)
		r.Kind = strings.TrimSpace(r.Kind)
		r.RuntimeState = strings.TrimSpace(strings.ToLower(r.RuntimeState))
		if r.RuntimeState == "" {
			r.RuntimeState = "unknown"
		}
		r.storagePath = path
		rows = append(rows, r)
	}
	rows, stalePaths := pruneExitedRuntimeSessions(rows, time.Now().Unix())
	for _, path := range stalePaths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			// Best-effort GC: keep serving records even if one file cannot be removed.
			continue
		}
	}
	slices.SortFunc(rows, func(a, b agentRuntimeSessionRecord) int {
		if a.UpdatedAt != b.UpdatedAt {
			if a.UpdatedAt > b.UpdatedAt {
				return -1
			}
			return 1
		}
		if cmp := strings.Compare(a.WorkspaceID, b.WorkspaceID); cmp != 0 {
			return cmp
		}
		return strings.Compare(a.SessionID, b.SessionID)
	})
	return rows, nil
}

func pruneExitedRuntimeSessions(rows []agentRuntimeSessionRecord, nowUnix int64) ([]agentRuntimeSessionRecord, []string) {
	exitedByWorkspace := map[string][]int{}
	for i := range rows {
		if rows[i].RuntimeState != "exited" {
			continue
		}
		exitedByWorkspace[rows[i].WorkspaceID] = append(exitedByWorkspace[rows[i].WorkspaceID], i)
	}

	keepExited := map[int]bool{}
	stalePaths := make([]string, 0)
	for _, idxs := range exitedByWorkspace {
		slices.SortFunc(idxs, func(a, b int) int {
			if rows[a].UpdatedAt != rows[b].UpdatedAt {
				if rows[a].UpdatedAt > rows[b].UpdatedAt {
					return -1
				}
				return 1
			}
			return strings.Compare(rows[a].SessionID, rows[b].SessionID)
		})
		for rank, idx := range idxs {
			isExpired := false
			if rows[idx].UpdatedAt > 0 {
				age := nowUnix - rows[idx].UpdatedAt
				isExpired = age > int64(agentRuntimeExitedRetention/time.Second)
			}
			if rank < agentRuntimeExitedKeepPerWorkspace && !isExpired {
				keepExited[idx] = true
				continue
			}
			stalePaths = append(stalePaths, rows[idx].storagePath)
		}
	}

	trimmed := make([]agentRuntimeSessionRecord, 0, len(rows))
	for i := range rows {
		r := rows[i]
		if r.RuntimeState == "exited" && !keepExited[i] {
			continue
		}
		trimmed = append(trimmed, r)
	}
	return trimmed, stalePaths
}

func printAgentRuntimeListTSV(out io.Writer, rows []agentRuntimeSessionRecord) {
	fmt.Fprintln(out, "session_id\tworkspace_id\texecution_scope\trepo_key\tkind\truntime_state\tstarted_at\tupdated_at\tpid\texit_code")
	for _, r := range rows {
		exitCode := "-"
		if r.ExitCode != nil {
			exitCode = strconv.Itoa(*r.ExitCode)
		}
		fmt.Fprintf(
			out,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			r.SessionID,
			r.WorkspaceID,
			r.ExecutionScope,
			r.RepoKey,
			r.Kind,
			r.RuntimeState,
			formatUnixTS(r.StartedAt),
			formatUnixTS(r.UpdatedAt),
			r.PID,
			exitCode,
		)
	}
}

func printAgentRuntimeListHuman(out io.Writer, rows []agentRuntimeSessionRecord, useColor bool) {
	body := make([]string, 0, len(rows))
	if len(rows) == 0 {
		body = append(body, fmt.Sprintf("%s(none)", uiIndent))
		printSection(out, "Agents:", body, sectionRenderOptions{
			blankAfterHeading: true,
			trailingBlank:     true,
		})
		return
	}
	maxCols := listTerminalWidth()
	for _, r := range rows {
		location := "workspace"
		if r.ExecutionScope == "repo" {
			location = "repo:" + r.RepoKey
		}
		line := fmt.Sprintf(
			"%s• %s  session:%s  location:%s  kind:%s  state:%s  updated:%s",
			uiIndent,
			r.WorkspaceID,
			r.SessionID,
			location,
			r.Kind,
			r.RuntimeState,
			formatUnixTS(r.UpdatedAt),
		)
		if useColor {
			line = styleMuted(line, useColor)
		}
		body = append(body, truncateDisplay(line, maxCols))
	}
	printSection(out, "Agents:", body, sectionRenderOptions{
		blankAfterHeading: true,
		trailingBlank:     true,
	})
}

func hashRootPath(root string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(root)))
	return hex.EncodeToString(sum[:8])
}

func writeFileAtomically(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}
	if _, err := tmp.Write(content); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}
