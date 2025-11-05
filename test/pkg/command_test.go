package command_test

import (
	"os"
	"os/exec"
	"testing"

	command "github.com/ryo-arima/goxcel/pkg"
)

func TestExecute_Help_NoError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"goxcel", "--help"}
	command.Execute()
}

// TestExecute_Error_Exits ensures Execute exits with code 1 on error path
func TestExecute_Error_Exits(t *testing.T) {
	if os.Getenv("WANT_EXECUTE_FAIL") == "1" {
		// In child process: cause cobra to return an error by passing unknown flag
		os.Args = []string{"goxcel", "--definitely-unknown-flag"}
		command.Execute()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_Error_Exits")
	cmd.Env = append(os.Environ(), "WANT_EXECUTE_FAIL=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected error (non-zero exit) from Execute")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if code := ee.ExitCode(); code != 1 {
			t.Fatalf("expected exit code 1, got %d", code)
		}
	}
}
