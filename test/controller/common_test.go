package controller_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/controller"
)

func TestGenerateCmd_DryRun_WithYAMLData(t *testing.T) {
	cmd := controller.InitGenerateCmd()
	// Prepare YAML data file
	dir := t.TempDir()
	datap := filepath.Join(dir, "data.yaml")
	if err := os.WriteFile(datap, []byte("ok: true\n"), 0644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}
	in := filepath.Join("..", ".testdata", "minimal.gxl")

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	cmd.SetArgs([]string{"--template", in, "--data", datap, "--dry-run"})
	if err := cmd.Execute(); err != nil {
		w.Close()
		io := new(bytes.Buffer)
		_, _ = io.ReadFrom(r)
		t.Fatalf("Execute: %v; out=%s", err, io.String())
	}
	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	if !bytes.Contains(buf.Bytes(), []byte("Workbook:")) {
		t.Fatalf("expected summary output, got: %s", buf.String())
	}
}
