package usecase_test

import (
	"context"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	usecase "github.com/ryo-arima/goxcel/pkg/usecase"
)

// build a simple sheet with various nodes to exercise renderer paths
func TestRenderSheet_GridAnchorMergeAndComponents(t *testing.T) {
	conf := config.NewBaseConfig()
	r := usecase.NewBookUsecase(conf)

	// SheetTag with anchor, absolute grid, sequential grid, merge, and components
	st := model.SheetTag{
		Name: "S",
		Nodes: []any{
			model.AnchorTag{Ref: "B2"},
			model.GridTag{ // absolute grid (with style hints)
				Ref:         "C3",
				Rows:        []model.GridRowTag{{Cells: []string{"**B**", "_i_"}}},
				FontName:    "Arial",
				FontSize:    12,
				FontColor:   "112233",
				FillColor:   "FFEEDD",
				BorderStyle: "thin",
				BorderColor: "445566",
				BorderSides: "top,left",
			},
			model.GridTag{ // sequential grid (base at anchor B2)
				Rows: []model.GridRowTag{{Cells: []string{"x", "y"}}},
			},
			model.MergeTag{Range: "A1:C1"},
			model.ImageTag{Ref: "E5", Src: "img.png", Width: 10, Height: 20},
			model.ShapeTag{Ref: "F6", Kind: "rect", Text: "t", Width: 5, Height: 6, Style: "s"},
			model.ChartTag{Ref: "G7", Type: "bar", DataRange: "A1:B2", Title: "T", Width: 8, Height: 9},
			model.PivotTag{Ref: "H8", SourceRange: "A1:B5", Rows: "r1, r2", Columns: "c1", Values: "v1, v2", Filters: "f1, f2"},
		},
		Config: &model.SheetConfigTag{DefaultRowHeight: 16, DefaultColumnWidth: 9.0},
	}

	// Wrap into GXL and render with minimal context
	gxl := &model.GXL{Sheets: []model.SheetTag{st}}
	book, err := r.Render(context.Background(), gxl, map[string]any{"data": 1})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if len(book.Sheets) != 1 {
		t.Fatalf("sheets=%d", len(book.Sheets))
	}

	sh := book.Sheets[0]
	if sh.Config.DefaultRowHeight != 16 || sh.Config.DefaultColumnWidth != 9.0 {
		t.Errorf("sheet config not applied: %+v", sh.Config)
	}
	// Ensure components were added
	if len(sh.Merges) != 1 || len(sh.Images) != 1 || len(sh.Shapes) != 1 || len(sh.Charts) != 1 || len(sh.Pivots) != 1 {
		t.Errorf("unexpected components counts: merges=%d images=%d shapes=%d charts=%d pivots=%d", len(sh.Merges), len(sh.Images), len(sh.Shapes), len(sh.Charts), len(sh.Pivots))
	}

	// Check style merge from markdown (** and _)
	// Find the cell rendered from "**B**" at C3 anchor
	var styled *model.Cell
	for _, c := range sh.Cells {
		if c.Ref == "C3" {
			styled = c
			break
		}
	}
	if styled == nil || styled.Style == nil || !styled.Style.Bold {
		t.Errorf("expected bold style at C3, got %+v", styled)
	}
	if styled.Style.Italic { // first cell should be bold only
		t.Errorf("unexpected italic at C3")
	}
}

func TestRenderSheet_ForLoops_ArrayAndMap_AndErrors(t *testing.T) {
	conf := config.NewBaseConfig()
	sr := usecase.NewBookUsecase(conf)

	// Body prints name field
	body := []any{model.GridRowTag{Cells: []string{"{{ item.name }}"}}}
	st := model.SheetTag{
		Name: "Loop",
		Nodes: []any{
			model.ForTag{Each: "item in items", Body: body},                                              // []any case
			model.ForTag{Each: "m in maps", Body: []any{model.GridRowTag{Cells: []string{"{{ m.k }}"}}}}, // []map case
			model.ForTag{Each: "bad syntax", Body: body},                                                 // error path
		},
	}
	// Render
	gxl := &model.GXL{Sheets: []model.SheetTag{st}}
	data := map[string]any{
		"items": []any{map[string]any{"name": "A"}, map[string]any{"name": "B"}},
		"maps":  []map[string]any{{"k": "X"}, {"k": "Y"}},
	}
	_, _ = sr.Render(context.Background(), gxl, data) // ignore error from bad syntax to exercise branch
}

