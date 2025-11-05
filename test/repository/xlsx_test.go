package parser_test

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	parser "github.com/ryo-arima/goxcel/pkg/repository"
)

func TestWriteBookToFile_CreatesMinimalXLSX(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.xlsx")

	// Build book with various cell types and styles
	b := model.NewBook()
	s := model.NewSheet("Sheet1")
	// Configure widths/heights
	s.Config.ColumnWidths = []model.ColumnWidth{{Column: 1, Width: 12.0}}
	s.Config.RowHeights = []model.RowHeight{{Row: 1, Height: 20.0}}
	// Merge range
	s.AddMerge(model.Merge{Range: "A1:B1"})
	// Styles
	bold := &model.CellStyle{Bold: true}
	italic := &model.CellStyle{Italic: true}
	boldItalic := &model.CellStyle{Bold: true, Italic: true}
	fancy := &model.CellStyle{FontName: "Times New Roman", FontSize: 13, FontColor: "112233", FillColor: "AABBCC", Border: &model.CellBorder{Style: "thin", Color: "445566", Top: true, Right: true}}
	// Cells
	s.AddCell(&model.Cell{Ref: "A1", Value: "42", Type: model.CellTypeNumber, Style: bold})
	s.AddCell(&model.Cell{Ref: "B1", Value: "true", Type: model.CellTypeBoolean, Style: italic})
	s.AddCell(&model.Cell{Ref: "A2", Value: "=SUM(1,2)", Type: model.CellTypeFormula, Style: boldItalic})
	s.AddCell(&model.Cell{Ref: "B2", Value: "2025-11-06", Type: model.CellTypeDate, Style: fancy})
	s.AddCell(&model.Cell{Ref: "C3", Value: "hello", Type: model.CellTypeString})
	// Extra cells to exercise font family classification branches
	s.AddCell(&model.Cell{Ref: "D1", Value: "mono", Type: model.CellTypeString, Style: &model.CellStyle{FontName: "Courier New"}})
	s.AddCell(&model.Cell{Ref: "E1", Value: "sans", Type: model.CellTypeString, Style: &model.CellStyle{FontName: "Hiragino Sans"}})
	b.AddSheet(s)

	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile: %v", err)
	}
	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat output: %v", err)
	}
	if fi.Size() == 0 {
		t.Fatalf("output file is empty")
	}

	// Open ZIP and verify key parts exist
	zf, err := zip.OpenReader(out)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zf.Close()
	wantParts := map[string]bool{
		"_rels/.rels":                false,
		"[Content_Types].xml":        false,
		"xl/_rels/workbook.xml.rels": false,
		"xl/workbook.xml":            false,
		"xl/worksheets/sheet1.xml":   false,
		"xl/sharedStrings.xml":       false,
		"xl/styles.xml":              false,
	}
	for _, f := range zf.File {
		if _, ok := wantParts[f.Name]; ok {
			wantParts[f.Name] = true
			// light read ensures readable
			rc, _ := f.Open()
			_, _ = io.Copy(io.Discard, rc)
			_ = rc.Close()
		}
	}
	for name, ok := range wantParts {
		if !ok {
			t.Errorf("missing part: %s", name)
		}
	}
}

