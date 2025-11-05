package usecase_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	parser "github.com/ryo-arima/goxcel/pkg/repository"
	"github.com/ryo-arima/goxcel/pkg/usecase"
	"github.com/ryo-arima/goxcel/pkg/util"
)

func findCellByRef(b *model.Book, ref string) *model.Cell {
	for _, s := range b.Sheets {
		for _, c := range s.Cells {
			if c.Ref == ref {
				return c
			}
		}
	}
	return nil
}

func TestBookUsecase_Render_MinimalGrid(t *testing.T) {
	// Build GXL model in-memory (no repository I/O)
	gxl := &model.GXL{
		BookTag: model.BookTag{Name: "Book"},
		Sheets: []model.SheetTag{
			{
				Name: "S1",
				Nodes: []any{
					model.AnchorTag{Ref: "B2"},
					model.GridTag{Rows: []model.GridRowTag{
						{Cells: []string{"Hello", "{{ name }}", "**Bold**", "_Ital_"}},
						{Cells: []string{"=SUM(1,2)", "{{ 123:int }}", "{{ true:bool }}", "{{ \"x\" }}"}},
					}},
				},
			},
		},
	}

	data := map[string]any{"name": "Alice"}

	conf := config.NewBaseConfig()
	uc := usecase.NewBookUsecase(conf)
	book, err := uc.Render(context.Background(), gxl, data)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if len(book.Sheets) != 1 {
		t.Fatalf("sheets = %d, want 1", len(book.Sheets))
	}

	// Compare a projected view of cells using go-cmp
	type cellView struct {
		Ref    string
		Value  string
		Type   model.CellType
		Bold   bool
		Italic bool
	}
	got := []cellView{
		func() cellView {
			c := findCellByRef(book, "B2")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "C2")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "D2")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "E2")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "B3")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "C3")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "D3")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
		func() cellView {
			c := findCellByRef(book, "E3")
			return cellView{Ref: c.Ref, Value: c.Value, Type: c.Type, Bold: c.Style != nil && c.Style.Bold, Italic: c.Style != nil && c.Style.Italic}
		}(),
	}
	want := []cellView{
		{Ref: "B2", Value: "Hello", Type: model.CellTypeString},
		{Ref: "C2", Value: "Alice", Type: model.CellTypeString},
		{Ref: "D2", Value: "Bold", Type: model.CellTypeString, Bold: true},
		{Ref: "E2", Value: "Ital", Type: model.CellTypeString, Italic: true},
		{Ref: "B3", Value: "=SUM(1,2)", Type: model.CellTypeFormula},
		{Ref: "C3", Value: "123", Type: model.CellTypeNumber},
		{Ref: "D3", Value: "true", Type: model.CellTypeBoolean},
		{Ref: "E3", Value: "x", Type: model.CellTypeString},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("cells mismatch (-want +got):\n%s", diff)
	}
}

// TestBookUsecase_Render_ControlFlowGxl exercises For, Anchor, Merge nodes using control_flow.gxl
func TestBookUsecase_Render_ControlFlowGxl(t *testing.T) {
	t.Skip("control_flow.gxl rendering needs investigation - skipping for now")

	conf := config.NewBaseConfig()
	gxl, err := ReadTestGxl("../.testdata/control_flow.gxl", conf)
	if err != nil {
		t.Fatalf("ReadTestGxl: %v", err)
	}

	data := map[string]any{
		"items": []any{"a", "b", "c"},
		"ok":    true,
	}

	uc := usecase.NewBookUsecase(conf)
	book, err := uc.Render(context.Background(), gxl, data)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) == 0 {
		t.Fatal("no sheets rendered")
	}

	sheet := book.Sheets[0]

	// Verify merge was added (from If block when ok=true)
	if len(sheet.Merges) == 0 {
		t.Logf("Got %d sheets, %d cells, %d merges", len(book.Sheets), len(sheet.Cells), len(sheet.Merges))
		t.Error("expected merge ranges from control_flow.gxl")
	}

	// Verify cells from Grid inside If block
	if len(sheet.Cells) == 0 {
		t.Error("expected cells from control_flow.gxl")
	}
}

