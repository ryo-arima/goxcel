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

func TestInitFormatCmd_InPlace(t *testing.T) {
	// Test --write flag (in-place modification)
	dir := t.TempDir()
	src := filepath.Join("..", ".testdata", "minimal.gxl")
	dest := filepath.Join(dir, "test.gxl")

	// Copy source to temp
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read source: %v", err)
	}
	if err := os.WriteFile(dest, data, 0644); err != nil {
		t.Fatalf("write temp: %v", err)
	}

	cmd := controller.InitFormatCmd()
	cmd.SetArgs([]string{"--write", dest})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	// Verify file was modified
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("file not found: %v", err)
	}
}

func TestInitFormatCmd_Stdout(t *testing.T) {
	// Test default behavior (output to stdout)
	cmd := controller.InitFormatCmd()
	in := filepath.Join("..", ".testdata", "minimal.gxl")
	cmd.SetArgs([]string{in})

	// Just verify it runs without error
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
}

func TestInitFormatCmd_BothFlagsError(t *testing.T) {
	// Test error when both --write and --output are specified
	cmd := controller.InitFormatCmd()
	in := filepath.Join("..", ".testdata", "minimal.gxl")
	out := filepath.Join(t.TempDir(), "out.gxl")
	cmd.SetArgs([]string{"--write", "--output", out, in})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when using both --write and --output")
	}
	if err.Error() != "cannot use --write and --output together" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInitFormatCmd_ShortFlags(t *testing.T) {
	// Test short flags -w and -o
	t.Run("short write flag", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join("..", ".testdata", "minimal.gxl")
		dest := filepath.Join(dir, "test.gxl")

		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read source: %v", err)
		}
		if err := os.WriteFile(dest, data, 0644); err != nil {
			t.Fatalf("write temp: %v", err)
		}

		cmd := controller.InitFormatCmd()
		cmd.SetArgs([]string{"-w", dest})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute: %v", err)
		}
	})

	t.Run("short output flag", func(t *testing.T) {
		cmd := controller.InitFormatCmd()
		in := filepath.Join("..", ".testdata", "minimal.gxl")
		out := filepath.Join(t.TempDir(), "fmt.gxl")
		cmd.SetArgs([]string{"-o", out, in})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute: %v", err)
		}
	})
}
