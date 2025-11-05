package usecase_test

import (
	"testing"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	usecase "github.com/ryo-arima/goxcel/pkg/usecase"
)

// TestExpandMustache tests the ExpandMustache wrapper function (currently 0% coverage)
func TestExpandMustache_BasicInterpolation(t *testing.T) {
	conf := config.NewBaseConfig()
	bu := usecase.NewBookUsecase(conf)

	gxl := &model.GXL{
		Sheets: []model.SheetTag{
			{
				Name: "Test",
				Nodes: []any{
					model.GridTag{
						Rows: []model.GridRowTag{
							{Cells: []string{"{{ name }}", "{{ age }}", "{{ active }}"}},
						},
					},
				},
			},
		},
	}

	data := map[string]any{
		"name":   "Alice",
		"age":    30,
		"active": true,
	}

	book, err := bu.Render(nil, gxl, data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if len(book.Sheets) != 1 {
		t.Fatalf("Expected 1 sheet, got %d", len(book.Sheets))
	}

	sheet := book.Sheets[0]
	if len(sheet.Cells) != 3 {
		t.Fatalf("Expected 3 cells, got %d", len(sheet.Cells))
	}

	expectedValues := map[string]string{
		"A1": "Alice",
		"B1": "30",
		"C1": "true",
	}

	for _, cell := range sheet.Cells {
		expected, ok := expectedValues[cell.Ref]
		if !ok {
			t.Errorf("Unexpected cell ref: %s", cell.Ref)
			continue
		}
		if cell.Value != expected {
			t.Errorf("Cell %s: expected value %q, got %q", cell.Ref, expected, cell.Value)
		}
	}
}
