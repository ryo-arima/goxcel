package model

// XML namespace constants for Office Open XML format
const (
	// Package namespaces
	XMLNsPackageRelationships = "http://schemas.openxmlformats.org/package/2006/relationships"
	XMLNsPackageContentTypes  = "http://schemas.openxmlformats.org/package/2006/content-types"

	// SpreadsheetML namespace
	XMLNsSpreadsheetML = "http://schemas.openxmlformats.org/spreadsheetml/2006/main"

	// Office document relationships namespace
	XMLNsOfficeDocRelationships = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"

	// Relationship types
	XMLRelTypeOfficeDocument = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
	XMLRelTypeWorksheet      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet"
	XMLRelTypeStyles         = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
	XMLRelTypeSharedStrings  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings"
)

// CellType indicates the kind of cell value; a hint for writer.
type CellType string

const (
	CellTypeString  CellType = "string"  // Text string
	CellTypeNumber  CellType = "number"  // Numeric value
	CellTypeBoolean CellType = "boolean" // Boolean (true/false)
	CellTypeDate    CellType = "date"    // Date/time value
	CellTypeFormula CellType = "formula" // Excel formula
	CellTypeAuto    CellType = "auto"    // Auto-detect type
)

// CellStyle represents the formatting applied to a cell
type CellStyle struct {
	Bold      bool
	Italic    bool
	Underline bool

	// Color and font
	FontName  string // e.g., "Arial", "Calibri"
	FontSize  int    // Font size in points (e.g., 11, 14)
	FontColor string // RGB hex color: "FF0000" for red
	FillColor string // Cell background RGB hex: "FFFF00" for yellow

	// Alignment (future)
	HAlign string // "left", "center", "right"
	VAlign string // "top", "middle", "bottom"

	// Border (optional)
	Border *CellBorder
}

// CellBorder represents per-cell border settings
type CellBorder struct {
	Style  string // thin, medium, thick, dashed, dotted, double
	Color  string // RGB hex without #
	Top    bool
	Right  bool
	Bottom bool
	Left   bool
}

// ColumnWidth represents column width settings
type ColumnWidth struct {
	Column int     // Column number (1-based)
	Width  float64 // Width in Excel units (default ~8.43)
}

// RowHeight represents row height settings
type RowHeight struct {
	Row    int     // Row number (1-based)
	Height float64 // Height in points (default 15)
}

// SheetConfig represents sheet-level configuration
type SheetConfig struct {
	DefaultRowHeight   float64       // Default row height in points
	DefaultColumnWidth float64       // Default column width in Excel units
	ColumnWidths       []ColumnWidth // Specific column widths
	RowHeights         []RowHeight   // Specific row heights
	FreezePane         string        // Cell reference for freeze panes (e.g., "B2")
	ShowGridLines      bool          // Show/hide grid lines
	ShowRowColHeaders  bool          // Show/hide row/column headers
}

// Cell represents a single cell value.
type Cell struct {
	Ref   string     // A1 reference
	Value string     // text value (after {{ }} expansion)
	Type  CellType   // hint for writer
	Style *CellStyle // formatting options
}

// Merge represents a merged cell range like "A1:C1".
type Merge struct {
	Range string
}

// Image represents an image placement on a sheet.
type Image struct {
	Ref      string
	Source   string
	WidthPx  int
	HeightPx int
	Options  map[string]string
}

// Shape represents a drawable shape with optional text.
type Shape struct {
	Ref      string
	Kind     string
	Text     string
	WidthPx  int
	HeightPx int
	Style    string
	Options  map[string]string
}

// Chart represents a chart with data range.
type Chart struct {
	Ref       string
	Type      string
	DataRange string
	Title     string
	WidthPx   int
	HeightPx  int
	Options   map[string]string
}

// PivotTable represents a pivot table definition.
type PivotTable struct {
	Ref         string
	SourceRange string
	Rows        []string
	Columns     []string
	Values      []string
	Filters     []string
	Options     map[string]string
}

// Sheet contains grid cells and drawing objects.
type Sheet struct {
	Name   string
	Cells  []*Cell
	Merges []Merge
	Images []Image
	Shapes []Shape
	Charts []Chart
	Pivots []PivotTable
	Config *SheetConfig // Sheet-level configuration
}

// NewSheet creates a sheet with a given name.
func NewSheet(name string) *Sheet {
	return &Sheet{
		Name: name,
		Config: &SheetConfig{
			DefaultRowHeight:   15.0,
			DefaultColumnWidth: 8.43,
			ShowGridLines:      true,
			ShowRowColHeaders:  true,
		},
	}
}