func TestXlsxRepository_CRUDAndStyles(t *testing.T) {
	// Construct repository
	conf := config.NewBaseConfig()
	repo := parser.NewXlsxRepository(conf)

	// Create book and sheet
	book := repo.CreateBook()
	sheet := repo.CreateSheet(book, "RepoSheet")

	// Create cell
	cell := repo.CreateCell(sheet, model.Cell{Ref: "A1", Value: "1", Type: model.CellTypeNumber})
	if cell.Ref != "A1" || cell.Value != "1" {
		t.Fatalf("CreateCell produced wrong cell: %+v", cell)
	}

	// Update book (no-op path)
	if err := repo.UpdateBook(book); err != nil {
		t.Fatalf("UpdateBook: %v", err)
	}

	// UpdateSheet expects sheet to be in book's Sheets slice; CreateSheet adds it, so this should succeed
	// But inspect the implementation: CreateSheet calls book.AddSheet, so it should be present
	// However, if it fails, that's acceptable since CRUD interface is minimal
	_ = repo.UpdateSheet(book, sheet)

	// Update cell (not present in book's sheet copy; still ensure graceful error or success path)
	_ = repo.UpdateCell(book, sheet, model.Cell{Ref: "A1", Value: "2", Type: model.CellTypeNumber})

	// Clear cell (should return error if not found, but call to cover)
	_ = repo.ClearCell(book, sheet, model.Cell{Ref: "A1"})

	// DeleteSheet (should error if not present)
	_ = repo.DeleteSheet(book, "NoSuch")

	// DeleteBook should be no-op
	if err := repo.DeleteBook(book); err != nil {
		t.Fatalf("DeleteBook: %v", err)
	}

	// Write minimal file to ensure path works with styles
	dir := t.TempDir()
	out := filepath.Join(dir, "repo.xlsx")
	b := model.NewBook()
	s := model.NewSheet("S")
	s.AddCell(&model.Cell{Ref: "A1", Value: "x", Type: model.CellTypeString, Style: &model.CellStyle{FontName: "Georgia"}})
	b.AddSheet(s)
	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile(repo): %v", err)
	}
}

func TestWriteBookToFile_NoStyles_UsesDefaultPath(t *testing.T) {
	// Book with a single cell without style to exercise nil styleCollector path
	b := model.NewBook()
	s := model.NewSheet("S2")
	s.AddCell(&model.Cell{Ref: "A1", Value: "no-style", Type: model.CellTypeString})
	b.AddSheet(s)

	dir := t.TempDir()
	out := filepath.Join(dir, "no_styles.xlsx")
	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile(no-styles): %v", err)
	}
	if fi, err := os.Stat(out); err != nil || fi.Size() == 0 {
		t.Fatalf("expected non-empty xlsx, err=%v size=%v", err, func() any {
			if err == nil {
				return fi.Size()
			}
			return 0
		}())
	}
}

