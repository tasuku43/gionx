package cli

import (
	"testing"
)

func TestFindAttachDetachTrigger_ControlByte(t *testing.T) {
	i, n := findAttachDetachTrigger([]byte("abc\x1dxyz"))
	if i != 3 || n != 1 {
		t.Fatalf("idx=%d len=%d, want idx=3 len=1", i, n)
	}
}

func TestFindAttachDetachTrigger_CSIU(t *testing.T) {
	i, n := findAttachDetachTrigger([]byte("ab\x1b[93;5uzz"))
	if i != 2 || n != len("\x1b[93;5u") {
		t.Fatalf("idx=%d len=%d, want idx=2 len=%d", i, n, len("\x1b[93;5u"))
	}
}

func TestTrailingAttachControlPrefixLength(t *testing.T) {
	got := trailingAttachControlPrefixLength([]byte("ab\x1b[93;"))
	if got != len("\x1b[93;") {
		t.Fatalf("hold=%d, want=%d", got, len("\x1b[93;"))
	}
}

func TestTrailingAttachControlPrefixLength_WheelPrefix(t *testing.T) {
	got := trailingAttachControlPrefixLength([]byte("ab\x1b[<64;10;"))
	if got != len("\x1b[<64;10;") {
		t.Fatalf("hold=%d, want=%d", got, len("\x1b[<64;10;"))
	}
}

func TestTrailingAttachControlPrefixLength_ZeroForNoPrefix(t *testing.T) {
	got := trailingAttachControlPrefixLength([]byte("hello"))
	if got != 0 {
		t.Fatalf("hold=%d, want=0", got)
	}
}

func TestTranslateAttachWheelToPaging_SGR(t *testing.T) {
	got := string(translateAttachWheelToPaging([]byte("\x1b[<64;10;20M\x1b[<65;10;21M")))
	if got != "\x1b[5~\x1b[6~" {
		t.Fatalf("translated=%q", got)
	}
}

func TestTranslateAttachWheelToPaging_URXVT(t *testing.T) {
	got := string(translateAttachWheelToPaging([]byte("\x1b[64;10;20M\x1b[65;10;21M")))
	if got != "\x1b[5~\x1b[6~" {
		t.Fatalf("translated=%q", got)
	}
}

func TestTranslateAttachWheelToPaging_KeepOtherInput(t *testing.T) {
	got := string(translateAttachWheelToPaging([]byte("abc\x1b[Adef")))
	if got != "abc\x1b[Adef" {
		t.Fatalf("translated=%q", got)
	}
}
