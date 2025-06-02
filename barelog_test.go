// barelog_test.go
package barelog

import (
	"bytes"
	"context"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// helper to capture stdout during a test, and restore it afterward
func captureStdout(f func()) string {
	// Save the existing stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the code that writes to stdout
	f()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = old

	// Read everything that was written
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

// Test that a logger formats a single WARN message correctly (level, timestamp, message).
func TestLoggerBasicFormatting(t *testing.T) {
	logger := New(DEBUG) // accept all levels
	r, w, _ := os.Pipe()
	logger.out = w // redirect output into w
	logger.Warn("something happened")
	w.Close()

	// Read output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// It should contain the colored "WARN" tag
	if !strings.Contains(output, "WARN") {
		t.Errorf("expected 'WARN' in output, got: %q", output)
	}

	// It should contain our message
	if !strings.Contains(output, "something happened") {
		t.Errorf("expected message 'something happened' in output, got: %q", output)
	}

	// Timestamp should match pattern "YYYY-MM-DD HH:MM:SS"
	// Example: 2025-06-02 15:04:05
	re := regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\]`)
	if !re.MatchString(output) {
		t.Errorf("timestamp did not match expected format, got: %q", output)
	}
}

// Test that messages below the logger’s level are suppressed.
func TestLoggerLevelFiltering(t *testing.T) {
	logger := New(WARN) // only WARN and above
	r, w, _ := os.Pipe()
	logger.out = w

	logger.Info("info should not appear")
	logger.Debug("debug should not appear")
	logger.Warn("this is a warning")
	logger.Error("this is an error")
	w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if strings.Contains(output, "info should not appear") {
		t.Error("Info-level message should have been filtered out but was printed")
	}
	if strings.Contains(output, "debug should not appear") {
		t.Error("Debug-level message should have been filtered out but was printed")
	}
	if !strings.Contains(output, "this is a warning") {
		t.Error("Warn-level message should have appeared but did not")
	}
	if !strings.Contains(output, "this is an error") {
		t.Error("Error-level message should have appeared but did not")
	}
}

// Test the global-wrapper functions and SetGlobal.
func TestGlobalLoggerFunctions(t *testing.T) {
	// Create a fresh logger at DEBUG
	custom := New(DEBUG)
	r, w, _ := os.Pipe()
	custom.out = w

	// Override the global logger
	SetGlobal(custom)

	// Use the global wrappers
	Debug("global debug")
	Info("global info")
	Warn("global warn")
	Error("global error")
	w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "global debug") {
		t.Error("expected 'global debug' in output")
	}
	if !strings.Contains(output, "global info") {
		t.Error("expected 'global info' in output")
	}
	if !strings.Contains(output, "global warn") {
		t.Error("expected 'global warn' in output")
	}
	if !strings.Contains(output, "global error") {
		t.Error("expected 'global error' in output")
	}
}

// Test that WithContext and FromContext correctly return the logger stored in context.
func TestContextLogger(t *testing.T) {
	// Create a new logger at INFO
	ctxLogger := New(INFO)
	r, w, _ := os.Pipe()
	ctxLogger.out = w

	// Put it into a context
	ctx := WithContext(context.Background(), ctxLogger)

	// Retrieve from context and log
	retrieved := FromContext(ctx)
	retrieved.Info("info via context")
	retrieved.Debug("debug via context (should not appear, since level is INFO)")

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "info via context") {
		t.Error("expected 'info via context' in output, but did not see it")
	}
	if strings.Contains(output, "debug via context") {
		t.Error("did not expect 'debug via context' (below INFO level) to appear")
	}
}

// Test Init(): setting BARELOG_LEVEL should configure the global logger accordingly.
func TestInitFromEnv(t *testing.T) {
	// 1) Test with BARELOG_LEVEL=debug
	os.Setenv("BARELOG_LEVEL", "debug")
	defer os.Unsetenv("BARELOG_LEVEL")

	// Capture stdout from the global logger
	output := captureStdout(func() {
		Init()             // reads the env var
		Debug("helloDIAG") // since level should be DEBUG, this must appear
		Info("helloINFO")
	})

	if !strings.Contains(output, "helloDIAG") {
		t.Error("with BARELOG_LEVEL=debug, Debug(...) should have been printed")
	}
	if !strings.Contains(output, "helloINFO") {
		t.Error("with BARELOG_LEVEL=debug, Info(...) should have been printed")
	}

	// 2) Test with BARELOG_LEVEL=error (only errors should appear)
	os.Setenv("BARELOG_LEVEL", "error")
	output2 := captureStdout(func() {
		Init()
		Info("shouldNOTappear")
		Error("onlyError")
	})
	if strings.Contains(output2, "shouldNOTappear") {
		t.Error("with BARELOG_LEVEL=error, Info(...) should not have been printed")
	}
	if !strings.Contains(output2, "onlyError") {
		t.Error("with BARELOG_LEVEL=error, Error(...) should have been printed")
	}

	// 3) Test with unknown level (falls back to INFO with a stderr warning). We’ll just check INFO-level behavior.
	os.Setenv("BARELOG_LEVEL", "foobar")
	// redirect stderr so we don’t spam test output
	oldStderr := os.Stderr
	devNull, _ := os.Open(os.DevNull)
	os.Stderr = devNull
	defer func() {
		os.Stderr = oldStderr
	}()
	output3 := captureStdout(func() {
		Init()
		Info("fallbackWorks")
		Debug("debugShouldNotAppear")
	})
	if !strings.Contains(output3, "fallbackWorks") {
		t.Error("with unknown BARELOG_LEVEL, should fall back to INFO and print Info(...)")
	}
	if strings.Contains(output3, "debugShouldNotAppear") {
		t.Error("with fallback=INFO, Debug(...) should not appear")
	}
}