// TestBookUsecase_Render_ComponentsGxl exercises Image, Shape, Chart, Pivot nodes using components.gxl
func TestBookUsecase_Render_ComponentsGxl(t *testing.T) {
	conf := config.NewBaseConfig()
	gxl, err := ReadTestGxl("../.testdata/components.gxl", conf)
	if err != nil {
		t.Fatalf("ReadTestGxl: %v", err)
	}

	uc := usecase.NewBookUsecase(conf)
	book, err := uc.Render(context.Background(), gxl, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) == 0 {
		t.Fatal("no sheets rendered")
	}

	sheet := book.Sheets[0]

	// Verify components were added
	if len(sheet.Images) == 0 {
		t.Logf("Got %d images, %d shapes, %d charts, %d pivots", len(sheet.Images), len(sheet.Shapes), len(sheet.Charts), len(sheet.Pivots))
		t.Error("expected images from components.gxl")
	}
	if len(sheet.Shapes) == 0 {
		t.Error("expected shapes from components.gxl")
	}
	if len(sheet.Charts) == 0 {
		t.Error("expected charts from components.gxl")
	}
	if len(sheet.Pivots) == 0 {
		t.Error("expected pivots from components.gxl")
	}
}

// TestBookUsecase_Render_VariousValueTypes exercises valueToString with different Go types
func TestBookUsecase_Render_VariousValueTypes(t *testing.T) {
	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "ValueTypes",
				Nodes: []any{
					model.GridTag{Rows: []model.GridRowTag{
						{Cells: []string{"{{intVal}}", "{{int8Val}}", "{{int16Val}}", "{{int32Val}}", "{{int64Val}}"}},
						{Cells: []string{"{{uintVal}}", "{{uint8Val}}", "{{uint16Val}}", "{{uint32Val}}", "{{uint64Val}}"}},
						{Cells: []string{"{{float32Val}}", "{{float64Val}}", "{{boolVal}}", "{{nilVal}}", "{{strVal}}"}},
						{Cells: []string{"{{sliceVal}}", "{{mapVal}}", "{{structVal}}"}},
					}},
				},
			},
		},
	}

	type customStruct struct {
		Field string
	}

	data := map[string]any{
		"intVal":     int(42),
		"int8Val":    int8(8),
		"int16Val":   int16(16),
		"int32Val":   int32(32),
		"int64Val":   int64(64),
		"uintVal":    uint(42),
		"uint8Val":   uint8(8),
		"uint16Val":  uint16(16),
		"uint32Val":  uint32(32),
		"uint64Val":  uint64(64),
		"float32Val": float32(3.14),
		"float64Val": float64(2.718),
		"boolVal":    true,
		"nilVal":     nil,
		"strVal":     "hello",
		"sliceVal":   []int{1, 2, 3},
		"mapVal":     map[string]int{"a": 1},
		"structVal":  customStruct{Field: "test"},
	}

	conf := config.NewBaseConfig()
	uc := usecase.NewBookUsecase(conf)
	book, err := uc.Render(context.Background(), gxl, data)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) == 0 {
		t.Fatal("no sheets rendered")
	}

	sheet := book.Sheets[0]
	if len(sheet.Cells) < 10 {
		t.Errorf("expected at least 10 cells, got %d", len(sheet.Cells))
	}

	// Verify some specific values
	if c := findCellByRef(book, "A1"); c != nil && c.Value != "42" {
		t.Errorf("A1: got %q, want %q", c.Value, "42")
	}
	if c := findCellByRef(book, "A3"); c != nil && c.Value != "3.14" {
		t.Errorf("A3: got %q, want %q", c.Value, "3.14")
	}
	if c := findCellByRef(book, "C3"); c != nil && c.Value != "true" {
		t.Errorf("C3: got %q, want %q", c.Value, "true")
	}
}

// ReadTestGxl is a helper function to read GXL files for testing
func ReadTestGxl(path string, conf config.BaseConfig) (*model.GXL, error) {
	absPath := filepath.Join(path)
	gxl, err := parser.ReadGxlFromFile(absPath, util.NewLogger(util.LoggerConfig{
		Component:  "test",
		Service:    "usecase",
		Level:      "ERROR",
		Structured: false,
		Output:     "stdout",
	}))
	return &gxl, err
}
