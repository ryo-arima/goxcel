package parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

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

	// Collect all unique styles from cells
	styleCollector := newStyleCollector()
	for _, sheet := range book.Sheets {
		for _, cell := range sheet.Cells {
			if cell.Style != nil {
				styleCollector.AddStyle(cell.Style)
			}
		}
	}

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
		if err := writeSheetWithStyles(zipWriter, sheet, i+1, styleCollector); err != nil {
			return err
		}
	}

	// Write xl/sharedStrings.xml (empty for now)
	if err := writeSharedStrings(zipWriter); err != nil {
		return err
	}

	// Write xl/styles.xml with collected styles
	if err := writeStylesWithCollector(zipWriter, styleCollector); err != nil {
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

// createXMLCell creates an XMLCell based on the cell type
func createXMLCell(cell *model.Cell) model.XMLCell {
	return createXMLCellWithStyle(cell, nil)
}

func createXMLCellWithStyle(cell *model.Cell, styleCollector *styleCollector) model.XMLCell {
	var styleID int
	if styleCollector != nil {
		styleID = styleCollector.GetStyleID(cell.Style)
	} else {
		styleID = getCellStyleID(cell.Style)
	}

	switch cell.Type {
	case model.CellTypeNumber:
		return createNumberCell(cell, styleID)
	case model.CellTypeBoolean:
		return createBooleanCell(cell, styleID)
	case model.CellTypeFormula:
		return createFormulaCell(cell, styleID)
	case model.CellTypeDate:
		return createDateCell(cell, styleID)
	default:
		return createStringCell(cell, styleID)
	}
}

// createNumberCell creates a numeric cell
func createNumberCell(cell *model.Cell, styleID int) model.XMLCell {
	xmlCell := model.XMLCell{
		R: cell.Ref,
		V: &cell.Value,
	}
	applyStyle(&xmlCell, styleID)
	return xmlCell
}

// createBooleanCell creates a boolean cell
func createBooleanCell(cell *model.Cell, styleID int) model.XMLCell {
	boolValue := convertToExcelBoolean(cell.Value)
	xmlCell := model.XMLCell{
		R: cell.Ref,
		T: "b",
		V: &boolValue,
	}
	applyStyle(&xmlCell, styleID)
	return xmlCell
}

// createFormulaCell creates a formula cell
func createFormulaCell(cell *model.Cell, styleID int) model.XMLCell {
	formulaText := stripLeadingEquals(cell.Value)
	xmlCell := model.XMLCell{
		R: cell.Ref,
		F: &model.XMLFormula{Text: formulaText},
	}
	applyStyle(&xmlCell, styleID)
	return xmlCell
}

// createDateCell creates a date cell
func createDateCell(cell *model.Cell, styleID int) model.XMLCell {
	// TODO: Convert to Excel date serial number and apply date format
	xmlCell := model.XMLCell{
		R:  cell.Ref,
		T:  "inlineStr",
		IS: &model.XMLIS{T: cell.Value},
	}
	applyStyle(&xmlCell, styleID)
	return xmlCell
}

// createStringCell creates a string cell
func createStringCell(cell *model.Cell, styleID int) model.XMLCell {
	xmlCell := model.XMLCell{
		R:  cell.Ref,
		T:  "inlineStr",
		IS: &model.XMLIS{T: cell.Value},
	}
	applyStyle(&xmlCell, styleID)
	return xmlCell
}

// convertToExcelBoolean converts string boolean to Excel format (0 or 1)
func convertToExcelBoolean(value string) string {
	if value == "true" || value == "TRUE" || value == "1" {
		return "1"
	}
	return "0"
}

// stripLeadingEquals removes leading = from formula
func stripLeadingEquals(formula string) string {
	if len(formula) > 0 && formula[0] == '=' {
		return formula[1:]
	}
	return formula
}

// applyStyle applies style ID to cell if non-zero
func applyStyle(xmlCell *model.XMLCell, styleID int) {
	if styleID > 0 {
		xmlCell.S = &styleID
	}
}

// getCellStyleID returns the style ID based on cell formatting
func getCellStyleID(style *model.CellStyle) int {
	if style == nil {
		return 0 // Normal style
	}

	if style.Bold && style.Italic {
		return 3 // Bold + Italic
	} else if style.Bold {
		return 1 // Bold only
	} else if style.Italic {
		return 2 // Italic only
	}

	return 0 // Normal
}

func writeSheet(zw *zip.Writer, sheet *model.Sheet, sheetNum int) error {
	return writeSheetWithStyles(zw, sheet, sheetNum, nil)
}

func writeSheetWithStyles(zw *zip.Writer, sheet *model.Sheet, sheetNum int, styleCollector *styleCollector) error {
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

	// Apply default row height and column width via sheetFormatPr
	if sheet.Config != nil {
		sfp := &model.XMLSheetFormatPr{DefaultRowHeight: sheet.Config.DefaultRowHeight}
		if sheet.Config.DefaultColumnWidth > 0 {
			sfp.DefaultColWidth = sheet.Config.DefaultColumnWidth
		}
		// Optionally set baseColWidth if desired; omit to let Excel infer
		worksheet.SheetFormatPr = sfp
	}

	// Add column widths if configured
	if sheet.Config != nil && len(sheet.Config.ColumnWidths) > 0 {
		cols := &model.XMLCols{
			Col: []model.XMLCol{},
		}
		for _, cw := range sheet.Config.ColumnWidths {
			cols.Col = append(cols.Col, model.XMLCol{
				Min:         cw.Column,
				Max:         cw.Column,
				Width:       cw.Width,
				CustomWidth: true,
			})
		}
		worksheet.Cols = cols
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

		// Set row height if configured
		if sheet.Config != nil {
			for _, rh := range sheet.Config.RowHeights {
				if rh.Row == row {
					xmlRow.Height = rh.Height
					xmlRow.CustomHeight = true
					break
				}
			}
		}

		for _, cell := range cells {
			xmlCell := createXMLCellWithStyle(cell, styleCollector)
			xmlRow.Cells = append(xmlRow.Cells, xmlCell)
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
			Count: 4,
			Font: []model.XMLFont{
				// Font 0: Normal
				{
					Sz:   model.XMLFontSize{Val: "11"},
					Name: model.XMLFontName{Val: "Calibri"},
				},
				// Font 1: Bold
				{
					Sz:   model.XMLFontSize{Val: "11"},
					Name: model.XMLFontName{Val: "Calibri"},
					B:    &model.XMLBold{},
				},
				// Font 2: Italic
				{
					Sz:   model.XMLFontSize{Val: "11"},
					Name: model.XMLFontName{Val: "Calibri"},
					I:    &model.XMLItalic{},
				},
				// Font 3: Bold + Italic
				{
					Sz:   model.XMLFontSize{Val: "11"},
					Name: model.XMLFontName{Val: "Calibri"},
					B:    &model.XMLBold{},
					I:    &model.XMLItalic{},
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
			Count: 4,
			Xf: []model.XMLXf{
				// Style 0: Normal
				{
					NumFmtID: 0,
					FontID:   0,
					FillID:   0,
					BorderID: 0,
				},
				// Style 1: Bold
				{
					NumFmtID: 0,
					FontID:   1,
					FillID:   0,
					BorderID: 0,
				},
				// Style 2: Italic
				{
					NumFmtID: 0,
					FontID:   2,
					FillID:   0,
					BorderID: 0,
				},
				// Style 3: Bold + Italic
				{
					NumFmtID: 0,
					FontID:   3,
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

func writeStylesWithCollector(zw *zip.Writer, sc *styleCollector) error {
	w, err := zw.Create("xl/styles.xml")
	if err != nil {
		return err
	}

	// Build fonts, fills, and borders dynamically from collected styles
	fonts := []model.XMLFont{}
	fills := []model.XMLFill{{PatternFill: model.XMLPatternFill{PatternType: "none"}}}
	borders := []model.XMLBorder{{
		Left:   model.XMLBorderSide{},
		Right:  model.XMLBorderSide{},
		Top:    model.XMLBorderSide{},
		Bottom: model.XMLBorderSide{},
	}}
	xfs := []model.XMLXf{}

	for i, style := range sc.styles {
		fontName := "Calibri"
		fontSize := "11"

		if style != nil {
			if style.FontName != "" {
				fontName = style.FontName
			}
			if style.FontSize > 0 {
				fontSize = fmt.Sprintf("%d", style.FontSize)
			}
		}

		font := model.XMLFont{
			Sz:   model.XMLFontSize{Val: fontSize},
			Name: model.XMLFontName{Val: fontName},
		}

		if style != nil {
			if style.Bold {
				font.B = &model.XMLBold{}
			}
			if style.Italic {
				font.I = &model.XMLItalic{}
			}
			if style.Underline {
				font.U = &model.XMLUnderline{}
			}
			if style.FontColor != "" {
				font.Color = &model.XMLFontColor{RGB: "FF" + style.FontColor}
			}
			// Add family/charset hints for better compatibility (e.g., LibreOffice)
			fam := classifyFontFamily(fontName)
			if fam > 0 {
				font.Family = &model.XMLFontFamily{Val: fam}
			}
			font.Charset = &model.XMLFontCharset{Val: 0}
		}

		fonts = append(fonts, font)

		// Create fill for background color
		fillID := 0
		if style != nil && style.FillColor != "" {
			fill := model.XMLFill{
				PatternFill: model.XMLPatternFill{
					PatternType: "solid",
					FgColor:     &model.XMLFillColor{RGB: "FF" + style.FillColor},
					BgColor:     &model.XMLBgColor{Indexed: 64},
				},
			}
			fills = append(fills, fill)
			fillID = len(fills) - 1
		}

		// Border
		borderID := 0
		if style != nil && style.Border != nil {
			mkSide := func(on bool) model.XMLBorderSide {
				if !on || style.Border.Style == "" {
					return model.XMLBorderSide{}
				}
				var color *model.XMLBorderColor
				if style.Border.Color != "" {
					color = &model.XMLBorderColor{RGB: "FF" + style.Border.Color}
				}
				return model.XMLBorderSide{Style: style.Border.Style, Color: color}
			}
			b := model.XMLBorder{
				Left:   mkSide(style.Border.Left),
				Right:  mkSide(style.Border.Right),
				Top:    mkSide(style.Border.Top),
				Bottom: mkSide(style.Border.Bottom),
			}
			borders = append(borders, b)
			borderID = len(borders) - 1
		}

		xf := model.XMLXf{
			NumFmtID:    0,
			FontID:      i,
			FillID:      fillID,
			BorderID:    borderID,
			ApplyFont:   true,
			ApplyFill:   fillID != 0,
			ApplyBorder: borderID != 0,
		}
		xfs = append(xfs, xf)
	}

	styleSheet := model.XMLStyleSheet{
		Xmlns: model.XMLNsSpreadsheetML,
		Fonts: model.XMLFonts{
			Count: len(fonts),
			Font:  fonts,
		},
		Fills: model.XMLFills{
			Count: len(fills),
			Fill:  fills,
		},
		Borders: model.XMLBorders{
			Count:  len(borders),
			Border: borders,
		},
		CellStyleXfs: model.XMLCellStyleXfs{
			Count: 1,
			Xf: []model.XMLXf{
				{NumFmtID: 0, FontID: 0, FillID: 0, BorderID: 0}, // base Normal
			},
		},
		CellXfs: model.XMLCellXfs{
			Count: len(xfs),
			Xf:    xfs,
		},
		CellStyles: model.XMLCellStyles{
			Count: 1,
			Cell: []model.XMLCellStyle{
				{Name: "Normal", XfID: 0, BuiltinID: 0},
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

// classifyFontFamily maps a font name to OOXML family code
// 0: unknown, 1: Roman (serif), 2: Swiss (sans-serif), 3: Modern (monospace), 4: Script, 5: Decorative
func classifyFontFamily(name string) int {
	n := strings.ToLower(name)
	// Sans-serif
	if strings.Contains(n, "arial") || strings.Contains(n, "calibri") || strings.Contains(n, "helvetica") || strings.Contains(n, "hiragino") || strings.Contains(n, "noto sans") || strings.Contains(n, "source sans") || strings.Contains(n, "yu gothic") || strings.Contains(n, "meiryo") || strings.Contains(n, "ms pgothic") || strings.Contains(n, "ms gothic") || strings.Contains(n, "liberation sans") {
		return 2
	}
	// Serif
	if strings.Contains(n, "times") || strings.Contains(n, "georgia") || strings.Contains(n, "noto serif") {
		return 1
	}
	// Monospace
	if strings.Contains(n, "courier") || strings.Contains(n, "consolas") || strings.Contains(n, "menlo") || strings.Contains(n, "monaco") || strings.Contains(n, "source code") {
		return 3
	}
	return 0
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

// styleCollector manages unique styles and generates style IDs
type styleCollector struct {
	styles   []*model.CellStyle
	styleMap map[string]int // style signature -> style ID
}

func newStyleCollector() *styleCollector {
	return &styleCollector{
		styles:   []*model.CellStyle{nil}, // ID 0 is always nil (default style)
		styleMap: make(map[string]int),
	}
}

func (sc *styleCollector) AddStyle(style *model.CellStyle) {
	if style == nil {
		return
	}

	sig := sc.styleSignature(style)
	if _, exists := sc.styleMap[sig]; !exists {
		styleID := len(sc.styles)
		sc.styleMap[sig] = styleID
		sc.styles = append(sc.styles, style)
	}
}

func (sc *styleCollector) GetStyleID(style *model.CellStyle) int {
	if style == nil {
		return 0
	}

	sig := sc.styleSignature(style)
	if id, exists := sc.styleMap[sig]; exists {
		return id
	}
	return 0
}

func (sc *styleCollector) styleSignature(style *model.CellStyle) string {
	// Border signature components
	var bSig string
	if style.Border != nil {
		bSig = fmt.Sprintf("|b:%s|c:%s|t:%t|r:%t|b:%t|l:%t",
			style.Border.Style, style.Border.Color,
			style.Border.Top, style.Border.Right, style.Border.Bottom, style.Border.Left)
	}
	return fmt.Sprintf("%v|%v|%v|%s|%d|%s|%s%s",
		style.Bold, style.Italic, style.Underline,
		style.FontName, style.FontSize,
		style.FontColor, style.FillColor,
		bSig)
}
