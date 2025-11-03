# goxcel

**Template-driven Excel generation library for Go**

Generate Excel (.xlsx) files from human-readable grid templates. Combines Markdown-like syntax with programmatic Excel generation power.

## Why goxcel?

- **Visual Templates**: Grid-oriented templates that look like your Excel output
- **Data-Driven**: Separate data from layout using JSON contexts
- **Pure Go**: No external dependencies, no C libraries
- **Type-Safe**: Strong typing with compile-time safety
- **Extensible**: Formulas, merges, images, charts, and more

## Features

✅ Grid templates with pipe-delimited syntax  
✅ Value interpolation (`{{ expr }}`)  
✅ Control structures (For loops, conditionals)  
✅ Excel formulas  
✅ Cell merging  
✅ Components (Images, Shapes, Charts)  
✅ Structured logging  
✅ CLI tool

## Quick Example

**Template (.gxl)**:
```xml
<Sheet name="Sales">
<Grid>| Product | Qty | Price |</Grid>
<For each="item in items">
<Grid>| {{ item.name }} | {{ item.qty }} | {{ item.price }} |</Grid>
</For>
</Sheet>
```

**Data (JSON)**:
```json
{"items": [{"name": "Apple", "qty": 10, "price": 100}]}
```

**Generate**:
```bash
goxcel generate --template report.gxl --data data.json --output report.xlsx
```

## Next Steps

- **New?** [Quick Start](./getting-started/quick-start.md)
- **Vision?** [Mission & Strategy](./vision-strategy.md)
- **Details?** [Specification](./specification/core-tags.md)
- **Contribute?** [Contributing](./development/contributing.md)

## Status & License

Active development. Core features stable. MIT License.