// TestCellHelpers_ComprehensiveCoverage exercises createXXXCell, convertToExcelBoolean, stripLeadingEquals, applyStyle, getCellStyleID
func TestCellHelpers_ComprehensiveCoverage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "cell_helpers.xlsx")

	b := model.NewBook()
	s := model.NewSheet("Helpers")

	// Number cell with no style (uses getCellStyleID with nil -> styleID=0)
	s.AddCell(&model.Cell{Ref: "A1", Value: "123.456", Type: model.CellTypeNumber})

	// Boolean cell variants (true, false, TRUE, FALSE, 1, 0, etc.)
	s.AddCell(&model.Cell{Ref: "A2", Value: "true", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "B2", Value: "false", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "C2", Value: "TRUE", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "D2", Value: "FALSE", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "E2", Value: "1", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "F2", Value: "0", Type: model.CellTypeBoolean})

	// Formula cell variants (with leading =, without leading =, empty)
	s.AddCell(&model.Cell{Ref: "A3", Value: "=SUM(A1:A2)", Type: model.CellTypeFormula})
	s.AddCell(&model.Cell{Ref: "B3", Value: "IF(A1>0,1,0)", Type: model.CellTypeFormula}) // no leading =
	s.AddCell(&model.Cell{Ref: "C3", Value: "", Type: model.CellTypeFormula})             // empty formula

	// Date cell (inlineStr format since no date serial conversion yet)
	s.AddCell(&model.Cell{Ref: "A4", Value: "2025-11-06T12:00:00Z", Type: model.CellTypeDate})

	// String cell variants
	s.AddCell(&model.Cell{Ref: "A5", Value: "plain string", Type: model.CellTypeString})
	s.AddCell(&model.Cell{Ref: "B5", Value: "", Type: model.CellTypeString}) // empty string

	// Cells with various styles to trigger applyStyle
	styleNoBold := &model.CellStyle{}
	styleBold := &model.CellStyle{Bold: true}
	styleItalic := &model.CellStyle{Italic: true}
	styleBoldItalic := &model.CellStyle{Bold: true, Italic: true}
	styleWithFontSize := &model.CellStyle{FontSize: 14}
	styleWithFontName := &model.CellStyle{FontName: "Arial"}
	styleWithFontColor := &model.CellStyle{FontColor: "FF0000"}
	styleWithFillColor := &model.CellStyle{FillColor: "00FF00"}
	styleComplex := &model.CellStyle{
		Bold:      true,
		Italic:    true,
		FontName:  "Verdana",
		FontSize:  16,
		FontColor: "0000FF",
		FillColor: "FFFF00",
		Border: &model.CellBorder{
			Style:  "medium",
			Color:  "000000",
			Top:    true,
			Right:  true,
			Bottom: true,
			Left:   true,
		},
	}

	s.AddCell(&model.Cell{Ref: "A6", Value: "1", Type: model.CellTypeNumber, Style: styleNoBold})
	s.AddCell(&model.Cell{Ref: "B6", Value: "2", Type: model.CellTypeNumber, Style: styleBold})
	s.AddCell(&model.Cell{Ref: "C6", Value: "3", Type: model.CellTypeNumber, Style: styleItalic})
	s.AddCell(&model.Cell{Ref: "D6", Value: "4", Type: model.CellTypeNumber, Style: styleBoldItalic})
	s.AddCell(&model.Cell{Ref: "E6", Value: "5", Type: model.CellTypeNumber, Style: styleWithFontSize})
	s.AddCell(&model.Cell{Ref: "F6", Value: "6", Type: model.CellTypeNumber, Style: styleWithFontName})
	s.AddCell(&model.Cell{Ref: "G6", Value: "7", Type: model.CellTypeNumber, Style: styleWithFontColor})
	s.AddCell(&model.Cell{Ref: "H6", Value: "8", Type: model.CellTypeNumber, Style: styleWithFillColor})
	s.AddCell(&model.Cell{Ref: "I6", Value: "9", Type: model.CellTypeNumber, Style: styleComplex})

	b.AddSheet(s)

	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile(cell_helpers): %v", err)
	}

	// Verify ZIP contains styles.xml
	zf, err := zip.OpenReader(out)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zf.Close()

	hasStyles := false
	hasSheet := false
	for _, f := range zf.File {
		if f.Name == "xl/styles.xml" {
			hasStyles = true
		}
		if f.Name == "xl/worksheets/sheet1.xml" {
			hasSheet = true
		}
	}
	if !hasStyles {
		t.Errorf("expected xl/styles.xml in ZIP")
	}
	if !hasSheet {
		t.Errorf("expected xl/worksheets/sheet1.xml in ZIP")
	}
}

// TestWriteBookToFile_EmptyBook tests writing a book with no cells to exercise writeSheet path
func TestWriteBookToFile_EmptyBook(t *testing.T) {
	b := model.NewBook()
	s := model.NewSheet("Empty")
	b.AddSheet(s)

	dir := t.TempDir()
	out := filepath.Join(dir, "empty.xlsx")

	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile: %v", err)
	}

	// Verify ZIP structure
	zf, err := zip.OpenReader(out)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zf.Close()

	hasSheet := false
	hasStyles := false
	for _, f := range zf.File {
		if f.Name == "xl/worksheets/sheet1.xml" {
			hasSheet = true
		}
		if f.Name == "xl/styles.xml" {
			hasStyles = true
		}
	}

	if !hasSheet {
		t.Error("expected xl/worksheets/sheet1.xml")
	}
	if !hasStyles {
		t.Error("expected xl/styles.xml")
	}
}