func (s *Sheet) AddCell(c *Cell)       { s.Cells = append(s.Cells, c) }
func (s *Sheet) AddMerge(m Merge)      { s.Merges = append(s.Merges, m) }
func (s *Sheet) AddImage(i Image)      { s.Images = append(s.Images, i) }
func (s *Sheet) AddShape(sh Shape)     { s.Shapes = append(s.Shapes, sh) }
func (s *Sheet) AddChart(ch Chart)     { s.Charts = append(s.Charts, ch) }
func (s *Sheet) AddPivot(p PivotTable) { s.Pivots = append(s.Pivots, p) }

// Book represents a workbook containing multiple sheets.
type Book struct {
	Sheets []*Sheet
}

// NewBook creates an empty workbook.
func NewBook() *Book { return &Book{Sheets: []*Sheet{}} }

// AddSheet appends a sheet to the workbook.
func (b *Book) AddSheet(s *Sheet) { b.Sheets = append(b.Sheets, s) }

// XML structures for XLSX file format

// XMLRelationships represents the Relationships XML structure
type XMLRelationships struct {
	XMLName       struct{} `xml:"Relationships"`
	Xmlns         string   `xml:"xmlns,attr"`
	Relationships []XMLRelationship
}

// XMLRelationship represents a single Relationship element
type XMLRelationship struct {
	XMLName struct{} `xml:"Relationship"`
	ID      string   `xml:"Id,attr"`
	Type    string   `xml:"Type,attr"`
	Target  string   `xml:"Target,attr"`
}

// XMLTypes represents the [Content_Types].xml structure
type XMLTypes struct {
	XMLName   struct{}      `xml:"Types"`
	Xmlns     string        `xml:"xmlns,attr"`
	Defaults  []XMLDefault  `xml:"Default"`
	Overrides []XMLOverride `xml:"Override"`
}

// XMLDefault represents a Default content type
type XMLDefault struct {
	XMLName     struct{} `xml:"Default"`
	Extension   string   `xml:"Extension,attr"`
	ContentType string   `xml:"ContentType,attr"`
}

// XMLOverride represents an Override content type
type XMLOverride struct {
	XMLName     struct{} `xml:"Override"`
	PartName    string   `xml:"PartName,attr"`
	ContentType string   `xml:"ContentType,attr"`
}

// XMLWorkbook represents the xl/workbook.xml structure
type XMLWorkbook struct {
	XMLName struct{}  `xml:"workbook"`
	Xmlns   string    `xml:"xmlns,attr"`
	XmlnsR  string    `xml:"xmlns:r,attr"`
	Sheets  XMLSheets `xml:"sheets"`
}

// XMLSheets contains sheet references
type XMLSheets struct {
	XMLName struct{}      `xml:"sheets"`
	Sheet   []XMLSheetRef `xml:"sheet"`
}

// XMLSheetRef represents a sheet reference in workbook.xml
type XMLSheetRef struct {
	XMLName struct{} `xml:"sheet"`
	Name    string   `xml:"name,attr"`
	SheetID int      `xml:"sheetId,attr"`
	RID     string   `xml:"r:id,attr"`
}

// XMLWorksheet represents the xl/worksheets/sheetN.xml structure
type XMLWorksheet struct {
	XMLName       struct{}          `xml:"worksheet"`
	Xmlns         string            `xml:"xmlns,attr"`
	SheetFormatPr *XMLSheetFormatPr `xml:"sheetFormatPr,omitempty"`
	SheetViews    *XMLSheetViews    `xml:"sheetViews,omitempty"`
	Cols          *XMLCols          `xml:"cols,omitempty"`
	SheetData     XMLSheetData      `xml:"sheetData"`
	MergeCells    *XMLMergeCells    `xml:"mergeCells,omitempty"`
}

// XMLSheetFormatPr holds default row/column settings for the sheet
type XMLSheetFormatPr struct {
	XMLName          struct{} `xml:"sheetFormatPr"`
	DefaultRowHeight float64  `xml:"defaultRowHeight,attr"`
	DefaultColWidth  float64  `xml:"defaultColWidth,attr,omitempty"`
	BaseColWidth     int      `xml:"baseColWidth,attr,omitempty"`
}

// XMLSheetViews contains sheet view settings
type XMLSheetViews struct {
	XMLName   struct{}       `xml:"sheetViews"`
	SheetView []XMLSheetView `xml:"sheetView"`
}

// XMLSheetView represents a sheet view
type XMLSheetView struct {
	XMLName           struct{} `xml:"sheetView"`
	WorkbookViewID    int      `xml:"workbookViewId,attr"`
	ShowGridLines     *bool    `xml:"showGridLines,attr,omitempty"`
	ShowRowColHeaders *bool    `xml:"showRowColHeaders,attr,omitempty"`
	Pane              *XMLPane `xml:"pane,omitempty"`
}

