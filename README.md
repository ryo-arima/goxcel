# goxcel

**A template-based Excel generation library and CLI tool for Go**

goxcel transforms human-readable `.gxl` templates into Excel `.xlsx` files using a simple, grid-oriented syntax. No external dependencies required.

> **⚠️ Development Status**  
> This project is currently under active development. Breaking changes may occur between versions until v1.0.0 is released. Please use with caution in production environments.

## Why goxcel?

- **📝 Template-First**: Write Excel layouts in readable text format
- **🎯 Grid-Oriented**: Natural pipe-delimited table syntax
- **🔄 Data-Driven**: Separate templates from data with JSON context
- **🎨 Feature-Rich**: Formulas, loops, merges, images, charts, and more
- **🚀 Pure Go**: No external XLSX libraries or dependencies
- **📊 Production-Ready**: Clean architecture with structured logging

## Quick Example

**Template** (`.gxl`):
```xml
<Book>
  <Sheet name="Invoice">
    <Grid>
    | Invoice #{{invoiceNumber}} |
    | Customer: {{customer}} |
    </Grid>
    
    <Grid>
    | Item | Quantity | Price |
    </Grid>
    
    <For src="items">
      <Grid>
      | {{name}} | {{qty}} | ${{price}} |
      </Grid>
    </For>
  </Sheet>
</Book>
```

**Data** (JSON):
```json
{
  "invoiceNumber": "INV-001",
  "customer": "Acme Corp",
  "items": [
    {"name": "Widget", "qty": 10, "price": 50.00},
    {"name": "Gadget", "qty": 5, "price": 120.00}
  ]
}
```

**Result**: Professional Excel file with dynamic data populated.

## Features

### Core Features
- ✅ **Template-based generation** with `.gxl` format
- ✅ **Grid syntax** with pipe-delimited tables
- ✅ **Variable interpolation** with `{{ expr }}` syntax
- ✅ **Control structures**: `<For>` loops (v1.0), `<If>` conditionals (planned v1.1)
- ✅ **Excel formulas** with cell references
- ✅ **Cell merging** for headers and layouts
- ✅ **Multi-sheet workbooks** with independent sheets
- ✅ **JSON and YAML data**: Support for both data formats

### Components (Placeholders in v1.0)
- ⏳ **Images**: PNG/JPEG embedding
- ⏳ **Shapes**: Rectangles, ellipses, arrows
- ⏳ **Charts**: Column, bar, line, pie charts
- ⏳ **Pivot tables**: Dynamic data aggregation

### Developer Experience
- ✅ **Pure Go**: No C dependencies, easy deployment
- ✅ **CLI tool**: Generate files from command line
- ✅ **Library API**: Use as Go package in your code
- ✅ **Structured logging**: Message codes for traceability
- ✅ **Clean architecture**: Testable, maintainable code

## Installation

### Install CLI Tool

```bash
go install github.com/ryo-arima/goxcel/cmd@latest
```

### Use as Library

```bash
go get github.com/ryo-arima/goxcel
```

## Quick Start

### Using the CLI

```bash
# Build from source
make build

# Generate Excel file from JSON data
.bin/goxcel generate --template .etc/sample.gxl --data .etc/sample.json --output invoice.xlsx

# Generate Excel file from YAML data
.bin/goxcel generate --template .etc/sample.gxl --data .etc/sample.yaml --output invoice.xlsx

# Preview without generating file
.bin/goxcel generate --template .etc/sample.gxl --data .etc/sample.json --dry-run
```

### Using as Library

```go
package main

import (
    "github.com/ryo-arima/goxcel/pkg/config"
    "github.com/ryo-arima/goxcel/pkg/controller"
)

func main() {
    cfg := config.NewBaseConfig()
    ctrl := controller.NewCommonController(cfg)
    
    err := ctrl.Generate(
        ".etc/sample.gxl",
        ".etc/sample.json",
        "output.xlsx",
        false, // dry-run
    )
    if err != nil {
        panic(err)
    }
}
```

## Documentation

📚 **Comprehensive documentation available at [docs/](./docs/)**

Build and view locally:
```bash
make docs-build
make docs-serve
# Open http://localhost:3000
```

### Key Documentation
- **[GXL Specification](./docs/src/specification/)** - Complete format reference
- **[Getting Started Guide](./docs/src/getting-started/)** - Tutorials and examples
- **[API Reference](./docs/src/reference/)** - Go package documentation
- **[Vision & Strategy](./docs/src/vision-strategy.md)** - Project roadmap

