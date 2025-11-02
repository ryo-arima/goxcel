package parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
)

// This file is intentionally left minimal. The renderer has moved to pkg/usecase.
// Keeping this file ensures the package directory compiles without mixed packages.

// XlsxRepository manages Excel workbook operations.
type XlsxRepository interface {
	CreateBook() model.Book
	CreateSheet(book model.Book, name string) model.Sheet
	CreateCell(sheet model.Sheet, cell model.Cell) model.Cell
	UpdateBook(book model.Book) error
	UpdateSheet(book model.Book, sheet model.Sheet) error
	UpdateCell(book model.Book, sheet model.Sheet, cell model.Cell) error
	DeleteBook(book model.Book) error
	DeleteSheet(book model.Book, sheetName string) error
	ClearCell(book model.Book, sheet model.Sheet, cell model.Cell) error
}

type xlsxRepository struct {
	Conf config.BaseConfig
}

// NewXlsxRepository creates a new XLSX repository.
func NewXlsxRepository(conf config.BaseConfig) XlsxRepository {
	return &xlsxRepository{
		Conf: conf,
	}
}

// CreateBook creates a new empty workbook.
func (r *xlsxRepository) CreateBook() model.Book {
	return *model.NewBook()
}

// CreateSheet creates a new sheet in the workbook.
func (r *xlsxRepository) CreateSheet(book model.Book, name string) model.Sheet {
	sheet := model.NewSheet(name)
	book.AddSheet(sheet)
	return *sheet
}

// CreateCell creates a new cell in the sheet.
func (r *xlsxRepository) CreateCell(sheet model.Sheet, cell model.Cell) model.Cell {
	newCell := &model.Cell{
		Ref:   cell.Ref,
		Value: cell.Value,
		Type:  cell.Type,
		Style: cell.Style,
	}
	sheet.AddCell(newCell)
	return *newCell
}

// UpdateBook updates book-level properties.
func (r *xlsxRepository) UpdateBook(book model.Book) error {
	// TODO: Implement book-level updates (metadata, properties, etc.)
	return nil
}

// UpdateSheet updates sheet-level properties.
func (r *xlsxRepository) UpdateSheet(book model.Book, sheet model.Sheet) error {
	// Find and update the sheet in the book
	for i, s := range book.Sheets {
		if s.Name == sheet.Name {
			book.Sheets[i] = &sheet
			return nil
		}
	}
	return fmt.Errorf("sheet %q not found in book", sheet.Name)
}

// UpdateCell updates an existing cell's value.
func (r *xlsxRepository) UpdateCell(book model.Book, sheet model.Sheet, cell model.Cell) error {
	// Find the sheet in the book
	var targetSheet *model.Sheet
	for _, s := range book.Sheets {
		if s.Name == sheet.Name {
			targetSheet = s
			break
		}
	}
	if targetSheet == nil {
		return fmt.Errorf("sheet %q not found in book", sheet.Name)
	}

	// Find and update the cell in the sheet
	for i, c := range targetSheet.Cells {
		if c.Ref == cell.Ref {
			targetSheet.Cells[i].Value = cell.Value
			targetSheet.Cells[i].Type = cell.Type
			targetSheet.Cells[i].Style = cell.Style
			return nil
		}
	}
	return fmt.Errorf("cell %q not found in sheet %q", cell.Ref, sheet.Name)
}

// DeleteBook removes a book from memory (cleanup).
func (r *xlsxRepository) DeleteBook(book model.Book) error {
	// Clear all sheets
	book.Sheets = nil
	return nil
}

