# Vision & Strategy

## Mission

### Our Mission

**To democratize Excel generation by making it as simple as writing Markdown.**

We believe that generating Excel files should not require extensive knowledge of complex APIs or deep understanding of the OOXML specification. Instead, developers should be able to:

1. **Visualize**: See the structure of their Excel output directly in the template
2. **Iterate**: Quickly modify layouts without rebuilding entire codebases
3. **Separate**: Keep data and presentation concerns cleanly separated
4. **Collaborate**: Enable non-programmers to understand and modify templates

### The Problem We Solve

Traditional Excel generation approaches suffer from several issues:

#### Verbose Code
```go
// Traditional approach - hard to visualize the output
sheet.SetCellValue("A1", "Name")
sheet.SetCellValue("B1", "Quantity")
sheet.SetCellValue("C1", "Price")
for i, item := range items {
    row := i + 2
    sheet.SetCellValue(fmt.Sprintf("A%d", row), item.Name)
    sheet.SetCellValue(fmt.Sprintf("B%d", row), item.Qty)
    sheet.SetCellValue(fmt.Sprintf("C%d", row), item.Price)
}
```

#### Tight Coupling
- Layout logic mixed with business logic
- Difficult to reuse templates across projects
- Hard to maintain when requirements change

#### Poor Discoverability
- No way to preview the structure without running code
- Difficult to review in code reviews
- Non-technical stakeholders cannot contribute

### Our Solution

goxcel provides a **template-first** approach:

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
3. **Template Marketplace**: A community-driven repository of reusable templates
4. **Enterprise Integration**: Seamless integration with BI tools and data platforms
5. **AI Assistance**: AI-powered template generation and optimization

### The Future We're Building

#### Phase 1: Excel Excellence (Current)
- Perfect the Excel generation experience
- Build a solid foundation with clean architecture
- Establish the GXL template syntax standard
- Grow an engaged community

#### Phase 2: Enhanced Capabilities
- Advanced styling and formatting
- Rich component library (charts, images, shapes)
- Template validation and error reporting
- Performance optimizations for large datasets

#### Phase 3: Multi-Format Support
- PDF generation from GXL templates
- Word document generation
- HTML/web output for previewing
- Format-agnostic template design

#### Phase 4: Ecosystem Growth
- Visual template editor (web-based)
- Template marketplace and sharing platform
- Integration plugins for popular frameworks
- Enterprise features (SSO, audit logs, governance)

#### Phase 5: Intelligence Layer
- AI-powered template suggestions
- Automated data-to-template mapping
- Performance optimization recommendations
- Natural language to template generation

### Long-term Goals

By 2030, we aim to:

- **100K+ Active Users**: Developers across industries using goxcel
- **10K+ Templates**: Community-contributed templates for common use cases
- **Enterprise Adoption**: Fortune 500 companies standardizing on GXL
- **Education**: Universities teaching document generation with goxcel
- **Certification**: Professional certification program for GXL experts

---

## Values

### Core Values

Our development philosophy is guided by these principles:

#### 1. Simplicity First
**Complex problems deserve simple solutions.**

- Prefer readable code over clever optimizations
- Keep the API surface minimal and intuitive
- Choose convention over configuration
- Remove features that add complexity without value

*Example*: Grid syntax uses pipes (`|`) because they're visually intuitive, not because they're the most efficient parser tokens.

#### 2. Developer Experience
**Happy developers build better software.**

- Optimize for the first-time user experience
- Provide clear, actionable error messages
- Include comprehensive documentation and examples
- Support common workflows out-of-the-box

*Example*: Logger outputs include message codes (MCode) for easy troubleshooting and documentation lookup.

#### 3. Reliability
**Trust is earned through consistency.**

- Maintain backward compatibility
- Write comprehensive tests
- Document breaking changes clearly
- Support Long-Term Support (LTS) versions

*Example*: Using only Go standard library ensures long-term stability and reduces dependency risks.

#### 4. Performance
**Speed matters, but not at the cost of clarity.**

- Optimize hot paths without sacrificing readability
- Provide streaming APIs for large datasets
- Benchmark regularly and prevent regressions
- Make performance characteristics predictable

*Example*: Pure Go implementation avoids CGO overhead while maintaining clean code.

#### 5. Openness
**Great software is built in the open.**

- Open source from day one
- Welcome contributions from everyone
- Transparent decision-making process
- Clear governance model