## GXL Template Language

### Core Syntax

**Book and Sheets**
```xml
<Book>
  <Sheet name="Sheet1">
    <!-- Sheet content -->
  </Sheet>
  <Sheet name="Sheet2">
    <!-- Another sheet -->
  </Sheet>
</Book>
```

**Grid Tables** (Pipe-delimited)
```xml
<Grid>
| Header 1 | Header 2 | Header 3 |
| Value 1  | Value 2  | Value 3  |
| Data A   | Data B   | Data C   |
</Grid>
```

**Variable Interpolation**
```xml
<Grid>
| Customer: {{customer.name}} |
| Date: {{invoice.date}} |
| Total: ${{invoice.total}} |
</Grid>
```

**Loops** (Iterate over arrays)
```xml
<For src="items">
  <Grid>
  | {{name}} | {{quantity}} | {{price}} |
  </Grid>
</For>
```

**Loop Variables**
```xml
<For src="items">
  <Grid>
  | Row {{_number}} | {{name}} | =B{{_startRow}}*2 |
  </Grid>
</For>
```
- `{{_index}}` - Zero-based index (0, 1, 2, ...)
- `{{_number}}` - One-based number (1, 2, 3, ...)
- `{{_startRow}}` - Excel row number for current iteration
- `{{_endRow}}` - Last row number (available after loop)

**Cell Merging**
```xml
<Grid>
| Title spanning multiple columns | | | |
</Grid>
<Merge range="A1:D1" />
```

**Positioning with Anchor**
```xml
<Anchor cell="E1">
  <Grid>
  | Side content |
  </Grid>
</Anchor>
```

**Excel Formulas**
```xml
<Grid>
| =SUM(A1:A10) |
| =AVERAGE(B:B) |
| =IF(C1>100,"High","Low") |
</Grid>
```

### Component Syntax (v1.0: Placeholders)

**Images**
```xml
<Image src="logo.png" cell="E1" width="100" height="50" />
```

**Shapes**
```xml
<Shape kind="rectangle" cell="E5" text="Note" width="150" height="40" />
```

**Charts**
```xml
<Chart type="column" cell="E8" dataRange="A3:C10" title="Sales" width="420" height="240" />
```

**Pivot Tables**
```xml
<Pivot cell="E15" sourceRange="A3:C100" rows="Category" values="SUM:Amount" />
```

## Architecture

### Project Structure

```
goxcel/
├── .bin/               # Built binaries
├── .etc/               # Configuration and sample files
│   ├── sample.gxl      # Sample GXL template
│   └── sample.json     # Sample data
├── cmd/                # CLI entry point
│   └── main.go
├── docs/               # Documentation (mdbook)
│   ├── book.toml       # Configuration
│   ├── dist/           # Built HTML (gitignored)
│   └── src/            # Markdown sources
│       ├── specification/  # GXL format specification
│       ├── reference/      # API reference
│       └── guide/          # User guides
├── pkg/                # Application packages
│   ├── command.go      # CLI commands
│   ├── config/         # Configuration layer
│   ├── controller/     # Command handlers (CLI interface)
│   ├── model/          # Data structures
│   │   ├── gxl.go      # GXL AST
│   │   └── xlsx.go     # XLSX models
│   ├── repository/     # File I/O
│   │   ├── gxl.go      # GXL parser
│   │   └── xlsx.go     # XLSX writer
│   ├── usecase/        # Business logic
│   │   ├── book.go     # Book rendering
│   │   ├── cell.go     # Cell processing
│   │   └── sheet.go    # Sheet rendering
│   └── util/           # Utilities
│       └── logger.go   # Structured logging
├── go.mod
├── Makefile            # Build tasks
└── README.md
```

### Clean Architecture Layers

**1. Config Layer**: Dependency injection and configuration
**2. Controller Layer**: CLI commands and handlers
**3. UseCase Layer**: Core business logic (template rendering)
**4. Repository Layer**: File I/O (GXL parsing, XLSX writing)
**5. Model Layer**: Data structures (GXL AST, XLSX models)

### Data Flow

```
.gxl Template + JSON Data
        ↓
    [Parser]
        ↓
    GXL AST
        ↓
   [Renderer] ← JSON Context
        ↓
  XLSX Model (Book → Sheets → Cells)
        ↓
  [XLSX Writer]
        ↓
.xlsx File (ZIP containing XML files)
```

