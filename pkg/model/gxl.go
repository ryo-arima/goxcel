package model

// GXL represents the root structure of a .gxl template file.
type GXL struct {
	HeaderTag HeaderTag
	BookTag   BookTag
	Sheets    []SheetTag
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

// StyleTag represents <Style> for applying formatting.
type StyleTag struct {
	Selector   string
	Name       string
	ID         string
	Class      string
	Properties map[string]string
}

// IfTag represents <If cond="..."> for conditional rendering.
type IfTag struct {
	Cond string
	Then []any
	Else []any
}
