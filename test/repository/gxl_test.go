package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	parser "github.com/ryo-arima/goxcel/pkg/repository"
	"github.com/ryo-arima/goxcel/pkg/util"
)

func TestReadGxlFromFile_Minimal(t *testing.T) {
	// Use test data file under test/.testdata
	path := filepath.Join("..", ".testdata", "minimal.gxl")

	lg := util.NewLogger(util.LoggerConfig{Component: "test", Service: "repo", Level: "DEBUG", Structured: false, Output: "stdout"})
	gxl, err := parser.ReadGxlFromFile(path, lg)
	if err != nil {
		t.Fatalf("ReadGxlFromFile: %v", err)
	}
	if len(gxl.Sheets) != 1 {
		t.Fatalf("sheets = %d, want 1", len(gxl.Sheets))
	}
	if gxl.Sheets[0].Name != "S1" {
		t.Fatalf("sheet name = %q, want S1", gxl.Sheets[0].Name)
	}
	// Grid rows should be parsed (2 rows)
	found := false
	for _, n := range gxl.Sheets[0].Nodes {
		if grid, ok := n.(model.GridTag); ok {
			found = true
			want := []model.GridRowTag{{Cells: []string{"A", "B"}}, {Cells: []string{"1", "2"}}}
			if diff := cmp.Diff(want, grid.Rows); diff != "" {
				t.Fatalf("grid rows mismatch (-want +got):\n%s", diff)
			}
		}
	}
	if !found {
		t.Fatalf("expected GridTag in nodes")
	}
}

func TestGxlRepository_Methods_FormatAndRead(t *testing.T) {
	// Success path
	path := filepath.Join("..", ".testdata", "minimal.gxl")
	conf := config.NewBaseConfigWithFile(path)
	repo := parser.NewGxlRepository(conf)
	if _, err := repo.FormatGxl(); err != nil {
		t.Fatalf("FormatGxl: %v", err)
	}
	if _, err := repo.ReadGxl(); err != nil {
		t.Fatalf("ReadGxl: %v", err)
	}

	// Error path: empty file path
	bad := config.NewBaseConfigWithFile("")
	repo2 := parser.NewGxlRepository(bad)
	if _, err := repo2.FormatGxl(); err == nil {
		t.Fatalf("expected error on empty path in FormatGxl")
	}
	if _, err := repo2.ReadGxl(); err == nil {
		t.Fatalf("expected error on empty path in ReadGxl")
	}
}

func TestFormatGxl_AlignsStyledGrid(t *testing.T) {
	path := filepath.Join("..", ".testdata", "style_grid.gxl")
	conf := config.NewBaseConfigWithFile(path)
	repo := parser.NewGxlRepository(conf)
	b, err := repo.FormatGxl()
	if err != nil {
		t.Fatalf("FormatGxl: %v", err)
	}
	if len(b) == 0 || string(b[:5]) != "<?xml" {
		t.Fatalf("unexpected formatted output prefix: %q", string(b[:5]))
	}
	// Lightweight check that grid pipes exist after formatting
	if !contains(string(b), "| A | B |") {
		t.Fatalf("expected aligned grid line in output")
	}
}

func TestParse_ControlFlowNodes(t *testing.T) {
	path := filepath.Join("..", ".testdata", "control_flow.gxl")
	lg := util.NewLogger(util.LoggerConfig{Component: "test", Service: "repo", Level: "DEBUG", Structured: false, Output: "stdout"})
	gxl, err := parser.ReadGxlFromFile(path, lg)
	if err != nil {
		t.Fatalf("ReadGxlFromFile: %v", err)
	}
	if len(gxl.Sheets) != 1 {
		t.Fatalf("sheets=%d, want 1", len(gxl.Sheets))
	}
	st := gxl.Sheets[0]
	var hasFor, hasIf bool
	for _, n := range st.Nodes {
		switch n.(type) {
		case model.ForTag:
			hasFor = true
		case model.IfTag:
			hasIf = true
		}
	}
	if !hasFor || !hasIf {
		t.Fatalf("expected ForTag and IfTag, got for=%v if=%v", hasFor, hasIf)
	}
}

func TestParse_GridColors_Sanitized(t *testing.T) {
	path := filepath.Join("..", ".testdata", "style_colors.gxl")
	lg := util.NewLogger(util.LoggerConfig{Component: "test", Service: "repo", Level: "DEBUG", Structured: false, Output: "stdout"})
	gxl, err := parser.ReadGxlFromFile(path, lg)
	if err != nil {
		t.Fatalf("ReadGxlFromFile: %v", err)
	}
	if len(gxl.Sheets) != 1 {
		t.Fatalf("sheets=%d, want 1", len(gxl.Sheets))
	}
	// Find GridTag and assert colors are uppercased without '#'
	var grid model.GridTag
	found := false
	for _, n := range gxl.Sheets[0].Nodes {
		if gt, ok := n.(model.GridTag); ok {
			grid = gt
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected GridTag present")
	}
	if grid.FontColor != "AABBCC" {
		t.Fatalf("FontColor=%q, want AABBCC", grid.FontColor)
	}
	if grid.FillColor != "00FF00" {
		t.Fatalf("FillColor=%q, want 00FF00", grid.FillColor)
	}
	if grid.BorderColor != "112233" {
		t.Fatalf("BorderColor=%q, want 112233", grid.BorderColor)
	}
}

// local contains helper (avoid importing strings): simple substring check
func contains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestSanitizeColor_GridAttributes(t *testing.T) {
	path := filepath.Join("..", ".testdata", "style_colors.gxl")
	lg := util.NewLogger(util.LoggerConfig{Component: "test", Service: "repo", Level: "DEBUG", Structured: false, Output: "stdout"})
	gxl, err := parser.ReadGxlFromFile(path, lg)
	if err != nil {
		t.Fatalf("ReadGxlFromFile: %v", err)
	}
	if len(gxl.Sheets) == 0 {
		t.Fatalf("no sheets parsed")
	}
	var grid model.GridTag
	found := false
	for _, n := range gxl.Sheets[0].Nodes {
		if g, ok := n.(model.GridTag); ok {
			grid = g
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no GridTag found")
	}
	if grid.FontColor != "AABBCC" { // ensure upcased and no '#'
		t.Errorf("fontColor not sanitized: %q", grid.FontColor)
	}
	if grid.FillColor == "" {
		t.Errorf("fillColor expected non-empty")
	}
	if grid.BorderColor == "" {
		t.Errorf("borderColor expected non-empty")
	}
}
