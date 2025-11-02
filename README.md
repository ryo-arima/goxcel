# goxcel

An Excel (.xlsx) generation library for Go. It parses grid-oriented .gxl templates and produces .xlsx files without external dependencies.

## Features
- ✅ Generate Excel files using template-based approach
- ✅ Grid-oriented syntax with pipe-delimited tables
- ✅ Support for cell values, formulas, and merges
- ✅ Template variables with `{{ expr }}` syntax
- ✅ Control structures (`<For>`, `<If>`)
- ✅ Components: Images, Shapes, Charts, and Pivot Tables
- ✅ Pure Go implementation - no external XLSX libraries required
- ✅ Structured logging with message codes

## Prerequisites
- **Go 1.21 or higher**
- No external dependencies (uses only Go standard library)

## Installation

```bash
go get github.com/ryo-arima/goxcel
```

## Quick Start

### Build the CLI

```bash
go build -o goxcel ./cmd
```

### Generate Excel from Template

```bash
./goxcel generate --template .examples/sample.gxl --data .examples/sample-data.json --output output.xlsx
```

### Command Options

- `--template, -t`: Path to .gxl template file (required)
- `--data, -d`: Path to JSON data file (optional)
- `--output, -o`: Output .xlsx file path (optional)
- `--dry-run`: Print summary without writing file

## Directory Structure

```
goxcel/
├── cmd/                # CLI entry point
├── pkg/
│   ├── config/         # Configuration with logger support
│   ├── controller/     # CLI command handlers
│   ├── model/          # Data models (GXL AST, XLSX structures)
│   ├── repository/     # File I/O (GXL parser, XLSX writer)
│   ├── usecase/        # Business logic (template rendering)
│   └── util/           # Utilities (structured logger)
├── .examples/          # Template examples and sample data
└── test/               # Tests
```

## .gxl Template Format

### Basic Syntax

.gxl files use XML-like tags with grid-oriented content:

- **Sheet definition**: `<Sheet name="SheetName">...</Sheet>`
- **Grid tables**: Pipe-delimited rows `| cell1 | cell2 | cell3 |`
- **Variables**: `{{ variable.path }}` for data injection
- **Control structures**:
  - `<For each="item in items">...</For>` - Loop over array data
  - `<If cond="condition">...<Else>...</Else></If>` - Conditional rendering
- **Components**:
  - `<Merge range="A1:C1" />` - Merge cell ranges
  - `<Image ref="B3" src="path/to/image.png" width="120" height="60" />` - Insert images
  - `<Shape ref="D3" kind="rectangle" text="Hello" width="100" height="40" />` - Add shapes
  - `<Chart ref="F3" type="column" dataRange="A9:C20" title="Sales" width="420" height="240" />` - Create charts
  - `<Pivot ref="F15" sourceRange="A9:C200" rows="Name" values="SUM:Price" />` - Pivot tables

### Example Template

**.examples/sample.gxl**:

```xml
<Book>
<Sheet name="Components demo">
<Grid>
| Components demo |
</Grid>
<Merge range="A1:C1" />

<Grid>
| Name | Quantity | Price |
</Grid>

<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.qty }} | {{ item.price }} |
</Grid>
</For>

<Grid>
|  |  Total | =SUM(C4:C6) |
</Grid>

<Image ref="E2" src="assets/logo.png" width="120" height="60" />
<Shape ref="E5" kind="rectangle" text="Sample Shape" width="150" height="50" />
<Chart ref="E8" type="column" dataRange="A3:C6" title="Sales Chart" width="420" height="240" />
<Pivot ref="E16" sourceRange="A3:C6" rows="Name" values="SUM:Price" />
</Sheet>
</Book>
```

### Sample Data File

**.examples/sample-data.json**:

```json
{
  "items": [
    {"name": "Apple", "qty": 10, "price": 100},
    {"name": "Banana", "qty": 20, "price": 200},
    {"name": "Cherry", "qty": 30, "price": 300}
  ]
}
```

### Output

The above template and data will generate an Excel file with:
- 13 cells with data
- 1 merged cell range (A1:C1)
- 1 image placeholder
- 1 shape placeholder
- 1 chart placeholder
- 1 pivot table placeholder

## Implementation Details

### Architecture

**Clean Architecture** with clear layer separation:

- **Config Layer**: Configuration management with dependency injection
- **Controller Layer**: CLI command handlers with logging
- **UseCase Layer**: Business logic for template rendering
- **Repository Layer**: File I/O for GXL parsing and XLSX writing
- **Model Layer**: Data structures for GXL AST and XLSX XML

### XLSX Generation

The library generates XLSX files by:

1. **Creating ZIP structure**: XLSX is a ZIP archive containing XML files
2. **Generating XML files**:
   - `_rels/.rels` - Package relationships
   - `[Content_Types].xml` - Content type definitions
   - `xl/workbook.xml` - Workbook structure
   - `xl/worksheets/sheet*.xml` - Worksheet data
   - `xl/styles.xml` - Style definitions
   - `xl/sharedStrings.xml` - Shared string table
3. **Using encoding/xml**: Type-safe XML marshaling with struct tags
4. **Supporting Excel features**: Cells, formulas, merges, inline strings

### Structured Logging

Message code-based logging for traceability:

- `SY-*`: System layer (application lifecycle)
- `FS-*`: File system operations (read, write, directory)
- `R-*`: Repository layer (parsing, writing)
- `C-*`: Controller layer (command execution)
- `U-*`: UseCase layer (rendering, business logic)
- `M-*`: Model layer (validation, conversion)
- `GXL-*`: GXL processing (parsing, rendering)
- `XLSX-*`: XLSX processing (generation, reading)
- `XML-*`: XML processing (marshaling, unmarshaling)

Example log output:
```
[2024-01-01T12:00:00Z] INFO [C-I1] goxcel/cli: Starting generate command {"template":".examples/sample.gxl","output":"output.xlsx"}
[2024-01-01T12:00:00Z] DEBUG [GXL-P1] goxcel/cli: GXL template parsed successfully {"sheets":1}
[2024-01-01T12:00:00Z] INFO [U-R1] goxcel/cli: Rendering template
[2024-01-01T12:00:00Z] INFO [R-W1] goxcel/cli: Writing XLSX file {"output":"output.xlsx"}
[2024-01-01T12:00:00Z] INFO [C-C1] goxcel/cli: Successfully generated XLSX file {"output":"output.xlsx"}
```

## GXL Template Model

### AST Structure

```go
type GXL struct {
    Sheets []SheetTag
}

type SheetTag struct {
    Name  string
    Nodes []any  // Union of all node types
}
```

### Node Types

- `GridTag` - Contains rows of cells (pipe-delimited tables)
- `GridRowTag` - Single row with cell values
- `MergeTag` - Merge cell range definition
- `ForTag` - Loop construct with iteration variable
- `IfTag` - Conditional rendering with Then/Else branches
- `ImageTag` - Image placement with dimensions
- `ShapeTag` - Shape drawing with style
- `ChartTag` - Chart definition with data range
- `PivotTag` - Pivot table configuration

## XLSX Model

### Book Structure

```go
type Book struct {
    Sheets []*Sheet
}

type Sheet struct {
    Name    string
    Cells   []*Cell
    Merges  []*Merge
    Images  []*Image
    Shapes  []*Shape
    Charts  []*Chart
    Pivots  []*PivotTable
}

type Cell struct {
    Ref   string    // A1 reference
    Value string    // Cell value
    Type  CellType  // String, number, formula, etc.
    Style string    // Style name/ID
}
```

## Current Status

### ✅ Implemented
- GXL template parsing (XML-based with encoding/xml)
- Grid table parsing (pipe-delimited format)
- Template variable expansion (`{{ expr }}`)
- Control structures (`<For>`, `<If>`)
- Component placeholders (Image, Shape, Chart, Pivot)
- XLSX generation without external libraries
- ZIP-based XLSX file structure
- XML marshaling for all XLSX components
- Cell references and formulas
- Merge cell ranges
- Structured logging system
- CLI with generate command

### 🚧 Planned
- Style definitions (fonts, colors, borders, fills)
- Column width and row height
- Cell formatting (number formats, alignment)
- Actual image embedding (currently placeholders)
- Chart data series configuration
- Pivot table field configuration
- Named ranges
- Data validation
- Conditional formatting
- Page setup (headers, footers, print settings)
- More comprehensive test coverage

## Development

### Run Tests

```bash
go test ./...
```

### Build and Test

```bash
# Build CLI
go build -o goxcel ./cmd

# Test with sample template
./goxcel generate --template .examples/sample.gxl --data .examples/sample-data.json --output output.xlsx

# Dry run (no file output)
./goxcel generate --template .examples/sample.gxl --data .examples/sample-data.json --dry-run
```

### Enable Debug Logging

Modify `pkg/config/base.go` to set log level to "DEBUG":

```go
logger := util.NewLogger(util.LoggerConfig{
    Component:    "goxcel",
    Service:      "cli",
    Level:        "DEBUG",  // Change from INFO to DEBUG
    Structured:   false,
    EnableCaller: false,
    Output:       "stdout",
})
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2024 Ryo Arima

## References

- [Office Open XML Standard](http://officeopenxml.com/)
- [ECMA-376 Standard](https://www.ecma-international.org/publications-and-standards/standards/ecma-376/)
- [SpreadsheetML Reference](https://docs.microsoft.com/en-us/openspecs/office_standards/ms-xlsx/)