// XMLPane represents freeze panes
type XMLPane struct {
	XMLName     struct{} `xml:"pane"`
	XSplit      int      `xml:"xSplit,attr,omitempty"`
	YSplit      int      `xml:"ySplit,attr,omitempty"`
	TopLeftCell string   `xml:"topLeftCell,attr,omitempty"`
	ActivePane  string   `xml:"activePane,attr,omitempty"`
	State       string   `xml:"state,attr,omitempty"`
}

// XMLCols contains column definitions
type XMLCols struct {
	XMLName struct{} `xml:"cols"`
	Col     []XMLCol `xml:"col"`
}

// XMLCol represents a column definition
type XMLCol struct {
	XMLName     struct{} `xml:"col"`
	Min         int      `xml:"min,attr"`
	Max         int      `xml:"max,attr"`
	Width       float64  `xml:"width,attr"`
	CustomWidth bool     `xml:"customWidth,attr,omitempty"`
}

// XMLSheetData contains rows
type XMLSheetData struct {
	XMLName struct{} `xml:"sheetData"`
	Rows    []XMLRow `xml:"row"`
}

// XMLRow represents a row in the worksheet
type XMLRow struct {
	XMLName      struct{}  `xml:"row"`
	R            int       `xml:"r,attr"`
	Height       float64   `xml:"ht,attr,omitempty"`
	CustomHeight bool      `xml:"customHeight,attr,omitempty"`
	Cells        []XMLCell `xml:"c"`
}

// XMLCell represents a cell in a row
type XMLCell struct {
	XMLName struct{}    `xml:"c"`
	R       string      `xml:"r,attr"`
	S       *int        `xml:"s,attr,omitempty"`
	T       string      `xml:"t,attr,omitempty"`
	IS      *XMLIS      `xml:"is,omitempty"`
	V       *string     `xml:"v,omitempty"`
	F       *XMLFormula `xml:"f,omitempty"`
}

// XMLFormula represents a cell formula
type XMLFormula struct {
	XMLName struct{} `xml:"f"`
	Text    string   `xml:",chardata"`
}

// XMLIS represents inline string
type XMLIS struct {
	XMLName struct{} `xml:"is"`
	T       string   `xml:"t"`
}

// XMLMergeCells represents merged cells
type XMLMergeCells struct {
	XMLName struct{}       `xml:"mergeCells"`
	Count   int            `xml:"count,attr"`
	Merge   []XMLMergeCell `xml:"mergeCell"`
}

// XMLMergeCell represents a single merged cell range
type XMLMergeCell struct {
	XMLName struct{} `xml:"mergeCell"`
	Ref     string   `xml:"ref,attr"`
}

// XMLSharedStrings represents xl/sharedStrings.xml
type XMLSharedStrings struct {
	XMLName     struct{} `xml:"sst"`
	Xmlns       string   `xml:"xmlns,attr"`
	Count       int      `xml:"count,attr"`
	UniqueCount int      `xml:"uniqueCount,attr"`
}

// XMLStyleSheet represents xl/styles.xml
type XMLStyleSheet struct {
	XMLName      struct{}        `xml:"styleSheet"`
	Xmlns        string          `xml:"xmlns,attr"`
	Fonts        XMLFonts        `xml:"fonts"`
	Fills        XMLFills        `xml:"fills"`
	Borders      XMLBorders      `xml:"borders"`
	CellStyleXfs XMLCellStyleXfs `xml:"cellStyleXfs"`
	CellXfs      XMLCellXfs      `xml:"cellXfs"`
	CellStyles   XMLCellStyles   `xml:"cellStyles"`
}

// XMLFonts contains font definitions
type XMLFonts struct {
	XMLName struct{}  `xml:"fonts"`
	Count   int       `xml:"count,attr"`
	Font    []XMLFont `xml:"font"`
}

// XMLFont represents a font
type XMLFont struct {
	XMLName struct{}        `xml:"font"`
	Sz      XMLFontSize     `xml:"sz"`
	Name    XMLFontName     `xml:"name"`
	Family  *XMLFontFamily  `xml:"family,omitempty"`
	Charset *XMLFontCharset `xml:"charset,omitempty"`
	Color   *XMLFontColor   `xml:"color,omitempty"`
	B       *XMLBold        `xml:"b,omitempty"`
	I       *XMLItalic      `xml:"i,omitempty"`
	U       *XMLUnderline   `xml:"u,omitempty"`
}

// XMLFontColor represents font color
type XMLFontColor struct {
	XMLName struct{} `xml:"color"`
	RGB     string   `xml:"rgb,attr,omitempty"`
	Theme   int      `xml:"theme,attr,omitempty"`
}

