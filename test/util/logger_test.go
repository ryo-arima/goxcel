package util_test

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ryo-arima/goxcel/pkg/util"
)

type logEntry struct {
	Level   string                 `json:"level"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields"`
}

func newTempLogger(t *testing.T) (util.Logger, string) {
	t.Helper()
	tmpFile, err := os.CreateTemp(t.TempDir(), "logger-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	// Close immediately; logger will open by path
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()

	lg := util.NewLogger(util.LoggerConfig{
		Component:    "test",
		Service:      "logger",
		Level:        "DEBUG",
		Structured:   true,
		EnableCaller: false,
		Output:       tmpPath,
	})
	return lg, tmpPath
}

func readLastJSONLine(t *testing.T, path string) logEntry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open log file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var last string
	for scanner.Scan() {
		last = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}
	var e logEntry
	if err := json.Unmarshal([]byte(last), &e); err != nil {
		t.Fatalf("unmarshal: %v; line=%q", err, last)
	}
	return e
}

func TestLogger_DEBUG_IncludesFields(t *testing.T) {
	lg, path := newTempLogger(t)
	lg.DEBUG(util.SYS1, "debug with fields", map[string]interface{}{"k": "v"})
	// tiny delay to ensure write
	time.Sleep(10 * time.Millisecond)
	entry := readLastJSONLine(t, path)
	if entry.Level != "DEBUG" {
		t.Fatalf("level = %s, want DEBUG", entry.Level)
	}
	if entry.Code != util.SYS1.Code {
		t.Fatalf("code = %s, want %s", entry.Code, util.SYS1.Code)
	}
	if entry.Fields == nil || entry.Fields["k"] != "v" {
		t.Fatalf("fields missing or wrong: %+v", entry.Fields)
	}
}

func TestLogger_INFO_NoFieldsProvided(t *testing.T) {
	lg, path := newTempLogger(t)
	lg.INFO(util.SYS1, "info message")
	time.Sleep(10 * time.Millisecond)
	entry := readLastJSONLine(t, path)
	if entry.Level != "INFO" {
		t.Fatalf("level = %s, want INFO", entry.Level)
	}
	if entry.Fields != nil {
		t.Fatalf("expected no fields, got: %+v", entry.Fields)
	}
}

func TestLogger_WARN_NoFields(t *testing.T) {
	lg, path := newTempLogger(t)
	lg.WARN(util.SYS1, "warn message")
	time.Sleep(10 * time.Millisecond)
	e := readLastJSONLine(t, path)
	if e.Level != "WARN" {
		t.Fatalf("level = %s, want WARN", e.Level)
	}
	if e.Fields != nil {
		t.Fatalf("expected no fields, got: %+v", e.Fields)
	}
}

func TestLogger_ERROR_NoFields(t *testing.T) {
	lg, path := newTempLogger(t)
	lg.ERROR(util.SYS3, "error message")
	time.Sleep(10 * time.Millisecond)
	e := readLastJSONLine(t, path)
	if e.Level != "ERROR" {
		t.Fatalf("level = %s, want ERROR", e.Level)
	}
	if e.Fields != nil {
		t.Fatalf("expected no fields, got: %+v", e.Fields)
	}
}

func TestLogger_FATAL_Exits(t *testing.T) {
	if os.Getenv("GOXCEL_HELPER_FATAL") == "1" {
		lg := util.NewLogger(util.LoggerConfig{Component: "test", Service: "util", Level: "DEBUG", Structured: false, Output: "stdout"})
		lg.FATAL(util.SYS3, "fatal test")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run", "TestLogger_FATAL_Exits")
	cmd.Env = append(os.Environ(), "GOXCEL_HELPER_FATAL=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected process to exit with error")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if ee.ExitCode() != 1 {
			t.Fatalf("expected exit code 1, got %d", ee.ExitCode())
		}
	} else {
		t.Fatalf("unexpected error type: %T: %v", err, err)
	}
}
