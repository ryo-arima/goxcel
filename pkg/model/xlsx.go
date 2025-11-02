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
	CellTypeString CellType = "string"
)

// Cell represents a single cell value.
type Cell struct {
	Ref   string   // A1 reference
	Value string   // text value (after {{ }} expansion)
	Type  CellType // hint for writer
	Style string   // optional style name/id applied
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
}

// NewSheet creates a sheet with a given name.
func NewSheet(name string) *Sheet { return &Sheet{Name: name} }

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
	XMLName    struct{}       `xml:"worksheet"`
	Xmlns      string         `xml:"xmlns,attr"`
	SheetData  XMLSheetData   `xml:"sheetData"`
	MergeCells *XMLMergeCells `xml:"mergeCells,omitempty"`
}

// XMLSheetData contains rows
type XMLSheetData struct {
	XMLName struct{} `xml:"sheetData"`
	Rows    []XMLRow `xml:"row"`
}

// XMLRow represents a row in the worksheet
type XMLRow struct {
	XMLName struct{}  `xml:"row"`
	R       int       `xml:"r,attr"`
	Cells   []XMLCell `xml:"c"`
}

// XMLCell represents a cell in a row
type XMLCell struct {
	XMLName struct{} `xml:"c"`
	R       string   `xml:"r,attr"`
	T       string   `xml:"t,attr,omitempty"`
	IS      *XMLIS   `xml:"is,omitempty"`
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
	XMLName struct{}   `xml:"styleSheet"`
	Xmlns   string     `xml:"xmlns,attr"`
	Fonts   XMLFonts   `xml:"fonts"`
	Fills   XMLFills   `xml:"fills"`
	Borders XMLBorders `xml:"borders"`
	CellXfs XMLCellXfs `xml:"cellXfs"`
}

// XMLFonts contains font definitions
type XMLFonts struct {
	XMLName struct{}  `xml:"fonts"`
	Count   int       `xml:"count,attr"`
	Font    []XMLFont `xml:"font"`
}

// XMLFont represents a font
type XMLFont struct {
	XMLName struct{}    `xml:"font"`
	Sz      XMLFontSize `xml:"sz"`
	Name    XMLFontName `xml:"name"`
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
	XMLName     struct{} `xml:"patternFill"`
	PatternType string   `xml:"patternType,attr"`
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
}

// XMLCellXfs contains cell format definitions
type XMLCellXfs struct {
	XMLName struct{} `xml:"cellXfs"`
	Count   int      `xml:"count,attr"`
	Xf      []XMLXf  `xml:"xf"`
}

// XMLXf represents a cell format
type XMLXf struct {
	XMLName  struct{} `xml:"xf"`
	NumFmtID int      `xml:"numFmtId,attr"`
	FontID   int      `xml:"fontId,attr"`
	FillID   int      `xml:"fillId,attr"`
	BorderID int      `xml:"borderId,attr"`
}