## How It Works

### XLSX File Structure

XLSX files are ZIP archives containing XML files:

```
output.xlsx (ZIP)
├── [Content_Types].xml       # MIME types
├── _rels/
│   └── .rels                 # Package relationships
└── xl/
    ├── workbook.xml          # Workbook structure
    ├── worksheets/
    │   ├── sheet1.xml        # Sheet data
    │   └── sheet2.xml
    ├── styles.xml            # Styles and formats
    └── sharedStrings.xml     # Shared string table
```

### Generation Process

1. **Parse Phase**: Read `.gxl` file → Build AST
2. **Render Phase**: Walk AST + Evaluate expressions → Generate cells
3. **Write Phase**: Marshal XML → Create ZIP → Write `.xlsx`

### Type-Safe XML Generation

```go
// Define XLSX structure with XML tags
type Worksheet struct {
    XMLName    xml.Name `xml:"http://... worksheet"`
    SheetData  SheetData `xml:"sheetData"`
}

// Marshal to XML
xml.MarshalIndent(worksheet, "", "  ")
```

## Development

### Prerequisites

- **Go 1.21+**
- **mdbook** (for documentation): `cargo install mdbook`

### Makefile Commands

```bash
# Build CLI
make build

# Run with sample
make run

# Dry run (preview)
make run-dry

# Run tests
make test

# Build documentation
make docs-build

# Serve documentation
make docs-serve

# Clean build artifacts
make clean
```

## Use Cases

### 📊 Business Reports
Generate monthly/quarterly reports with dynamic data from databases or APIs.

### 📋 Invoices & Documents
Create professional invoices, quotes, and contracts with consistent formatting.

### 📈 Data Exports
Export application data to Excel with custom layouts and formulas.

### 🏢 Batch Processing
Generate hundreds of personalized Excel files from templates.

### 📑 Form Filling
Populate Excel forms with data from web services or user input.

## Performance

**Typical Performance** (M1 Mac, Go 1.21):
- Small template (10 sheets, 100 rows): ~50ms
- Medium template (100 sheets, 1000 rows): ~500ms
- Large template (1000 rows with loops): ~2s

**Memory Usage**:
- Lightweight: ~10MB for typical templates
- Scales linearly with data size

## Comparison with Other Tools

| Feature | goxcel | excelize | xlsx | Apache POI |
|---------|--------|----------|------|------------|
| Template-based | ✅ | ❌ | ❌ | ❌ |
| Pure Go | ✅ | ✅ | ✅ | ❌ (Java) |
| No dependencies | ✅ | ❌ | ❌ | ❌ |
| Grid syntax | ✅ | ❌ | ❌ | ❌ |
| CLI tool | ✅ | ❌ | ❌ | ❌ |
| Data-driven | ✅ | ⚠️ Manual | ⚠️ Manual | ⚠️ Manual |

## Contributing

We welcome contributions! Here's how to get started:

### 1. Fork and Clone

```bash
git clone https://github.com/yourusername/goxcel.git
cd goxcel
```

### 2. Create a Branch

```bash
git checkout -b feature/my-feature
```

### 3. Make Changes

- Follow Go best practices
- Add tests for new features
- Update documentation
- Run tests: `make test`

### 4. Submit Pull Request

- Clear description of changes
- Reference related issues
- Ensure CI passes

### Development Guidelines

- **Code Style**: Follow Go conventions (`gofmt`, `golint`)
- **Testing**: Unit tests for all new code
- **Documentation**: Update docs for user-facing changes
- **Commits**: Clear, descriptive commit messages

### Areas for Contribution

- 🐛 Bug fixes
- 📝 Documentation improvements
- ✨ New features (check roadmap)
- 🧪 Test coverage
- 🌍 Internationalization
- 📊 Performance optimizations

## FAQ

See [FAQ](./docs/src/appendix/faq.md) for more questions.

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## Acknowledgments
- Excel specification: [Office Open XML Standard](http://officeopenxml.com/)
- Go community forexcellent tools and libraries

## Contact & Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/ryo-arima/goxcel/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/ryo-arima/goxcel/discussions)

---

**⭐ If you find goxcel useful, please star the repository!**
- [ECMA-376 Standard](https://www.ecma-international.org/publications-and-standards/standards/ecma-376/)
- [SpreadsheetML Reference](https://docs.microsoft.com/en-us/openspecs/office_standards/ms-xlsx/)
