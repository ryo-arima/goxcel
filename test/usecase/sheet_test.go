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