// TestRenderSheet_AllNodeTypes covers Anchor, Merge, Image, Shape, Chart, Pivot nodes
func TestRenderSheet_AllNodeTypes(t *testing.T) {
	conf := config.NewBaseConfig()
	uc := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "AllNodes",
				Nodes: []any{
					// Anchor node
					model.AnchorTag{Ref: "B2"},
					// Grid after anchor
					model.GridTag{
						Ref:     "B2",
						Content: "| Cell1 | Cell2 |\n",
						Rows:    []model.GridRowTag{{Cells: []string{"Cell1", "Cell2"}}},
					},
					// Merge node
					model.MergeTag{Range: "B2:C2"},
					// Image node
					model.ImageTag{
						Ref:    "D5",
						Src:    "/path/to/image.png",
						Width:  100,
						Height: 100,
					},
					// Shape node
					model.ShapeTag{
						Ref:    "F5",
						Kind:   "rectangle",
						Text:   "Shape Text",
						Width:  50,
						Height: 30,
						Style:  "fill:blue",
					},
					// Chart node
					model.ChartTag{
						Ref:       "H5",
						Type:      "bar",
						DataRange: "A1:B10",
						Width:     300,
						Height:    200,
						Title:     "Test Chart",
					},
					// Pivot node
					model.PivotTag{
						Ref:         "K5",
						SourceRange: "A1:D20",
						Rows:        "Category",
						Columns:     "Region",
						Values:      "Sales",
					},
					// For loop with Grid inside
					model.ForTag{
						Each: "item in items",
						Body: []any{
							model.GridTag{
								Content: "| {{item}} |\n",
								Rows:    []model.GridRowTag{{Cells: []string{"{{item}}"}}},
							},
						},
					},
				},
			},
		},
	}

	data := map[string]any{
		"items": []any{"A", "B", "C"},
	}

	book, err := uc.Render(context.Background(), gxl, data)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) == 0 {
		t.Fatal("no sheets rendered")
	}

	sheet := book.Sheets[0]

	// Verify merge was added
	if len(sheet.Merges) == 0 {
		t.Error("expected merge ranges to be added")
	}

	// Verify image was added
	if len(sheet.Images) == 0 {
		t.Error("expected images to be added")
	}

	// Verify shape was added
	if len(sheet.Shapes) == 0 {
		t.Error("expected shapes to be added")
	}

	// Verify chart was added
	if len(sheet.Charts) == 0 {
		t.Error("expected charts to be added")
	}

	// Verify pivot was added
	if len(sheet.Pivots) == 0 {
		t.Error("expected pivots to be added")
	}

	// Verify cells from For loop (3 items)
	itemCells := 0
	for _, cell := range sheet.Cells {
		if cell.Value == "A" || cell.Value == "B" || cell.Value == "C" {
			itemCells++
		}
	}
	if itemCells < 3 {
		t.Errorf("expected at least 3 cells from For loop, got %d", itemCells)
	}
}

// TestRenderSheet_GridWithRef exercises handleGridWithRef path
func TestRenderSheet_GridWithRef(t *testing.T) {
	conf := config.NewBaseConfig()
	uc := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "GridWithRef",
				Nodes: []any{
					model.GridTag{
						Ref: "C3",
						Rows: []model.GridRowTag{
							{Cells: []string{"A", "B"}},
							{Cells: []string{"C", "D"}},
						},
						// Style attributes to exercise gridTagToStyle
						FontName:    "Arial",
						FontSize:    12,
						FontColor:   "FF0000",
						FillColor:   "00FF00",
						BorderStyle: "thin",
						BorderColor: "0000FF",
						BorderSides: "all",
					},
				},
			},
		},
	}

	book, err := uc.Render(context.Background(), gxl, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) == 0 {
		t.Fatal("no sheets rendered")
	}

	sheet := book.Sheets[0]

	// Verify cells at C3, D3, C4, D4
	expectedRefs := map[string]bool{"C3": false, "D3": false, "C4": false, "D4": false}
	for _, cell := range sheet.Cells {
		if _, ok := expectedRefs[cell.Ref]; ok {
			expectedRefs[cell.Ref] = true
		}
	}

	for ref, found := range expectedRefs {
		if !found {
			t.Errorf("cell %s not found", ref)
		}
	}

	// Verify style was applied
	hasStyle := false
	for _, cell := range sheet.Cells {
		if cell.Style != nil && cell.Style.FontName == "Arial" {
			hasStyle = true
			break
		}
	}
	if !hasStyle {
		t.Error("expected cells with style")
	}
}

func TestRenderSheet_LoopConstructs(t *testing.T) {
	// Test handleFor, renderLoop, renderMapLoop, createLoopScope (low coverage)
	conf := config.NewBaseConfig()
	r := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "Loops",
				Nodes: []any{
					// Array loop
					model.ForTag{
						Each: "item in items",
						Body: []any{model.GridTag{Rows: []model.GridRowTag{{Cells: []string{"{{ item }}"}}}}},
					},
					// Map loop
					model.ForTag{
						Each: "k,v in mapping",
						Body: []any{model.GridTag{Rows: []model.GridRowTag{{Cells: []string{"{{ k }}", "{{ v }}"}}}}},
					},
					// Nested loop
					model.ForTag{
						Each: "outer in nested",
						Body: []any{
							model.ForTag{
								Each: "inner in outer.items",
								Body: []any{model.GridTag{Rows: []model.GridRowTag{{Cells: []string{"{{ inner }}"}}}}},
							},
						},
					},
				},
			},
		},
	}

	data := map[string]any{
		"items":   []any{"a", "b", "c"},
		"mapping": map[string]any{"key1": "val1", "key2": "val2"},
		"nested": []any{
			map[string]any{"items": []any{"x", "y"}},
			map[string]any{"items": []any{"z"}},
		},
	}

	book, err := r.Render(context.Background(), gxl, data)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) != 1 {
		t.Fatalf("expected 1 sheet, got %d", len(book.Sheets))
	}

	sheet := book.Sheets[0]
	if len(sheet.Cells) == 0 {
		t.Fatal("expected cells from loops")
	}

	// Verify we have cells from array loop
	hasA := false
	for _, cell := range sheet.Cells {
		if cell.Value == "a" {
			hasA = true
			break
		}
	}
	if !hasA {
		t.Error("expected cell with value 'a' from array loop")
	}
}