// XMLFontFamily represents font family classification (0 unknown, 1 Roman, 2 Swiss, 3 Modern, 4 Script, 5 Decorative)
type XMLFontFamily struct {
	XMLName struct{} `xml:"family"`
	Val     int      `xml:"val,attr"`
}

// XMLFontCharset represents font charset (0 default)
type XMLFontCharset struct {
	XMLName struct{} `xml:"charset"`
	Val     int      `xml:"val,attr"`
}

// XMLBold represents bold font
type XMLBold struct {
	XMLName struct{} `xml:"b"`
}

// XMLItalic represents italic font
type XMLItalic struct {
	XMLName struct{} `xml:"i"`
}

// XMLUnderline represents underline font
type XMLUnderline struct {
	XMLName struct{} `xml:"u"`
}

// XMLFontSize represents font size
type XMLFontSize struct {
	XMLName struct{} `xml:"sz"`
	Val     string   `xml:"val,attr"`
}

// XMLFontName represents font name
type XMLFontName struct {
	XMLName struct{} `xml:"name"`
	Val     string   `xml:"val,attr"`
}

// XMLFills contains fill definitions
type XMLFills struct {
	XMLName struct{}  `xml:"fills"`
	Count   int       `xml:"count,attr"`
	Fill    []XMLFill `xml:"fill"`
}

// XMLFill represents a fill
type XMLFill struct {
	XMLName     struct{}       `xml:"fill"`
	PatternFill XMLPatternFill `xml:"patternFill"`
}

// XMLPatternFill represents pattern fill
type XMLPatternFill struct {
	XMLName     struct{}      `xml:"patternFill"`
	PatternType string        `xml:"patternType,attr"`
	FgColor     *XMLFillColor `xml:"fgColor,omitempty"`
	BgColor     *XMLBgColor   `xml:"bgColor,omitempty"`
}

// XMLFillColor represents foreground fill color
type XMLFillColor struct {
	RGB     string `xml:"rgb,attr,omitempty"`
	Indexed int    `xml:"indexed,attr,omitempty"`
}

// XMLBgColor represents background fill color
type XMLBgColor struct {
	RGB     string `xml:"rgb,attr,omitempty"`
	Indexed int    `xml:"indexed,attr,omitempty"`
}

// XMLBorders contains border definitions
type XMLBorders struct {
	XMLName struct{}    `xml:"borders"`
	Count   int         `xml:"count,attr"`
	Border  []XMLBorder `xml:"border"`
}

// XMLBorder represents a border
type XMLBorder struct {
	XMLName struct{}      `xml:"border"`
	Left    XMLBorderSide `xml:"left"`
	Right   XMLBorderSide `xml:"right"`
	Top     XMLBorderSide `xml:"top"`
	Bottom  XMLBorderSide `xml:"bottom"`
}

// XMLBorderSide represents a border side
type XMLBorderSide struct {
	Style string          `xml:"style,attr,omitempty"`
	Color *XMLBorderColor `xml:"color,omitempty"`
}

// XMLBorderColor represents border color
type XMLBorderColor struct {
	XMLName struct{} `xml:"color"`
	RGB     string   `xml:"rgb,attr,omitempty"`
	Indexed int      `xml:"indexed,attr,omitempty"`
}

// XMLCellXfs contains cell format definitions
type XMLCellXfs struct {
	XMLName struct{} `xml:"cellXfs"`
	Count   int      `xml:"count,attr"`
	Xf      []XMLXf  `xml:"xf"`
}

// XMLXf represents a cell format
type XMLXf struct {
	XMLName     struct{} `xml:"xf"`
	NumFmtID    int      `xml:"numFmtId,attr"`
	FontID      int      `xml:"fontId,attr"`
	FillID      int      `xml:"fillId,attr"`
	BorderID    int      `xml:"borderId,attr"`
	ApplyFont   bool     `xml:"applyFont,attr,omitempty"`
	ApplyFill   bool     `xml:"applyFill,attr,omitempty"`
	ApplyBorder bool     `xml:"applyBorder,attr,omitempty"`
}

// XMLCellStyleXfs represents base (named) styles
type XMLCellStyleXfs struct {
	XMLName struct{} `xml:"cellStyleXfs"`
	Count   int      `xml:"count,attr"`
	Xf      []XMLXf  `xml:"xf"`
}

// XMLCellStyles represents the list of style names
type XMLCellStyles struct {
	XMLName struct{}       `xml:"cellStyles"`
	Count   int            `xml:"count,attr"`
	Cell    []XMLCellStyle `xml:"cellStyle"`
}

// XMLCellStyle represents a named cell style
type XMLCellStyle struct {
	XMLName   struct{} `xml:"cellStyle"`
	Name      string   `xml:"name,attr"`
	XfID      int      `xml:"xfId,attr"`
	BuiltinID int      `xml:"builtinId,attr,omitempty"`
}
