package controller_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/controller"
	"github.com/ryo-arima/goxcel/pkg/model"
)

func TestGenerateCmd_BasicExecution(t *testing.T) {
	dir := t.TempDir()

	// Use GXL template from .testdata (relative to test/controller)
	gxlPath := filepath.Join("..", ".testdata", "simple.gxl")

	// Generate output
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath, "--output", outputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRunGenerate_DirectCall(t *testing.T) {
	dir := t.TempDir()

	// Use GXL template from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "direct_call.gxl")

	// Call RunGenerate directly
	outputPath := filepath.Join(dir, "output.xlsx")
	if err := controller.RunGenerate(gxlPath, "", outputPath, false); err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRunGenerate_WithJSONData(t *testing.T) {
	dir := t.TempDir()

	// Use GXL template from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "with_variables.gxl")

	// Create JSON data
	jsonContent := `{"name": "Test", "value": 123}`
	jsonPath := filepath.Join(dir, "data.json")
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	// Call RunGenerate with JSON data
	outputPath := filepath.Join(dir, "output.xlsx")
	if err := controller.RunGenerate(gxlPath, jsonPath, outputPath, false); err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRunGenerate_WithYAMLData(t *testing.T) {
	dir := t.TempDir()

	// Use GXL template from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "yaml_template.gxl")

	// Create YAML data
	yamlContent := `title: TestTitle`
	yamlPath := filepath.Join(dir, "data.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml: %v", err)
	}

	// Call RunGenerate with YAML data
	outputPath := filepath.Join(dir, "output.xlsx")
	if err := controller.RunGenerate(gxlPath, yamlPath, outputPath, false); err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRunGenerate_DryRun(t *testing.T) {
	// Use GXL template from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "dry_run.gxl")

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call RunGenerate in dry-run mode
	err := controller.RunGenerate(gxlPath, "", "", true)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Verify summary was printed
	output := buf.String()
	if !strings.Contains(output, "Workbook:") {
		t.Errorf("expected summary output, got: %s", output)
	}
}

func TestRunGenerate_ErrorInvalidTemplate(t *testing.T) {
	// Use invalid GXL from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "invalid.gxl")
	dir := t.TempDir()

	// Should return error
	outputPath := filepath.Join(dir, "output.xlsx")
	err := controller.RunGenerate(gxlPath, "", outputPath, false)
	if err == nil {
		t.Error("expected error for invalid template, got nil")
	}
}

func TestRunGenerate_ErrorInvalidJSON(t *testing.T) {
	dir := t.TempDir()

	// Use GXL template from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "error_test.gxl")

	// Create invalid JSON
	invalidJson := `{invalid json`
	jsonPath := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(jsonPath, []byte(invalidJson), 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	// Should return error
	outputPath := filepath.Join(dir, "output.xlsx")
	err := controller.RunGenerate(gxlPath, jsonPath, outputPath, false)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGenerateCmd_DryRun_WithYAMLData(t *testing.T) {
	dir := t.TempDir()

	// Use GXL template from .testdata (relative path corrected)
	gxlPath := filepath.Join("..", ".testdata", "yaml_dry_run.gxl")

	// Prepare YAML data file
	datap := filepath.Join(dir, "data.yaml")
	if err := os.WriteFile(datap, []byte("ok: true\n"), 0644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath, "--data", datap, "--dry-run"})

	err := cmd.Execute()
	w.Close()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("Execute: %v; out=%s", err, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "Workbook:") {
		t.Fatalf("expected summary output with 'Workbook:', got: %s", output)
	}
	// Check that printBookSummary was called (even with 0 or more sheets)
	if !strings.Contains(output, "sheets") {
		t.Errorf("expected sheets count in summary, got: %s", output)
	}
}

func TestPrintBookSummary(t *testing.T) {
	// Create a test book
	book := model.NewBook()

	sheet1 := model.NewSheet("TestSheet1")
	for i := 0; i < 5; i++ {
		sheet1.AddCell(&model.Cell{Ref: "A" + string(rune('1'+i)), Value: "test", Type: model.CellTypeString})
	}
	sheet1.AddMerge(model.Merge{Range: "A1:B2"})
	sheet1.AddMerge(model.Merge{Range: "C3:D4"})
	sheet1.AddImage(model.Image{Ref: "E5", Source: "test.png"})
	book.AddSheet(sheet1)

	sheet2 := model.NewSheet("TestSheet2")
	for i := 0; i < 10; i++ {
		sheet2.AddCell(&model.Cell{Ref: "A" + string(rune('1'+i)), Value: "test", Type: model.CellTypeString})
	}
	sheet2.AddShape(model.Shape{Ref: "F6", Kind: "rectangle"})
	sheet2.AddShape(model.Shape{Ref: "G7", Kind: "circle"})
	sheet2.AddShape(model.Shape{Ref: "H8", Kind: "triangle"})
	sheet2.AddChart(model.Chart{Ref: "I9", Type: "bar"})
	sheet2.AddPivot(model.PivotTable{Ref: "J10", SourceRange: "A1:D10"})
	book.AddSheet(sheet2)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	controller.PrintBookSummary(book)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output
	if !strings.Contains(output, "Workbook: 2 sheets") {
		t.Errorf("expected 'Workbook: 2 sheets', got: %s", output)
	}
	if !strings.Contains(output, "TestSheet1") {
		t.Errorf("expected 'TestSheet1' in output, got: %s", output)
	}
	if !strings.Contains(output, "TestSheet2") {
		t.Errorf("expected 'TestSheet2' in output, got: %s", output)
	}
	if !strings.Contains(output, "5 cells") {
		t.Errorf("expected '5 cells' for TestSheet1, got: %s", output)
	}
	if !strings.Contains(output, "10 cells") {
		t.Errorf("expected '10 cells' for TestSheet2, got: %s", output)
	}
}

func TestInitGenerateCmd_FlagsExist(t *testing.T) {
	cmd := controller.InitGenerateCmd()
	if cmd == nil {
		t.Fatal("InitGenerateCmd returned nil")
	}

	// Verify flags exist
	flags := []string{"template", "data", "output", "dry-run"}
	for _, flag := range flags {
		if f := cmd.Flags().Lookup(flag); f == nil {
			t.Errorf("flag %q not found", flag)
		}
	}
}

func TestGenerateCmd_WithJSONData(t *testing.T) {
	dir := t.TempDir()

	// Create GXL template with variable
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Data">
    <grid ref="A1">
      <row>
        <cell>{{name}}</cell>
        <cell>{{value}}</cell>
      </row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "template.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Create JSON data file
	jsonContent := `{
  "name": "TestName",
  "value": 12345
}`
	jsonPath := filepath.Join(dir, "data.json")
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	// Generate with data
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath, "--data", jsonPath, "--output", outputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestGenerateCmd_WithYMLExtension(t *testing.T) {
	dir := t.TempDir()

	// Create GXL template
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="YML">
    <grid ref="A1">
      <row><cell>{{key}}</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Create YML data file
	ymlContent := `key: value`
	ymlPath := filepath.Join(dir, "data.yml")
	if err := os.WriteFile(ymlPath, []byte(ymlContent), 0644); err != nil {
		t.Fatalf("failed to write yml: %v", err)
	}

	// Generate with YML extension
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath, "--data", ymlPath, "--output", outputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestGenerateCmd_AutoDetectFormat(t *testing.T) {
	dir := t.TempDir()

	// Create GXL template
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Auto">
    <grid ref="A1">
      <row><cell>{{data}}</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Create data file with non-standard extension
	dataContent := `{"data": "test"}`
	dataPath := filepath.Join(dir, "data.txt")
	if err := os.WriteFile(dataPath, []byte(dataContent), 0644); err != nil {
		t.Fatalf("failed to write data: %v", err)
	}

	// Generate (should auto-detect JSON)
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath, "--data", dataPath, "--output", outputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestGenerateCmd_ErrorInvalidJSON(t *testing.T) {
	dir := t.TempDir()

	// Create valid GXL
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>{{value}}</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Create invalid JSON
	invalidJson := `{invalid json`
	jsonPath := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(jsonPath, []byte(invalidJson), 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	// Should fail with invalid JSON
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath, "--data", jsonPath, "--output", outputPath})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGenerateCmd_PositionalArgument(t *testing.T) {
	dir := t.TempDir()

	// Create GXL template
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="PosArg">
    <grid ref="A1">
      <row><cell>Position</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Use positional argument for template
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{gxlPath, "--output", outputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestGenerateCmd_ShorthandFlags(t *testing.T) {
	dir := t.TempDir()

	// Create GXL template
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Short">
    <grid ref="A1"><row><cell>{{val}}</cell></row></grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Create data
	jsonContent := `{"val": "test"}`
	jsonPath := filepath.Join(dir, "data.json")
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	// Use shorthand flags
	outputPath := filepath.Join(dir, "output.xlsx")
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"-t", gxlPath, "-d", jsonPath, "-o", outputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestGenerateCmd_EmptyOutputNoDryRun(t *testing.T) {
	dir := t.TempDir()

	// Create GXL template
	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Empty">
    <grid ref="A1">
      <row><cell>Test</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	// Run without output path (should print summary)
	cmd := controller.InitGenerateCmd()
	cmd.SetArgs([]string{"--template", gxlPath})

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
		t.Errorf("expected summary output, got: %s", buf.String())
	}
}

// TestRunGenerate_ExistingOutputDirectory tests writing to an existing directory
func TestRunGenerate_ExistingOutputDirectory(t *testing.T) {
	dir := t.TempDir()

	// Create output directory beforehand
	outDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>Existing Dir</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	outputPath := filepath.Join(outDir, "output.xlsx")

	// Execute - should succeed even though directory exists
	err := controller.RunGenerate(gxlPath, "", outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// TestRunGenerate_NestedOutputDirectory tests creating nested directories
func TestRunGenerate_NestedOutputDirectory(t *testing.T) {
	dir := t.TempDir()

	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>Nested Dir</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Output path with multiple nested directories
	outputPath := filepath.Join(dir, "level1", "level2", "level3", "output.xlsx")

	err := controller.RunGenerate(gxlPath, "", outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// TestRunGenerate_EmptyDataPath tests with empty data path
func TestRunGenerate_EmptyDataPath(t *testing.T) {
	dir := t.TempDir()

	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>No Data</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	outputPath := filepath.Join(dir, "output.xlsx")

	// Test with empty string data path
	err := controller.RunGenerate(gxlPath, "", outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate failed: %v", err)
	}

	// Test with whitespace-only data path
	err = controller.RunGenerate(gxlPath, "   ", outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate with whitespace data path failed: %v", err)
	}
}

// TestRunGenerate_AutoDetectYAML tests auto-detection with .yml extension
func TestRunGenerate_AutoDetectYAML(t *testing.T) {
	dir := t.TempDir()

	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>{{ name }}</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Test with .yml extension
	ymlData := `name: YAML Test`
	ymlPath := filepath.Join(dir, "data.yml")
	if err := os.WriteFile(ymlPath, []byte(ymlData), 0644); err != nil {
		t.Fatalf("failed to write yml: %v", err)
	}

	outputPath := filepath.Join(dir, "output.xlsx")

	err := controller.RunGenerate(gxlPath, ymlPath, outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate with .yml failed: %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// TestRunGenerate_AutoDetectNoExtension tests auto-detection without file extension
func TestRunGenerate_AutoDetectNoExtension(t *testing.T) {
	dir := t.TempDir()

	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>{{ value }}</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Test with no extension (JSON format)
	jsonData := `{"value": "Auto JSON"}`
	dataPath := filepath.Join(dir, "datafile")
	if err := os.WriteFile(dataPath, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to write data: %v", err)
	}

	outputPath := filepath.Join(dir, "output.xlsx")

	err := controller.RunGenerate(gxlPath, dataPath, outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate with no extension failed: %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// TestRunGenerate_AutoDetectFallbackYAML tests auto-detection falling back to YAML
func TestRunGenerate_AutoDetectFallbackYAML(t *testing.T) {
	dir := t.TempDir()

	gxlContent := `<?xml version="1.0" encoding="UTF-8"?>
<gxl>
  <sheet name="Test">
    <grid ref="A1">
      <row><cell>{{ value }}</cell></row>
    </grid>
  </sheet>
</gxl>`
	gxlPath := filepath.Join(dir, "test.gxl")
	if err := os.WriteFile(gxlPath, []byte(gxlContent), 0644); err != nil {
		t.Fatalf("failed to write gxl: %v", err)
	}

	// Test with no extension (YAML format)
	yamlData := `value: Auto YAML`
	dataPath := filepath.Join(dir, "datafile")
	if err := os.WriteFile(dataPath, []byte(yamlData), 0644); err != nil {
		t.Fatalf("failed to write data: %v", err)
	}

	outputPath := filepath.Join(dir, "output.xlsx")

	err := controller.RunGenerate(gxlPath, dataPath, outputPath, false)
	if err != nil {
		t.Fatalf("RunGenerate with YAML fallback failed: %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}