*Example*: MIT license allows commercial use without restrictions.

#### 6. Pragmatism
**Ship working software over perfect designs.**

- Deliver incremental value frequently
- Validate assumptions with real users
- Refactor based on actual usage patterns
- Balance idealism with practical constraints

*Example*: Component placeholders (Image, Chart) are implemented first, full rendering comes later based on demand.

### Design Principles

#### Principle of Least Surprise
Users should be able to predict behavior from reading the template. No magic, no hidden behaviors.

#### Progressive Disclosure
Start simple, add complexity only when needed. Basic use cases should be trivial, advanced use cases should be possible.

#### Explicit Over Implicit
Be explicit about what's happening. Implicit behavior leads to confusion and bugs.

#### Composability
Small pieces that combine well are better than large monolithic features.

---

## Roadmap

### Version 1.0 (Current - Q4 2024)

**Goal**: Stable foundation for Excel generation

- [x] Core template parsing (GXL format)
- [x] Grid-based cell layout
- [x] Value interpolation (`{{ expr }}`)
- [x] For loops
- [x] Excel formula support
- [x] Cell merging
- [x] Component placeholders
- [x] Structured logging
- [x] CLI tool
- [x] Comprehensive documentation

### Version 1.1 (Q1 2025)

**Goal**: Enhanced control flow and styling

- [ ] Conditional rendering (`<If>` statements)
- [ ] Anchor positioning system
- [ ] Basic style system (colors, fonts, alignment)
- [ ] Cell borders and formatting
- [ ] Number formatting (currency, dates, percentages)
- [ ] Data validation
- [ ] Template validation and better error messages

### Version 1.2 (Q2 2025)

**Goal**: Rich components and advanced features

- [ ] Full image rendering (PNG, JPEG)
- [ ] Basic chart rendering (column, bar, line, pie)
- [ ] Shape rendering (rectangles, arrows)
- [ ] Named ranges
- [ ] Workbook-level settings
- [ ] Sheet protection
- [ ] Hyperlinks

### Version 2.0 (Q3 2025)

**Goal**: Enterprise-ready features

- [ ] Pivot table generation
- [ ] Advanced chart types (scatter, combo, stock)
- [ ] Conditional formatting
- [ ] Data tables and what-if analysis
- [ ] Custom functions
- [ ] Template inheritance and composition
- [ ] Streaming mode for large datasets
- [ ] Performance optimizations
- [ ] Multi-sheet references

### Version 2.1 (Q4 2025)

**Goal**: Developer tooling

- [ ] Template linter
- [ ] VS Code extension (syntax highlighting)
- [ ] Template debugger
- [ ] Performance profiler
- [ ] Template testing framework
- [ ] Migration tools
- [ ] Template gallery website

### Version 3.0 (2026)

**Goal**: Multi-format support

- [ ] PDF generation from GXL templates
- [ ] HTML preview mode
- [ ] Word document generation
- [ ] CSV export
- [ ] Format-agnostic template design
- [ ] Unified styling system across formats

### Version 4.0 (2027+)

**Goal**: Ecosystem and intelligence

- [ ] Visual template editor (web-based)
- [ ] Template marketplace
- [ ] AI-powered template generation
- [ ] Natural language queries
- [ ] Integration plugins (Spring Boot, Django, Rails)
- [ ] Cloud service (hosted rendering)
- [ ] Enterprise features (SSO, RBAC, audit logs)

### Community Milestones

- **Q4 2024**: 100 GitHub stars
- **Q2 2025**: 500 stars, 10 contributors
- **Q4 2025**: 1,000 stars, 50 contributors, first conference talk
- **2026**: 5,000 stars, 100+ contributors, enterprise adoption
- **2027**: Template marketplace launch
- **2030**: 10,000+ active users, industry standard

### How to Influence the Roadmap

We welcome community input! Here's how you can help shape goxcel's future:

1. **GitHub Issues**: Request features or report bugs
2. **Discussions**: Share use cases and requirements
3. **Pull Requests**: Contribute code or documentation
4. **Voting**: Upvote issues that matter to you
5. **Sponsorship**: Financial support accelerates development

### Commitment to Stability

- **Semantic Versioning**: We follow SemVer strictly
- **LTS Versions**: Major versions supported for 2 years
- **Deprecation Policy**: 6-month notice before removing features
- **Migration Guides**: Detailed guides for breaking changes
- **Backward Compatibility**: Maintained within major versions