func TestRenderSheet_ConditionalRendering(t *testing.T) {
	// Test if tag conditional rendering
	conf := config.NewBaseConfig()
	r := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "Conditionals",
				Nodes: []any{
					model.IfTag{
						Cond: "show",
						Then: []any{model.GridTag{Rows: []model.GridRowTag{{Cells: []string{"visible"}}}}},
					},
					model.IfTag{
						Cond: "hide",
						Then: []any{model.GridTag{Rows: []model.GridRowTag{{Cells: []string{"hidden"}}}}},
					},
				},
			},
		},
	}

	data := map[string]any{
		"show": true,
		"hide": false,
	}

	book, err := r.Render(context.Background(), gxl, data)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	sheet := book.Sheets[0]

	// Should have "visible" but not "hidden"
	hasVisible := false
	hasHidden := false
	for _, cell := range sheet.Cells {
		if cell.Value == "visible" {
			hasVisible = true
		}
		if cell.Value == "hidden" {
			hasHidden = true
		}
	}

	if !hasVisible {
		t.Error("expected 'visible' cell")
	}
	if hasHidden {
		t.Error("unexpected 'hidden' cell")
	}
}

func TestRenderSheet_EdgeCases(t *testing.T) {
	// Test edge cases: empty grids, invalid refs, etc.
	conf := config.NewBaseConfig()
	r := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "EdgeCases",
				Nodes: []any{
					// Empty grid
					model.GridTag{Rows: []model.GridRowTag{}},
					// Grid with empty rows
					model.GridTag{Rows: []model.GridRowTag{{Cells: []string{}}}},
					// Grid with empty cells
					model.GridTag{Rows: []model.GridRowTag{{Cells: []string{"", "", ""}}}},
					// Invalid ref (should fall back to sequential)
					model.GridTag{Ref: "INVALID", Rows: []model.GridRowTag{{Cells: []string{"test"}}}},
					// Merge with valid range
					model.MergeTag{Range: "A1:B2"},
				},
			},
		},
	}

	book, err := r.Render(context.Background(), gxl, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if len(book.Sheets) != 1 {
		t.Fatalf("expected 1 sheet, got %d", len(book.Sheets))
	}

	sheet := book.Sheets[0]

	// Should have at least the "test" cell
	hasTest := false
	for _, cell := range sheet.Cells {
		if cell.Value == "test" {
			hasTest = true
			break
		}
	}
	if !hasTest {
		t.Error("expected 'test' cell")
	}

	// Should have merge
	if len(sheet.Merges) != 1 {
		t.Errorf("expected 1 merge, got %d", len(sheet.Merges))
	}
}

func TestRenderSheet_ParseA1RefVariations(t *testing.T) {
	// Test parseA1Ref with various formats (currently 35.2% coverage)
	conf := config.NewBaseConfig()
	r := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "A1Refs",
				Nodes: []any{
					// Various valid A1 references
					model.GridTag{Ref: "A1", Rows: []model.GridRowTag{{Cells: []string{"A1"}}}},
					model.GridTag{Ref: "Z10", Rows: []model.GridRowTag{{Cells: []string{"Z10"}}}},
					model.GridTag{Ref: "AA100", Rows: []model.GridRowTag{{Cells: []string{"AA100"}}}},
					model.GridTag{Ref: "AB999", Rows: []model.GridRowTag{{Cells: []string{"AB999"}}}},
					// Invalid refs should fall back to sequential
					model.GridTag{Ref: "123", Rows: []model.GridRowTag{{Cells: []string{"invalid1"}}}},
					model.GridTag{Ref: "", Rows: []model.GridRowTag{{Cells: []string{"empty"}}}},
				},
			},
		},
	}

	book, err := r.Render(context.Background(), gxl, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	sheet := book.Sheets[0]
	if len(sheet.Cells) == 0 {
		t.Fatal("expected cells")
	}

	// Verify specific refs exist
	expectedRefs := map[string]bool{
		"A1":    false,
		"Z10":   false,
		"AA100": false,
		"AB999": false,
	}

	for _, cell := range sheet.Cells {
		if _, ok := expectedRefs[cell.Ref]; ok {
			expectedRefs[cell.Ref] = true
		}
	}

	for ref, found := range expectedRefs {
		if !found {
			t.Errorf("expected cell at ref %s", ref)
		}
	}
}
