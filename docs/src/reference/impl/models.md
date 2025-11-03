# Data Models

## GXL Models (`pkg/model/gxl.go`)

AST nodes representing parsed template structure.

**Key types**:
- `GXL`: Root template structure
- `SheetTag`: Sheet definition
- `GridTag`: Table/grid with rows
- `ForTag`: Loop iteration
- `AnchorTag`: Position anchor
- `MergeTag`: Cell merge range

## XLSX Models (`pkg/model/xlsx.go`)

Output structures for Excel generation.

**Key types**:
- `Book`: Workbook container
- `Sheet`: Worksheet with cells
- `Cell`: Individual cell (ref, value, type, style)
- `CellType`: String, Number, Boolean, Date, Formula, Auto
- `CellStyle`: Bold, Italic, Underline flags

## Model Flow

```
GXL (Template) → Usecase → XLSX (Output)
```

See detailed model definitions:
- [GXL Models](./models/gxl.md)
- [XLSX Models](./models/xlsx.md)
