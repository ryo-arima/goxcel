# File Format

This section describes the physical structure of GXL template files.

---

## File Extension

### Primary Extension
- **`.gxl`** - GXL template file

### Alternative Extensions (Future)
- `.gxl.md` - GXL with Markdown emphasis
- `.gxl.xml` - GXL with XML emphasis

---

## Character Encoding

### Required Encoding
- **UTF-8** without BOM (Byte Order Mark)

### Why UTF-8?
- Universal character support (international characters, emoji, etc.)
- Backward compatible with ASCII
- Industry standard for text files
- Git-friendly (no encoding issues)

### BOM Policy
- **Do NOT use BOM** (Byte Order Mark)
- BOM causes parsing issues in many tools
- UTF-8 BOM is unnecessary and discouraged

### Example
```
# Correct UTF-8 file (no BOM)
<Book>
<Sheet name="日本語シート">
...
```

---

## Line Endings

### Recommended
- **LF (Line Feed)** `\n` - Unix/Linux/macOS style

### Also Supported
- **CRLF (Carriage Return + Line Feed)** `\r\n` - Windows style

### Rationale
- LF is the Git default
- LF works across all platforms
- Most modern editors handle both automatically

### Git Configuration
Configure `.gitattributes` to normalize line endings:
```
*.gxl text eol=lf
```

---

## File Structure

### Overall Organization

A GXL file consists of:
1. Optional header comments
2. One `<Book>` root element
3. One or more `<Sheet>` elements
4. Content within sheets (Grid, components, control structures)

### Visual Structure

```
┌─────────────────────────────────┐
│ <!-- Optional Comments -->      │
├─────────────────────────────────┤
│ <Book>                          │
│   ├─ <Sheet name="Sheet1">      │
│   │    ├─ <Grid> ... </Grid>    │
│   │    ├─ <For> ... </For>      │
│   │    └─ <Merge ... />         │
│   │                              │
│   └─ <Sheet name="Sheet2">      │
│        └─ <Grid> ... </Grid>    │
│ </Book>                          │
└─────────────────────────────────┘
```

---

## Comments

### XML-Style Comments

```xml
<!-- This is a comment -->

<!--
  Multi-line comment
  can span multiple lines
-->
```

### Comment Rules
- Comments are **ignored** during parsing
- Can appear anywhere outside tags
- Cannot appear inside tag attributes
- Cannot be nested

### Use Cases
```xml
<!-- Invoice Template v2.1 -->
<!-- Author: John Doe -->
<!-- Last Modified: 2024-11-03 -->

<Book>
<!-- Sales Data Sheet -->
<Sheet name="Sales">
  <!-- Header row with company logo -->
  <Grid>
  | Company Name | Date |
  </Grid>
  
  <!-- TODO: Add monthly breakdown -->
</Sheet>
</Book>
```

---

## Whitespace Handling

### General Rules
- **Significant**: Whitespace inside `<Grid>` cells
- **Insignificant**: Whitespace outside tags and around tag names
- **Trimmed**: Leading/trailing whitespace in cell content

### Cell Content Trimming

```xml
<Grid>
|   Hello   |   World   |
</Grid>
```

Equivalent to:
```xml
<Grid>
| Hello | World |
</Grid>
```

Both produce cells with content `"Hello"` and `"World"` (no extra spaces).

### Preserving Whitespace

To preserve leading/trailing spaces, use expressions:
```xml
<Grid>
| {{ "  Hello  " }} |
</Grid>
```

### Indentation

Indentation is **insignificant** and used for readability:

```xml
<Book>
  <Sheet name="Example">
    <Grid>
    | A | B |
    </Grid>
  </Sheet>
</Book>
```

Same as:
```xml
<Book>
<Sheet name="Example">
<Grid>
| A | B |
</Grid>
</Sheet>
</Book>
```

---

## Automatic Formatting (CLI)

goxcel includes a template formatter available via the CLI:

```bash
goxcel format <file.gxl>           # pretty-print to stdout
goxcel format -w <file.gxl>        # overwrite in place
goxcel format -o out.gxl <file.gxl>
```

Rules applied by the formatter:

- Indentation: tags are indented with 2 spaces per nesting level
- Empty elements: when a tag has no text and no children, it is inlined on one line
  - Example: `<Merge range="A1:C1"> </Merge>`
- Blank lines: consecutive blank lines outside content are collapsed (no double blank lines)
- Grid alignment: inside `<Grid>`, pipe-delimited rows are aligned so that `|` characters line up by column
- Preservation: non-whitespace character data and XML comments are preserved

Before and After:

Before
```xml
<Grid>

  | A |  B   |C|
  | 1| 22 |333|

</Grid>
<Merge range="A1:C1">
</Merge>
```

After
```xml
<Grid>
  | A | B  | C   |
  | 1 | 22 | 333 |
</Grid>
<Merge range="A1:C1"> </Merge>
```

Note: Grid alignment uses character count (runes) for width; full-width East Asian display widths are not accounted for.

