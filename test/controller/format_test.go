package controller_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/controller"
)

func TestInitFormatCmd_WritesOutput(t *testing.T) {
	cmd := controller.InitFormatCmd()
	in := filepath.Join("..", ".testdata", "minimal.gxl")
	out := filepath.Join(t.TempDir(), "fmt.gxl")
	cmd.SetArgs([]string{"-o", out, in})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output not written: %v", err)
	}
}
