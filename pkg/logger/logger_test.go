package logger

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w
	os.Stderr = w

	outCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outCh <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return <-outCh
}

func TestInitLogger_DevModeEnablesDebug(t *testing.T) {
	prevLogger := slog.Default()
	defer slog.SetDefault(prevLogger)

	output := captureOutput(t, func() {
		InitLogger("dev")
		slog.Debug("debug message", "key", "value")
	})

	require.Contains(t, output, "level=DEBUG")
	require.Contains(t, output, `msg="debug message"`)
	require.Contains(t, output, "key=value")
}

func TestInitLogger_ProdModeSuppressesDebug(t *testing.T) {
	prevLogger := slog.Default()
	defer slog.SetDefault(prevLogger)

	output := captureOutput(t, func() {
		InitLogger("prod")
		slog.Debug("hidden debug")
		slog.Info("visible info")
	})

	require.NotContains(t, output, "hidden debug")
	require.Contains(t, output, "visible info")
}

func TestInitLogger_UnknownModeFallsBackToInfo(t *testing.T) {
	prevLogger := slog.Default()
	defer slog.SetDefault(prevLogger)

	output := captureOutput(t, func() {
		tempHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(tempHandler))

		InitLogger("staging")
		slog.Debug("should not appear")
		slog.Info("info message")
	})

	require.Contains(t, output, "unknown mode")
	require.Contains(t, output, `msg="info message"`)
	require.NotContains(t, output, "should not appear")
}