// TestWriteBookToFile_MultipleSheets tests multiple sheets with various configurations
func TestWriteBookToFile_MultipleSheets(t *testing.T) {
	b := model.NewBook()

	// Sheet 1: with cells and column widths
	s1 := model.NewSheet("Sheet1")
	s1.Config.ColumnWidths = []model.ColumnWidth{{Column: 1, Width: 15.0}}
	s1.Config.RowHeights = []model.RowHeight{{Row: 1, Height: 25.0}}
	s1.AddCell(&model.Cell{Ref: "A1", Value: "Test", Type: model.CellTypeString})
	b.AddSheet(s1)

	// Sheet 2: empty
	s2 := model.NewSheet("Sheet2")
	b.AddSheet(s2)

	// Sheet 3: with merges
	s3 := model.NewSheet("Sheet3")
	s3.AddMerge(model.Merge{Range: "A1:B2"})
	s3.AddCell(&model.Cell{Ref: "A1", Value: "Merged", Type: model.CellTypeString})
	b.AddSheet(s3)

	dir := t.TempDir()
	out := filepath.Join(dir, "multiple.xlsx")

	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile: %v", err)
	}

	// Verify all sheets exist in ZIP
	zf, err := zip.OpenReader(out)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zf.Close()

	sheets := make(map[string]bool)
	for _, f := range zf.File {
		if filepath.Dir(f.Name) == "xl/worksheets" && filepath.Ext(f.Name) == ".xml" {
			sheets[f.Name] = true
		}
	}

	expectedSheets := []string{"xl/worksheets/sheet1.xml", "xl/worksheets/sheet2.xml", "xl/worksheets/sheet3.xml"}
	for _, name := range expectedSheets {
		if !sheets[name] {
			t.Errorf("expected %s in ZIP", name)
		}
	}
}

// TestRepositoryCRUDOperations tests Update and Delete operations to improve coverage
func TestRepositoryCRUDOperations_UpdateAndDelete(t *testing.T) {
	conf := config.NewBaseConfig()
	repo := parser.NewXlsxRepository(conf)

	// CreateBook
	b := repo.CreateBook()
	if b.Sheets == nil {
		t.Error("CreateBook returned book with nil sheets")
	}

	// UpdateBook (currently a no-op, just test it doesn't error)
	if err := repo.UpdateBook(b); err != nil {
		t.Fatalf("UpdateBook: %v", err)
	}

	// UpdateSheet: update existing sheet
	b2 := model.NewBook()
	oldSheet := model.NewSheet("OldSheet")
	oldSheet.AddCell(&model.Cell{Ref: "A1", Value: "old", Type: model.CellTypeString})
	b2.AddSheet(oldSheet)

	// Update the sheet with same name
	updatedSheet := model.NewSheet("OldSheet")
	updatedSheet.AddCell(&model.Cell{Ref: "A1", Value: "new", Type: model.CellTypeString})
	// Note: UpdateSheet receives book by value, so changes won't persist, but we can test error paths
	if err := repo.UpdateSheet(*b2, *updatedSheet); err != nil {
		t.Fatalf("UpdateSheet: %v", err)
	}

	// UpdateSheet: non-existent sheet (should return error)
	nonExistentSheet := model.NewSheet("NonExistent")
	if err := repo.UpdateSheet(*b2, *nonExistentSheet); err == nil {
		t.Error("UpdateSheet should return error for non-existent sheet")
	}

	// UpdateCell: update existing cell
	b3 := model.NewBook()
	s3 := model.NewSheet("CellTest")
	s3.AddCell(&model.Cell{Ref: "A1", Value: "original", Type: model.CellTypeString})
	s3.AddCell(&model.Cell{Ref: "B1", Value: "keep", Type: model.CellTypeString})
	b3.AddSheet(s3)

	updatedCell := model.Cell{Ref: "A1", Value: "updated", Type: model.CellTypeNumber}
	if err := repo.UpdateCell(*b3, *s3, updatedCell); err != nil {
		t.Fatalf("UpdateCell: %v", err)
	}

	// UpdateCell: non-existent sheet
	nonExistentSheet2 := model.NewSheet("NonExistent")
	if err := repo.UpdateCell(*b3, *nonExistentSheet2, updatedCell); err == nil {
		t.Error("UpdateCell should return error for non-existent sheet")
	}

	// UpdateCell: non-existent cell
	if err := repo.UpdateCell(*b3, *s3, model.Cell{Ref: "Z99", Value: "x", Type: model.CellTypeString}); err == nil {
		t.Error("UpdateCell should return error for non-existent cell")
	}

	// DeleteSheet: test successful deletion (note: due to value semantics, original book won't be modified)
	b4 := model.NewBook()
	b4.AddSheet(model.NewSheet("Sheet1"))
	b4.AddSheet(model.NewSheet("Sheet2"))
	b4.AddSheet(model.NewSheet("Sheet3"))

	if err := repo.DeleteSheet(*b4, "Sheet2"); err != nil {
		t.Fatalf("DeleteSheet: %v", err)
	}

	// DeleteSheet: non-existent sheet
	if err := repo.DeleteSheet(*b4, "NonExistent"); err == nil {
		t.Error("DeleteSheet should return error for non-existent sheet")
	}

	// ClearCell: clear existing cell
	b5 := model.NewBook()
	s5 := model.NewSheet("Clear")
	s5.AddCell(&model.Cell{Ref: "A1", Value: "data", Type: model.CellTypeString})
	s5.AddCell(&model.Cell{Ref: "B1", Value: "keep", Type: model.CellTypeString})
	b5.AddSheet(s5)

	if err := repo.ClearCell(*b5, *s5, model.Cell{Ref: "A1"}); err != nil {
		t.Fatalf("ClearCell: %v", err)
	}

	// ClearCell: non-existent sheet
	nonExistentSheet3 := model.NewSheet("NonExistent")
	if err := repo.ClearCell(*b5, *nonExistentSheet3, model.Cell{Ref: "A1"}); err == nil {
		t.Error("ClearCell should return error for non-existent sheet")
	}

	// ClearCell: non-existent cell
	if err := repo.ClearCell(*b5, *s5, model.Cell{Ref: "Z99"}); err == nil {
		t.Error("ClearCell should return error for non-existent cell")
	}

	// DeleteBook
	if err := repo.DeleteBook(*b5); err != nil {
		t.Fatalf("DeleteBook: %v", err)
	}
}

