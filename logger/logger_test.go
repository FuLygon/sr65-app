package logger

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	log.Out = &buf
	err := errors.New("error message")

	Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Errorf("info was incorrect, got: %s, want: %s.", buf.String(), "info message")
	}

	buf.Reset()
	Warn("warn message", err)
	if !strings.Contains(buf.String(), "warn message") {
		t.Errorf("warn was incorrect, got: %s, want: %s.", buf.String(), "warn message")
	}
	if !strings.Contains(buf.String(), fmt.Sprintf(`error="%s"`, err)) {
		t.Errorf("warn was incorrect, got: %s, missing error field.", buf.String())
	}

	buf.Reset()
	Error("error message", err)
	if !strings.Contains(buf.String(), "error message") {
		t.Errorf("error was incorrect, got: %s, want: %s.", buf.String(), "error message")
	}
	if !strings.Contains(buf.String(), fmt.Sprintf(`error="%s"`, err)) {
		t.Errorf("error was incorrect, got: %s, missing error field.", buf.String())
	}
}
