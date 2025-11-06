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

func TestLogger_LevelFiltering(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "logger-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()

	// Create logger with INFO level (should filter out DEBUG)
	lg := util.NewLogger(util.LoggerConfig{
		Component:    "test",
		Service:      "level",
		Level:        "INFO",
		Structured:   true,
		EnableCaller: false,
		Output:       tmpPath,
	})

	lg.DEBUG(util.SYS1, "this should not appear")
	lg.INFO(util.SYS1, "this should appear")
	time.Sleep(10 * time.Millisecond)

	// Read all lines
	f, err := os.Open(tmpPath)
	if err != nil {
		t.Fatalf("open log file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() {
		lines++
		var e logEntry
		if err := json.Unmarshal(scanner.Bytes(), &e); err == nil {
			if e.Level == "DEBUG" {
				t.Error("DEBUG log should have been filtered")
			}
		}
	}

	if lines == 0 {
		t.Error("expected at least one log line")
	}
}

func TestLogger_OutputStdout(t *testing.T) {
	// Test that stdout output doesn't crash
	lg := util.NewLogger(util.LoggerConfig{
		Component:    "test",
		Service:      "stdout",
		Level:        "INFO",
		Structured:   false,
		EnableCaller: false,
		Output:       "stdout",
	})

	// These should not crash
	lg.INFO(util.SYS1, "stdout test")
	lg.WARN(util.SYS2, "warning to stdout")
	lg.ERROR(util.SYS3, "error to stdout")
}

func TestLogger_StructuredFormatting(t *testing.T) {
	lg, path := newTempLogger(t)

	// Test with various field types
	lg.INFO(util.FSR1, "structured test", map[string]interface{}{
		"string": "value",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"nil":    nil,
		"array":  []string{"a", "b", "c"},
		"map":    map[string]int{"x": 1, "y": 2},
	})
	time.Sleep(10 * time.Millisecond)

	entry := readLastJSONLine(t, path)
	if entry.Level != "INFO" {
		t.Errorf("level = %s, want INFO", entry.Level)
	}
	if entry.Fields == nil {
		t.Fatal("expected fields to be present")
	}
	if entry.Fields["string"] != "value" {
		t.Errorf("string field = %v, want 'value'", entry.Fields["string"])
	}
	if entry.Fields["int"] != float64(42) { // JSON unmarshals numbers as float64
		t.Errorf("int field = %v, want 42", entry.Fields["int"])
	}
	if entry.Fields["bool"] != true {
		t.Errorf("bool field = %v, want true", entry.Fields["bool"])
	}
}

func TestLogger_UnstructuredFormat(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "logger-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()

	// Create unstructured logger
	lg := util.NewLogger(util.LoggerConfig{
		Component:    "test",
		Service:      "unstructured",
		Level:        "DEBUG",
		Structured:   false,
		EnableCaller: false,
		Output:       tmpPath,
	})

	lg.INFO(util.CC1, "unstructured message", map[string]interface{}{"key": "value"})
	time.Sleep(10 * time.Millisecond)

	// Read the log file
	f, err := os.Open(tmpPath)
	if err != nil {
		t.Fatalf("open log file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			found = true
			// Unstructured log should not be JSON
			var e logEntry
			if json.Unmarshal([]byte(line), &e) == nil {
				t.Error("unstructured log should not be valid JSON")
			}
		}
	}

	if !found {
		t.Error("expected at least one log line")
	}
}

func TestLogger_WithCaller(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "logger-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()

	// Create logger with caller enabled
	lg := util.NewLogger(util.LoggerConfig{
		Component:    "test",
		Service:      "caller",
		Level:        "DEBUG",
		Structured:   true,
		EnableCaller: true,
		Output:       tmpPath,
	})

	lg.DEBUG(util.SYS1, "with caller info")
	time.Sleep(10 * time.Millisecond)

	entry := readLastJSONLine(t, tmpPath)
	if entry.Level != "DEBUG" {
		t.Errorf("level = %s, want DEBUG", entry.Level)
	}
	// Note: Caller information would be in a "caller" or "source" field
	// The exact field name depends on implementation
}

func TestMCode_FormatWithOptional(t *testing.T) {
	tests := []struct {
		name    string
		code    util.MCode
		message string
		wantMsg string
	}{
		{
			name:    "with optional message",
			code:    util.SYS1,
			message: "additional info",
			wantMsg: "Application started: additional info",
		},
		{
			name:    "without optional message",
			code:    util.FSR1,
			message: "",
			wantMsg: "File read success",
		},
		{
			name:    "empty optional message",
			code:    util.CC1,
			message: "",
			wantMsg: util.CC1.Message,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.FormatWithOptional(tt.message)
			if got != tt.wantMsg {
				t.Errorf("FormatWithOptional() = %s, want %s", got, tt.wantMsg)
			}
		})
	}
}

func TestMCode_String(t *testing.T) {
	tests := []struct {
		code     util.MCode
		wantCode string
		wantMsg  string
	}{
		{util.SYS1, "SY-S1", "Application started"},
		{util.SYS2, "SY-S2", "Application terminated successfully"},
		{util.SYS3, "SY-S3", "Application terminated with error"},
		{util.CC1, "C-C1", "Command execution success"},
		{util.FSR1, "FS-R1", "File read success"},
		{util.FSR2, "FS-R2", "File read failed"},
		{util.RP1, "R-P1", "Parse operation success"},
		{util.RP2, "R-P2", "Parse operation failed"},
	}

	for _, tt := range tests {
		t.Run(tt.wantCode, func(t *testing.T) {
			if tt.code.Code != tt.wantCode {
				t.Errorf("MCode.Code = %s, want %s", tt.code.Code, tt.wantCode)
			}
			if tt.code.Message != tt.wantMsg {
				t.Errorf("MCode.Message = %s, want %s", tt.code.Message, tt.wantMsg)
			}
		})
	}
}