## Case Sensitivity

### Tag Names
- **Case-sensitive**: `<Book>` ≠ `<book>`
- **Convention**: PascalCase (e.g., `<Book>`, `<Sheet>`, `<Grid>`)

### Attribute Names
- **Case-sensitive**: `name="Sheet1"` ≠ `Name="Sheet1"`
- **Convention**: camelCase (e.g., `name`, `dataRange`, `fillColor`)

### Attribute Values
- **Case-sensitive**: Depend on context
  - Sheet names: `"Sales"` ≠ `"sales"`
  - Cell references: `"A1"` = `"a1"` (normalized to uppercase)
  - Expressions: `{{ user.Name }}` ≠ `{{ user.name }}`

---

## File Size Limits

### Practical Limits
- **File size**: No hard limit (limited by available memory)
- **Sheets**: Recommended max 100 sheets per workbook
- **Rows per sheet**: Excel limit is 1,048,576
- **Columns per sheet**: Excel limit is 16,384 (XFD)

### Performance Considerations
- Large templates (>10MB) may be slow to parse
- Use streaming mode for large datasets (future feature)
- Consider splitting large workbooks into multiple files

---

## MIME Type

### Proposed MIME Type
- **`text/x-gxl`** - GXL template file

### HTTP Headers
```
Content-Type: text/x-gxl; charset=utf-8
Content-Disposition: attachment; filename="template.gxl"
```

### File Upload Detection
Web servers and applications should recognize `.gxl` extension:
```apache
# Apache .htaccess
AddType text/x-gxl .gxl
```

---

## Metadata (Future)

### Embedded Metadata (Planned)

```xml
<Book
  title="Sales Report"
  author="John Doe"
  version="1.0"
  created="2024-11-03"
>
  <Metadata>
    <Property name="department" value="Finance" />
    <Property name="confidential" value="true" />
  </Metadata>
  
  <Sheet name="Data">
    ...
  </Sheet>
</Book>
```

**Status**: Planned for GXL v0.2

---

## File Naming Conventions

### Recommended Naming
- Use descriptive names: `invoice-template.gxl`, `sales-report.gxl`
- Use kebab-case: `monthly-report.gxl` (not `MonthlyReport.gxl`)
- Avoid spaces: Use hyphens or underscores
- Be specific: Include purpose in name

### Version Suffixes
```
invoice-template-v1.gxl
invoice-template-v2.gxl
sales-report-2024.gxl
```

### Environment Suffixes
```
report-dev.gxl
report-staging.gxl
report-prod.gxl
```

---

## Validation

### Well-Formedness
A valid GXL file must be:
- Valid UTF-8 encoding
- Properly nested XML-like tags
- One root `<Book>` element
- At least one `<Sheet>` element

### Validation Tools (Future)
- `gxl-lint`: Syntax checker
- `gxl-format`: Auto-formatter
- Editor plugins: VS Code, Sublime, etc.

---

## Best Practices

### 1. Use Version Control
- Store `.gxl` files in Git
- Track changes with meaningful commit messages
- Use branches for template variations

### 2. Document Templates
- Add comments explaining complex logic
- Include authorship and version info
- Document expected data structure

### 3. Organize Large Templates
- Group related content with comments
- Use consistent indentation
- Split very large templates into multiple files (via includes, future feature)

### 4. Test Templates
- Test with sample data
- Validate generated Excel files
- Check for edge cases (empty arrays, null values)

---

## Example Template Structure

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!-- 
  Invoice Template v1.0
  Author: John Doe
  Created: 2024-11-03
  
  Required data structure:
  {
    "company": {"name": string, "address": string},
    "invoice": {"number": string, "date": string},
    "items": [{"name": string, "qty": number, "price": number}]
  }
-->

<Book>
  <!-- Invoice Sheet -->
  <Sheet name="Invoice">
    <!-- Company Header -->
    <Grid>
    | {{ company.name }} | |
    | {{ company.address }} | |
    </Grid>
    <Merge range="A1:B1" />
    <Merge range="A2:B2" />
    
    <!-- Invoice Details -->
    <Grid>
    | Invoice: {{ invoice.number }} | Date: {{ invoice.date }} |
    </Grid>
    
    <!-- Item Table -->
    <Grid>
    | Item | Qty | Price | Total |
    </Grid>
    
    <For each="item in items">
    <Grid>
    | {{ item.name }} | {{ item.qty }} | {{ item.price }} | ={{ item.qty }}*{{ item.price }} |
    </Grid>
    </For>
    
    <!-- Total Row -->
    <Grid>
    | | | Total: | =SUM(D6:D{{ 5 + items.length }}) |
    </Grid>
  </Sheet>
</Book>
```

---

## Related Sections

- [Core Tags](./core-tags.md) - Learn about `<Book>`, `<Sheet>`, `<Grid>`
- [Expressions](./expressions.md) - Value interpolation syntax
- [Validation Rules](./validation.md) - What makes a valid GXL file
