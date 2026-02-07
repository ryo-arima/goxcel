package model

// BookNodeType represents the type of node at book level
type BookNodeType int

const (
	BookNodeTypeImport BookNodeType = iota
	BookNodeTypeSheet
)

// BookNode represents a node at book level (Import or Sheet) with order preserved
type BookNode struct {
	Type   BookNodeType
	Import *ImportTag
	Sheet  *SheetTag
}

// GXL represents the root structure of a .gxl template file.
type GXL struct {
	HeaderTag HeaderTag
	BookTag   BookTag
	Imports   []ImportTag // Import tags at book level (deprecated - use BookNodes)
	Sheets    []SheetTag  // (deprecated - use BookNodes)
	BookNodes []BookNode  // Ordered book-level nodes (Import and Sheet in definition order)
}

// HeaderTag holds global metadata for the GXL template.
type HeaderTag struct {
	Title      string
	Version    string
	Encoding   string
	Properties map[string]string
}

// BookTag represents the <Book> element with its attributes.
type BookTag struct {
	Name       string
	Properties map[string]string
}

// SheetTag represents a <Sheet> element within a workbook.
type SheetTag struct {
	Name  string
	Nodes []any
	// Optional sheet-level configuration parsed from attributes or child tags
	Config *SheetConfigTag
}

// SheetConfigTag represents <SheetConfig> for sheet-level settings
type SheetConfigTag struct {
	DefaultRowHeight   float64
	DefaultColumnWidth float64
	FreezePane         string
	ShowGridLines      *bool
	ShowRowColHeaders  *bool
}

// ColumnTag represents <Column> for column width settings
type ColumnTag struct {
	Index int     // Column number (1-based)
	Width float64 // Width in Excel units
}

// RowTag represents <Row> for row height settings
type RowHeightTag struct {
	Index  int     // Row number (1-based)
	Height float64 // Height in points
}

// CellTag represents a single cell value within a grid row.
type CellTag struct {
	Value string
}

// RowTag represents a table row.
type RowTag struct {
	Cells []CellTag
}

// ColTag represents a column definition.
type ColTag struct {
	Index int
	Width int
}

// GridTag represents <Grid> containing pipe-delimited table rows.
type GridTag struct {
	Content string
	Rows    []GridRowTag
	Ref     string // Optional: Starting cell reference (e.g., "A1", "B5")
	// Optional style defaults applied to all cells in this Grid
	FontName    string
	FontSize    int
	FontColor   string // RGB hex without # (e.g., "FF0000")
	FillColor   string // RGB hex without # (e.g., "FFFF00")
	BorderStyle string // Border line style (thin, medium, thick, dashed, dotted, double)
	BorderColor string // RGB hex without #
	BorderSides string // comma-separated: all, top, right, bottom, left
}

// GridRowTag represents a single row parsed from Grid content.
type GridRowTag struct {
	Cells []string
}

// AnchorTag represents <Anchor ref="A1" /> to set the current cell position.
type AnchorTag struct {
	Ref string
}

// MergeTag represents <Merge range="A1:C1" /> to merge cells.
type MergeTag struct {
	Range string
}

// ImageTag represents <Image> for placing an image on the sheet.
type ImageTag struct {
	Ref    string
	Src    string
	Width  int
	Height int
}

// ShapeTag represents <Shape> for drawing shapes with optional text.
type ShapeTag struct {
	Ref    string
	Kind   string
	Text   string
	Width  int
	Height int
	Style  string
}

// ForTag represents <For each="item in items"> for iteration.
type ForTag struct {
	Each string
	Body []any
}

// TableTag represents <Table> containing rows and columns
type TableTag struct {
	Rows []TableRowTag
}

// TableRowTag represents <Row> inside <Table>
type TableRowTag struct {
	Each string        // Optional: "item in items" syntax for looping
	Cols []TableColTag // Child <Col> elements
}

// TableColTag represents <Col> inside <Row>
type TableColTag struct {
	Each    string // Optional: "item in items" syntax for looping  
	Content string // Cell content (no pipe needed)
}

// ChartTag represents <Chart> for embedding charts.
type ChartTag struct {
	Ref       string
	Type      string
	DataRange string
	Title     string
	Width     int
	Height    int
}

// PivotTag represents <Pivot> for pivot table definitions.
type PivotTag struct {
	Ref         string
	SourceRange string
	Rows        string
	Columns     string
	Values      string
	Filters     string
	Options     string
}

// StyleTag represents <Style> for cell styling
type StyleTag struct {
	Ref        string // Cell reference or range (e.g., "A1" or "A1:C3") - replaces Selector
	Selector   string // Deprecated: use Ref
	Name       string
	ID         string
	Class      string
	Properties map[string]string

	// Direct style attributes
	Bold      bool
	Italic    bool
	Underline bool
	FontName  string
	FontSize  int
	FontColor string // RGB hex without # (e.g., "FF0000")
	FillColor string // RGB hex without # (e.g., "FFFF00")
}

// IfTag represents <If cond="..."> for conditional rendering.
type IfTag struct {
	Cond string
	Then []any
	Else []any
}

// ImportTag represents <Import src="..." sheet="..." /> for importing external .gxl files.
type ImportTag struct {
	Src   string // Path to the external .gxl file (relative or absolute)
	Sheet string // Name of the sheet to import (required)
}