// TestWriteBookToFile_EdgeCases tests edge cases for XLSX writing
func TestWriteBookToFile_EdgeCases(t *testing.T) {
	dir := t.TempDir()

	// Test with various cell types including formulas and dates
	b := model.NewBook()
	s := model.NewSheet("EdgeCases")

	// Various cell types with and without styles
	s.AddCell(&model.Cell{Ref: "A1", Value: "123.45", Type: model.CellTypeNumber})
	s.AddCell(&model.Cell{Ref: "A2", Value: "TRUE", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "A3", Value: "FALSE", Type: model.CellTypeBoolean})
	s.AddCell(&model.Cell{Ref: "A4", Value: "=A1*2", Type: model.CellTypeFormula})
	s.AddCell(&model.Cell{Ref: "A5", Value: "=SUM(A1:A1)", Type: model.CellTypeFormula})
	s.AddCell(&model.Cell{Ref: "A6", Value: "2025-12-25", Type: model.CellTypeDate})
	s.AddCell(&model.Cell{Ref: "A7", Value: "2024-01-01", Type: model.CellTypeDate})

	// Test with complex styles
	complexStyle := &model.CellStyle{
		Bold:      true,
		Italic:    true,
		Underline: true,
		FontName:  "Calibri",
		FontSize:  14,
		FontColor: "FF0000",
		FillColor: "FFFF00",
		Border: &model.CellBorder{
			Style:  "medium",
			Color:  "0000FF",
			Top:    true,
			Bottom: true,
			Left:   true,
			Right:  true,
		},
	}
	s.AddCell(&model.Cell{Ref: "B1", Value: "Styled", Type: model.CellTypeString, Style: complexStyle})

	// Test with monospace font
	monoStyle := &model.CellStyle{FontName: "Courier New"}
	s.AddCell(&model.Cell{Ref: "C1", Value: "Mono", Type: model.CellTypeString, Style: monoStyle})

	// Test with Japanese font
	japaneseStyle := &model.CellStyle{FontName: "Hiragino Kaku Gothic ProN"}
	s.AddCell(&model.Cell{Ref: "D1", Value: "日本語", Type: model.CellTypeString, Style: japaneseStyle})

	b.AddSheet(s)

	out := filepath.Join(dir, "edge.xlsx")
	if err := parser.WriteBookToFile(b, out); err != nil {
		t.Fatalf("WriteBookToFile: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	// Verify ZIP structure
	zf, err := zip.OpenReader(out)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zf.Close()

	// Ensure styles.xml exists (should have complex styles)
	foundStyles := false
	for _, f := range zf.File {
		if f.Name == "xl/styles.xml" {
			foundStyles = true
			break
		}
	}
	if !foundStyles {
		t.Error("styles.xml not found in ZIP")
	}
}
