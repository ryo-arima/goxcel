//go:build skiptests

package command_test

import (
	"os"
	"testing"

	command "github.com/ryo-arima/goxcel/pkg"
)

func TestExecute_Help_NoError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"goxcel", "--help"}
	command.Execute()
}
