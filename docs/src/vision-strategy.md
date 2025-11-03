# Vision & Strategy

## Mission

**Make Excel generation as simple as writing Markdown** - developers should visualize output directly in templates without deep OOXML knowledge.

## Problems Solved

**Traditional approaches:**
- Verbose cell-by-cell code
- Layout mixed with business logic
- Poor preview capability
- Non-technical users cannot contribute

**Our solution:**
- Visual grid templates
- Data/presentation separation
- Template-first approach

```xml
<Grid>
| Name | Quantity | Price |
</Grid>
<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.qty }} | {{ item.price }} |
</Grid>
</For>
```

**Benefits:**
- ✅ Visual structure matches Excel output
- ✅ Data and layout are separated
- ✅ Templates can be versioned and reused
- ✅ Non-programmers can understand templates
- ✅ Easy to review and validate

---

## Vision

### Our Vision for the Future

**To become the standard template language for structured document generation across all formats.**

While goxcel starts with Excel, we envision a future where:

1. **Universal Templates**: A single template syntax generates Excel, PDF, Word, and HTML
2. **Visual Editors**: WYSIWYG editors that generate GXL templates
## Core Values

1. **Simplicity First**: Readable code, minimal API, convention over configuration
2. **Developer Experience**: Clear errors, comprehensive docs, easy onboarding
3. **Reliability**: Backward compatibility, comprehensive tests, stable releases
4. **Performance**: Optimize hot paths, streaming for large data, predictable behavior
5. **Openness**: Open source, transparent decisions, MIT license
6. **Pragmatism**: Ship working software, validate with users, iterate

## Roadmap

### v1.0 (Current - Q4 2024) ✓
Core features: Parsing, Grid layout, Interpolation, For loops, Formulas, Merging, CLI

### v1.1 (Q1 2025)
Conditionals, Anchor positioning, Basic styling, Cell formatting, Validation

### v1.2 (Q2 2025)
Images, Charts, Shapes, Named ranges, Sheet protection, Hyperlinks

### v2.0 (Q3 2025)
Pivot tables, Advanced charts, Conditional formatting, Streaming mode, Multi-sheet refs

### v3.0+ (2026+)
Multi-format (PDF, HTML, Word), Visual editor, AI assistance, Enterprise features

## Strategy

**Phase 1 (Now)**: Excel excellence with clean architecture  
**Phase 2 (2025)**: Enhanced capabilities and ecosystem growth  
**Phase 3 (2026)**: Multi-format support and marketplace  
**Phase 4 (2027+)**: Intelligence layer and enterprise adoption

**Commitment**: SemVer, 2-year LTS, 6-month deprecation notice, migration guides

