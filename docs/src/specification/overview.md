# GXL Format Specification

**Version:** 0.1 (Draft)  
**Status:** Evolving  
**Last Updated:** 2024-11-03

---

## Introduction

GXL (Grid eXcel Language) is a template format for describing Excel workbooks using a human-readable, Markdown-compatible syntax with custom XML-like tags. It enables developers to create Excel files by writing templates that visually represent the final spreadsheet output.

### What is GXL?

GXL stands for **Grid eXcel Language**. It is:

- **A template format**: Not a programming language, but a declarative template system
- **Human-readable**: Designed to be easy to read and write by humans
- **Grid-oriented**: Uses pipe-delimited syntax to represent tabular data visually
- **Markdown-compatible**: Can coexist with regular Markdown documentation
- **XML-inspired**: Uses custom tags similar to HTML/XML for structure

### Design Philosophy

GXL is built on these core principles:

1. **What You See Is What You Get**: The template structure should closely match the Excel output
2. **Separation of Concerns**: Data and presentation should be cleanly separated
3. **Progressive Enhancement**: Simple cases should be simple, complex cases should be possible
4. **Human-First**: Optimize for human readability over parser efficiency
5. **Extensibility**: Design for future features without breaking existing templates

---

## Who Should Use GXL?

### Target Audience

- **Backend Developers**: Building reporting systems, data exports, invoice generators
- **Data Engineers**: Creating data pipelines with Excel outputs
- **Full-Stack Developers**: Adding Excel export features to web applications
- **DevOps Engineers**: Generating operational reports and dashboards
- **Technical Writers**: Creating documentation with embedded data tables

### Use Cases

- **Reports**: Financial reports, sales reports, analytics dashboards
- **Invoices**: Customer invoices, purchase orders, receipts
- **Data Exports**: Database exports, API response dumps, log summaries
- **Templates**: Reusable document templates with variable data
- **Bulk Operations**: Mass generation of personalized documents

---

## Specification Structure

This specification is organized into the following sections:

1. **[Overview](./overview.md)** - High-level introduction and concepts
2. **[File Format](./file-format.md)** - File structure, encoding, and metadata
3. **[Core Tags](./core-tags.md)** - Book, Sheet, Grid, Anchor, Merge
4. **[Control Structures](./control-structures.md)** - For loops, If/Else conditionals
5. **[Expressions](./expressions.md)** - Value interpolation and expression language
6. **[Components](./components.md)** - Images, Shapes, Charts, Pivot Tables
7. **[Styling](./styling.md)** - Style system and formatting
8. **[Data Context](./data-context.md)** - How data flows through templates
9. **[Validation Rules](./validation.md)** - Constraints and error conditions
10. **[Rendering Semantics](./rendering.md)** - How templates are processed
11. **[Examples](./examples.md)** - Complete working examples

---

## Quick Example

Here's a minimal GXL template to get you started:

**Template (invoice.gxl):**
```xml
<Book>
<Sheet name="Invoice">
<Grid>
| Invoice #{{ invoiceNumber }} | Date: {{ date }} |
</Grid>

<Grid>
| Item | Quantity | Price | Total |
</Grid>

<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.qty }} | ${{ item.price }} | ={{ item.qty }}*{{ item.price }} |
</Grid>
</For>

<Grid>
| | | Total: | =SUM(D3:D{{ items.length + 2 }}) |
</Grid>
</Sheet>
</Book>
```

**Data (data.json):**
```json
{
  "invoiceNumber": "INV-2024-001",
  "date": "2024-11-03",
  "items": [
    {"name": "Widget A", "qty": 10, "price": 25.00},
    {"name": "Widget B", "qty": 5, "price": 50.00}
  ]
}
```

**Output:**
An Excel file with:
- Header row with invoice number and date
- Table with item details
- Automatic total calculation using Excel formula

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 (Draft) | 2024-11-03 | Initial specification |

---

## Compatibility

### Excel Version Support
- Excel 2007+ (.xlsx format)
- Office Open XML (OOXML) specification compliance

### Implementation Compatibility
- **goxcel v1.0.x**: Implements GXL 0.1 (Draft)
- Future implementations should follow this specification

---

## Feedback and Contributions

This specification is evolving based on community feedback. We welcome:

- Feature requests and suggestions
- Clarification questions
- Bug reports in specification
- Real-world use case examples

Please contribute via:
- GitHub Issues: Report problems or request features
- GitHub Discussions: Ask questions and share ideas
- Pull Requests: Propose changes to specification

---

## License

This specification is released under the Creative Commons Attribution 4.0 International License (CC BY 4.0).

You are free to:
- **Share**: Copy and redistribute the specification
- **Adapt**: Remix, transform, and build upon the specification

Under the following terms:
- **Attribution**: You must give appropriate credit

---

## Next Steps

Continue reading:
- [File Format Specification](./file-format.md) - Detailed file structure
- [Core Tags](./core-tags.md) - Learn about the fundamental tags
- [Examples](./examples.md) - See complete working examples
