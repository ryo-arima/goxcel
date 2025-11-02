# goxcel

**A template-driven Excel generation library for Go**

goxcel is a powerful yet simple library that generates Excel (.xlsx) files from human-readable grid templates. It combines the ease of Markdown-like syntax with the power of programmatic Excel generation.

## Why goxcel?

Traditional Excel generation libraries require you to write verbose code for each cell, row, and column. goxcel takes a different approach:

- **Visual Templates**: Write grid-oriented templates that look like the Excel sheets you want to create
- **Data-Driven**: Separate your data from your layout using JSON data contexts
- **Pure Go**: No external dependencies, no C libraries, just Go standard library
- **Type-Safe**: Strong typing with compile-time safety
- **Extensible**: Support for formulas, merges, images, charts, and more

## Key Features

- ✅ **Grid-based templates** with pipe-delimited syntax
- ✅ **Value interpolation** using `{{ expr }}` syntax
- ✅ **Control structures** (For loops, conditional rendering)
- ✅ **Excel formulas** with full expression support
- ✅ **Cell merging** for complex layouts
- ✅ **Component system** (Images, Shapes, Charts, Pivots)
- ✅ **Structured logging** with message codes
- ✅ **CLI tool** for quick generation

## Quick Example

**Template (.gxl)**:
```xml
<Book>
<Sheet name="Sales Report">
<Grid>
| Product | Quantity | Price |
</Grid>
<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.qty }} | {{ item.price }} |
</Grid>
</For>
<Grid>
| Total | | =SUM(C2:C4) |
</Grid>
</Sheet>
</Book>
```

**Data (JSON)**:
```json
{
  "items": [
    {"name": "Apple", "qty": 10, "price": 100},
    {"name": "Banana", "qty": 20, "price": 200}
  ]
}
```

**Generate**:
```bash
goxcel generate --template report.gxl --data data.json --output report.xlsx
```

## What's Next?

- **New to goxcel?** Start with the [Quick Start Guide](./getting-started/quick-start.md)
- **Want to understand the vision?** Read our [Mission](./vision/mission.md) and [Values](./vision/values.md)
- **Need detailed docs?** Check out the [Reference](./reference/gxl-spec.md)
- **Want to contribute?** See our [Contributing Guide](./development/contributing.md)

## Project Status

goxcel is in active development. Core features are stable and ready for use. See our [Roadmap](./vision/roadmap.md) for planned features.

## License

goxcel is released under the [MIT License](./appendix/license.md).