// DeleteSheet removes a sheet from the workbook by name.
func (r *xlsxRepository) DeleteSheet(book model.Book, sheetName string) error {
	var newSheets []*model.Sheet
	found := false
	for _, sheet := range book.Sheets {
		if sheet.Name != sheetName {
			newSheets = append(newSheets, sheet)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("sheet %q not found", sheetName)
	}

	book.Sheets = newSheets
	return nil
}

// ClearCell clears a cell's value in the sheet.
func (r *xlsxRepository) ClearCell(book model.Book, sheet model.Sheet, cell model.Cell) error {
	// Find the sheet in the book
	var targetSheet *model.Sheet
	for _, s := range book.Sheets {
		if s.Name == sheet.Name {
			targetSheet = s
			break
		}
	}
	if targetSheet == nil {
		return fmt.Errorf("sheet %q not found in book", sheet.Name)
	}

	// Find and clear the cell in the sheet
	for i, c := range targetSheet.Cells {
		if c.Ref == cell.Ref {
			targetSheet.Cells[i].Value = ""
			return nil
		}
	}
	return fmt.Errorf("cell %q not found in sheet %q", cell.Ref, sheet.Name)
}

// WriteBookToFile writes a Book to an XLSX file.
func WriteBookToFile(book *model.Book, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Write _rels/.rels
	if err := writeRels(zipWriter); err != nil {
		return err
	}

	// Write [Content_Types].xml
	if err := writeContentTypes(zipWriter, len(book.Sheets)); err != nil {
		return err
	}

	// Write xl/_rels/workbook.xml.rels
	if err := writeWorkbookRels(zipWriter, len(book.Sheets)); err != nil {
		return err
	}

	// Write xl/workbook.xml
	if err := writeWorkbook(zipWriter, book); err != nil {
		return err
	}

	// Write xl/worksheets/sheet*.xml
	for i, sheet := range book.Sheets {
		if err := writeSheet(zipWriter, sheet, i+1); err != nil {
			return err
		}
	}

	// Write xl/sharedStrings.xml (empty for now)
	if err := writeSharedStrings(zipWriter); err != nil {
		return err
	}

	// Write xl/styles.xml (minimal)
	if err := writeStyles(zipWriter); err != nil {
		return err
	}

	return nil
}

func writeRels(zw *zip.Writer) error {
	w, err := zw.Create("_rels/.rels")
	if err != nil {
		return err
	}

	rels := model.XMLRelationships{
		Xmlns: model.XMLNsPackageRelationships,
		Relationships: []model.XMLRelationship{
			{
				ID:     "rId1",
				Type:   model.XMLRelTypeOfficeDocument,
				Target: "xl/workbook.xml",
			},
		},
	}

	data, err := xml.MarshalIndent(rels, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeContentTypes(zw *zip.Writer, numSheets int) error {
	w, err := zw.Create("[Content_Types].xml")
	if err != nil {
		return err
	}

	types := model.XMLTypes{
		Xmlns: model.XMLNsPackageContentTypes,
		Defaults: []model.XMLDefault{
			{Extension: "rels", ContentType: "application/vnd.openxmlformats-package.relationships+xml"},
			{Extension: "xml", ContentType: "application/xml"},
		},
		Overrides: []model.XMLOverride{
			{PartName: "/xl/workbook.xml", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"},
			{PartName: "/xl/styles.xml", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"},
			{PartName: "/xl/sharedStrings.xml", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"},
		},
	}

	for i := 1; i <= numSheets; i++ {
		types.Overrides = append(types.Overrides, model.XMLOverride{
			PartName:    fmt.Sprintf("/xl/worksheets/sheet%d.xml", i),
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml",
		})
	}

	data, err := xml.MarshalIndent(types, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeWorkbookRels(zw *zip.Writer, sheetCount int) error {
	w, err := zw.Create("xl/_rels/workbook.xml.rels")
	if err != nil {
		return err
	}

	rels := model.XMLRelationships{
		Xmlns:         model.XMLNsPackageRelationships,
		Relationships: []model.XMLRelationship{},
	}

	for i := 1; i <= sheetCount; i++ {
		rels.Relationships = append(rels.Relationships, model.XMLRelationship{
			ID:     fmt.Sprintf("rId%d", i),
			Type:   model.XMLRelTypeWorksheet,
			Target: fmt.Sprintf("worksheets/sheet%d.xml", i),
		})
	}

	rels.Relationships = append(rels.Relationships,
		model.XMLRelationship{
			ID:     fmt.Sprintf("rId%d", sheetCount+1),
			Type:   model.XMLRelTypeStyles,
			Target: "styles.xml",
		},
		model.XMLRelationship{
			ID:     fmt.Sprintf("rId%d", sheetCount+2),
			Type:   model.XMLRelTypeSharedStrings,
			Target: "sharedStrings.xml",
		},
	)

	data, err := xml.MarshalIndent(rels, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeWorkbook(zw *zip.Writer, book *model.Book) error {
	w, err := zw.Create("xl/workbook.xml")
	if err != nil {
		return err
	}

	workbook := model.XMLWorkbook{
		Xmlns:  model.XMLNsSpreadsheetML,
		XmlnsR: model.XMLNsOfficeDocRelationships,
		Sheets: model.XMLSheets{
			Sheet: []model.XMLSheetRef{},
		},
	}

	for i, sheet := range book.Sheets {
		workbook.Sheets.Sheet = append(workbook.Sheets.Sheet, model.XMLSheetRef{
			Name:    sheet.Name,
			SheetID: i + 1,
			RID:     fmt.Sprintf("rId%d", i+1),
		})
	}

	data, err := xml.MarshalIndent(workbook, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeSheet(zw *zip.Writer, sheet *model.Sheet, sheetNum int) error {
	w, err := zw.Create(fmt.Sprintf("xl/worksheets/sheet%d.xml", sheetNum))
	if err != nil {
		return err
	}

	worksheet := model.XMLWorksheet{
		Xmlns: model.XMLNsSpreadsheetML,
		SheetData: model.XMLSheetData{
			Rows: []model.XMLRow{},
		},
	}

	// Group cells by row
	rowMap := make(map[int][]*model.Cell)
	for _, cell := range sheet.Cells {
		row, _, err := parseA1Ref(cell.Ref)
		if err != nil {
			continue
		}
		rowMap[row] = append(rowMap[row], cell)
	}

	// Write rows in order
	for row := 1; row <= 1000; row++ {
		cells, ok := rowMap[row]
		if !ok {
			continue
		}

		xmlRow := model.XMLRow{
			R:     row,
			Cells: []model.XMLCell{},
		}

		for _, cell := range cells {
			xmlRow.Cells = append(xmlRow.Cells, model.XMLCell{
				R: cell.Ref,
				T: "inlineStr",
				IS: &model.XMLIS{
					T: cell.Value,
				},
			})
		}

		worksheet.SheetData.Rows = append(worksheet.SheetData.Rows, xmlRow)
	}

	// Add merges if any
	if len(sheet.Merges) > 0 {
		mergeCells := &model.XMLMergeCells{
			Count: len(sheet.Merges),
			Merge: []model.XMLMergeCell{},
		}
		for _, merge := range sheet.Merges {
			mergeCells.Merge = append(mergeCells.Merge, model.XMLMergeCell{
				Ref: merge.Range,
			})
		}
		worksheet.MergeCells = mergeCells
	}

	data, err := xml.MarshalIndent(worksheet, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeSharedStrings(zw *zip.Writer) error {
	w, err := zw.Create("xl/sharedStrings.xml")
	if err != nil {
		return err
	}

	sst := model.XMLSharedStrings{
		Xmlns:       model.XMLNsSpreadsheetML,
		Count:       0,
		UniqueCount: 0,
	}

	data, err := xml.MarshalIndent(sst, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeStyles(zw *zip.Writer) error {
	w, err := zw.Create("xl/styles.xml")
	if err != nil {
		return err
	}

	styleSheet := model.XMLStyleSheet{
		Xmlns: model.XMLNsSpreadsheetML,
		Fonts: model.XMLFonts{
			Count: 1,
			Font: []model.XMLFont{
				{
					Sz:   model.XMLFontSize{Val: "11"},
					Name: model.XMLFontName{Val: "Calibri"},
				},
			},
		},
		Fills: model.XMLFills{
			Count: 1,
			Fill: []model.XMLFill{
				{
					PatternFill: model.XMLPatternFill{PatternType: "none"},
				},
			},
		},
		Borders: model.XMLBorders{
			Count: 1,
			Border: []model.XMLBorder{
				{
					Left:   model.XMLBorderSide{},
					Right:  model.XMLBorderSide{},
					Top:    model.XMLBorderSide{},
					Bottom: model.XMLBorderSide{},
				},
			},
		},
		CellXfs: model.XMLCellXfs{
			Count: 1,
			Xf: []model.XMLXf{
				{
					NumFmtID: 0,
					FontID:   0,
					FillID:   0,
					BorderID: 0,
				},
			},
		},
	}

	data, err := xml.MarshalIndent(styleSheet, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func parseA1Ref(ref string) (int, int, error) {
	if ref == "" {
		return 0, 0, fmt.Errorf("empty ref")
	}
	i := 0
	for i < len(ref) && ((ref[i] >= 'A' && ref[i] <= 'Z') || (ref[i] >= 'a' && ref[i] <= 'z')) {
		i++
	}
	if i == 0 || i == len(ref) {
		return 0, 0, fmt.Errorf("invalid ref: %s", ref)
	}
	colLetters := ref[:i]
	rowDigits := ref[i:]
	col := 0
	for _, ch := range colLetters {
		uc := ch
		if uc >= 'a' && uc <= 'z' {
			uc = uc - 'a' + 'A'
		}
		col = col*26 + int(uc-'A'+1)
	}
	row := 0
	for _, ch := range rowDigits {
		if ch < '0' || ch > '9' {
			return 0, 0, fmt.Errorf("invalid row in ref: %s", ref)
		}
		row = row*10 + int(ch-'0')
	}
	return row, col, nil
}
