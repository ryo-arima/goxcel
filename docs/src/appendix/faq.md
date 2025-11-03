# FAQ

## General

### What is goxcel?

goxcel is a template-based Excel file generator that uses GXL (Grid eXchange Language) format to define Excel workbook structures with data binding capabilities.

### What file formats does goxcel support?

- **Input**: GXL (`.gxl`) templates, JSON/YAML data files
- **Output**: Excel (`.xlsx`) files

### Is goxcel production-ready?

goxcel is currently at version 0.1.1 (Stable) with v1.0 features implemented including grid layouts, loops, expressions, and markdown styling.

## Templates

### How do I create a GXL template?

GXL templates are XML files with a `.gxl` extension. Basic structure:

```xml
<Book name="MyWorkbook">
  <Sheet name="Sheet1">
    <Grid>
      | Header1 | Header2 |
      | {{ .value1 }} | {{ .value2 }} |
    </Grid>
  </Sheet>
</Book>
```

### Can I use loops in templates?

Yes, use the `<For>` tag:

```xml
<For each="item in items">
  <Grid>
    | {{ .item.name }} | {{ .item.value }} |
  </Grid>
</For>
```

### How do I apply cell styling?

Use markdown syntax or type hints:
- `**bold text**` for bold
- `_italic text_` for italic
- `{{ .value:number }}` for type hints

## Data

### What data formats are supported?

JSON and YAML files are supported for data input.

### How do I access nested data?

Use dot notation: `{{ .user.profile.name }}`

### Can I use expressions?

Yes, mustache-style expressions `{{ .path }}` with automatic type inference.

## CLI

### How do I generate an Excel file?

```bash
goxcel generate --template template.gxl --data data.json --output output.xlsx
```

### Can I preview without creating a file?

Yes, use the `--dry-run` flag:

```bash
goxcel generate --template template.gxl --data data.json --dry-run
```

## Troubleshooting

### Template parsing fails

Check:
- XML syntax is valid
- All tags are properly closed
- Grid syntax uses `|` delimiters

### Data not appearing in output

Check:
- JSON/YAML data structure matches template paths
- Expression syntax is correct: `{{ .path }}`
- Data file is valid JSON/YAML

### Excel file won't open

Ensure:
- Output path is writable
- No special characters in filename
- Sufficient disk space

## Performance

### How large can templates be?

goxcel can handle templates with thousands of rows. Performance depends on:
- Number of loops and iterations
- Complexity of expressions
- Available system memory

### Can I generate multiple sheets?

Yes, define multiple `<Sheet>` tags in your template.

## Development

### Can I use goxcel as a Go library?

Yes:

```go
import "github.com/ryo-arima/goxcel/pkg/controller"

conf := config.NewBaseConfigWithFile("template.gxl")
ctrl := controller.NewCommonController(conf)
err := ctrl.Generate("template.gxl", "data.json", "output.xlsx", false)
```

### How do I contribute?

See the GitHub repository for contribution guidelines.
