package model_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/goxcel/pkg/model"
)

func TestNewBookAndSheet(t *testing.T) {
	b := model.NewBook()
	if len(b.Sheets) != 0 {
		t.Fatalf("new book sheets = %d, want 0", len(b.Sheets))
	}
	s := model.NewSheet("S1")
	if s.Name != "S1" {
		t.Fatalf("sheet name = %q, want S1", s.Name)
	}
	if s.Config == nil {
		t.Fatal("sheet config is nil")
	}
	wantCfg := &model.SheetConfig{
		DefaultRowHeight:   15.0,
		DefaultColumnWidth: 8.43,
		ShowGridLines:      true,
		ShowRowColHeaders:  true,
	}
	if diff := cmp.Diff(wantCfg, s.Config); diff != "" {
		t.Fatalf("sheet config mismatch (-want +got):\n%s", diff)
	}
	b.AddSheet(s)
	if len(b.Sheets) != 1 {
		t.Fatalf("book sheets = %d, want 1", len(b.Sheets))
	}
}
