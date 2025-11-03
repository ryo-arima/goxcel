# XLSX Models

## Core Structures

### Book
Workbook container.
```go
type Book struct {
    Sheets []*Sheet
}
```

### Sheet
Worksheet with cells and components.
```go
type Sheet struct {
    Name   string
    Cells  []*Cell
    Merges []Merge
    Images []Image
    Shapes []Shape
    Charts []Chart
}
```

### Cell
Individual cell with value and metadata.
```go
type Cell struct {
    Ref   string      // "A1", "B2", etc.
    Value string      // Cell content
    Type  CellType    // Number, String, Boolean, etc.
    Style *CellStyle  // Formatting (bold, italic)
}
```

### CellType
Enumeration of cell types.
```go
type CellType int

const (
    CellTypeAuto     // Infer from value
    CellTypeString   // Text
    CellTypeNumber   // Numeric
    CellTypeBoolean  // True/False
    CellTypeDate     // ISO date
    CellTypeFormula  // Excel formula
)
```

### CellStyle
Visual formatting.
```go
type CellStyle struct {
    Bold      bool
    Italic    bool
    Underline bool
}
```

## Type Inference

Automatic type detection:
- Starts with `=` → Formula
- `true`/`false` → Boolean
- Matches number pattern → Number
- Matches date pattern → Date
- Default → String
