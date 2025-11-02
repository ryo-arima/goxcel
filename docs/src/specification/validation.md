# Validation Rules

This section defines the constraints and error conditions for GXL templates.

---

## Document Structure

### Root Element

**Rule:** Every GXL file must have exactly one `<Book>` root element.

**Valid:**
```xml
<Book>
  <Sheet name="Sheet1">
    ...
  </Sheet>
</Book>
```

**Invalid:**
```xml
<!-- Missing <Book> -->
<Sheet name="Sheet1">
  ...
</Sheet>
```

```xml
<!-- Multiple <Book> elements -->
<Book>
  <Sheet name="Sheet1">...</Sheet>
</Book>
<Book>
  <Sheet name="Sheet2">...</Sheet>
</Book>
```

---

### Required Sheets

**Rule:** At least one `<Sheet>` element is required within `<Book>`.

**Valid:**
```xml
<Book>
  <Sheet name="Sheet1">
    <Grid>| Data |</Grid>
  </Sheet>
</Book>
```

**Invalid:**
```xml
<Book>
  <!-- No sheets -->
</Book>
```

---

## Sheet Names

### Uniqueness

**Rule:** Sheet names must be unique within a workbook.

**Valid:**
```xml
<Book>
  <Sheet name="Sales">...</Sheet>
  <Sheet name="Expenses">...</Sheet>
  <Sheet name="Summary">...</Sheet>
</Book>
```

**Invalid:**
```xml
<Book>
  <Sheet name="Data">...</Sheet>
  <Sheet name="Data">...</Sheet>  <!-- Duplicate name -->
</Book>
```

---

### Length

**Rule:** Sheet names must not exceed 31 characters (Excel limitation).

**Valid:**
```xml
<Sheet name="Sales Report 2024">...</Sheet>  <!-- 18 characters -->
```

**Invalid:**
```xml
<Sheet name="This is a very long sheet name that exceeds the limit">...</Sheet>  <!-- 58 characters -->
```

---

### Special Characters

**Rule:** Sheet names cannot contain: `\ / ? * [ ] :`

**Valid:**
```xml
<Sheet name="Sales & Marketing">...</Sheet>
<Sheet name="Q1-Q2 Comparison">...</Sheet>
<Sheet name="Profit (Net)">...</Sheet>
```

**Invalid:**
```xml
<Sheet name="Sales/Marketing">...</Sheet>  <!-- Contains / -->
<Sheet name="Q1:Q2">...</Sheet>  <!-- Contains : -->
<Sheet name="Data[1]">...</Sheet>  <!-- Contains [ ] -->
```

---

### Empty Names

**Rule:** Sheet names cannot be empty.

**Invalid:**
```xml
<Sheet name="">...</Sheet>
<Sheet>...</Sheet>  <!-- Missing name attribute -->
```

---

## Cell References

### A1 Notation

**Rule:** Cell references must use valid A1 notation.

**Format:** `[Column][Row]`
- Column: A-Z, AA-ZZ, AAA-XFD
- Row: 1-1048576

**Valid:**
- `A1`
- `Z100`
- `AA1`
- `XFD1048576`

**Invalid:**
- `1A` (row before column)
- `A0` (row must be >= 1)
- `A1048577` (exceeds row limit)
- `XFE1` (exceeds column limit)

---

### Range Notation

**Rule:** Ranges must use format `StartCell:EndCell`.

**Valid:**
- `A1:C10`
- `B2:D5`
- `A1:XFD1048576`

**Invalid:**
- `A1-C10` (use `:` not `-`)
- `C10:A1` (start must come before end)
- `A1:` (incomplete range)

---

## Tag Structure

### Properly Nested

**Rule:** Tags must be properly nested (no overlapping).

**Valid:**
```xml
<Book>
  <Sheet name="Sheet1">
    <For each="item in items">
      <Grid>
      | {{ item.name }} |
      </Grid>
    </For>
  </Sheet>
</Book>
```

**Invalid:**
```xml
<Book>
  <Sheet name="Sheet1">
    <For each="item in items">
      <Grid>
      | {{ item.name }} |
  </Sheet>
    </For>  <!-- Overlapping with </Sheet> -->
</Book>
```

---

### Closed Tags

**Rule:** All opening tags must have matching closing tags.

**Valid:**
```xml
<Book>
  <Sheet name="Sheet1">
    <Grid>| Data |</Grid>
  </Sheet>
</Book>
```

**Invalid:**
```xml
<Book>
  <Sheet name="Sheet1">
    <Grid>| Data |</Grid>
  <!-- Missing </Sheet> -->
</Book>
```

---

### Self-Closing Tags

**Rule:** Self-closing tags must end with `/>`.

**Valid:**
```xml
<Anchor ref="A1" />
<Merge range="A1:C1" />
<Image ref="B3" src="logo.png" />
```

**Invalid:**
```xml
<Anchor ref="A1">  <!-- Should be self-closing -->
<Merge range="A1:C1">  <!-- Should be self-closing -->
```

---

## Attributes

### Required Attributes

**Rule:** Required attributes must be present.

| Tag | Required Attributes |
|-----|---------------------|
| `<Sheet>` | `name` |
| `<Anchor>` | `ref` |
| `<Merge>` | `range` |
| `<For>` | `each` |
| `<If>` | `cond` |
| `<Image>` | `ref`, `src` |
| `<Shape>` | `ref`, `kind` |
| `<Chart>` | `ref`, `type`, `dataRange` |
| `<Pivot>` | `ref`, `sourceRange`, `values` |

**Invalid:**
```xml
<Sheet>...</Sheet>  <!-- Missing name -->
<Anchor />  <!-- Missing ref -->
<For>...</For>  <!-- Missing each -->
```

---

### Attribute Syntax

**Rule:** Attributes must use format `name="value"` with quotes.

**Valid:**
```xml
<Sheet name="Sales">
<Anchor ref="A1" />
```

**Invalid:**
```xml
<Sheet name=Sales>  <!-- Missing quotes -->
<Sheet name='Sales'>  <!-- Use double quotes -->
<Anchor ref=A1 />  <!-- Missing quotes -->
```

---

### Unknown Attributes

**Rule:** Unknown attributes generate warnings (but don't fail).

```xml
<Sheet name="Sales" unknownAttr="value">  <!-- Warning: unknownAttr -->
  ...
</Sheet>
```

---

## Expressions

### Balanced Braces

**Rule:** Expression braces must be balanced.

**Valid:**
```xml
| {{ value }} |
| {{ user.name }} |
| {{ items[0].price }} |
```

**Invalid:**
```xml
| {{ value |  <!-- Missing closing }} -->
| { value }} |  <!-- Missing opening { -->
| {{ value } |  <!-- Unbalanced -->
```

---

### Valid Paths

**Rule:** Expression paths must use valid identifier syntax.

**Valid:**
- `{{ user }}`
- `{{ user.name }}`
- `{{ items[0] }}`
- `{{ data.nested.property }}`

**Invalid:**
- `{{ user.123 }}` (identifiers can't start with number)
- `{{ user..name }}` (consecutive dots)
- `{{ user. }}` (trailing dot)

---

## Control Structures

### For Loop Syntax

**Rule:** `each` attribute must follow format `<var> in <path>`.

**Valid:**
```xml
<For each="item in items">
<For each="user in users">
<For each="row in data.rows">
```

**Invalid:**
```xml
<For each="item">  <!-- Missing 'in' clause -->
<For each="in items">  <!-- Missing variable -->
<For each="item of items">  <!-- Use 'in' not 'of' -->
```

---

### If Condition

**Rule:** `cond` attribute must contain a valid expression.

**Valid:**
```xml
<If cond="isActive">
<If cond="total > 1000">
<If cond="status == 'paid'">
```

**Invalid:**
```xml
<If cond="">  <!-- Empty condition -->
<If>  <!-- Missing cond attribute -->
```

---

## Grid Structure

### Pipe Delimiters

**Rule:** Grid rows must use `|` to delimit cells.

**Valid:**
```xml
<Grid>
| A | B | C |
| 1 | 2 | 3 |
</Grid>
```

**Acceptable (leading/trailing pipes optional):**
```xml
<Grid>
A | B | C
1 | 2 | 3
</Grid>
```

**Warning (inconsistent):**
```xml
<Grid>
| A | B | C |
A | B | C  <!-- Missing trailing pipe (acceptable but inconsistent) -->
</Grid>
```

---

## Merge Ranges

### Valid Ranges

**Rule:** Merge ranges must be rectangular and valid.

**Valid:**
```xml
<Merge range="A1:C1" />  <!-- Horizontal -->
<Merge range="A1:A3" />  <!-- Vertical -->
<Merge range="B2:D4" />  <!-- Rectangular -->
```

**Invalid:**
```xml
<Merge range="A1:A1" />  <!-- Single cell (no merge needed) -->
<Merge range="C1:A1" />  <!-- End before start -->
```

---

## Component Positioning

### No Overlaps (Recommended)

**Rule:** Components should not overlap data cells.

**Warning (overlap):**
```xml
<Grid>
| A | B | C | D |
| 1 | 2 | 3 | 4 |
</Grid>

<!-- Chart overlaps cells C1:D2 -->
<Chart ref="C1" type="column" dataRange="A1:B2" width="200" height="100" />
```

**Better:**
```xml
<Grid>
| A | B |
| 1 | 2 |
</Grid>

<!-- Chart in separate area -->
<Chart ref="D1" type="column" dataRange="A1:B2" width="200" height="100" />
```

---

## Validation Levels

### Error (Parsing Fails)

These violations prevent template parsing:
- Missing `<Book>` root element
- Unclosed tags
- Malformed attributes
- Invalid tag nesting

### Warning (Logged but Continues)

These violations generate warnings:
- Unknown attributes
- Unknown tags (future extensibility)
- Empty arrays in `<For>` loops
- Undefined variables in expressions

### Info (Best Practice)

These are recommendations:
- Component overlaps
- Very long sheet names (< 31 but > 20)
- Deeply nested loops (> 3 levels)
- Large data ranges (> 10,000 rows)

---

## Validation Tools

### Command-Line Validation (Planned)

```bash
goxcel validate template.gxl
goxcel validate --strict template.gxl
goxcel validate --schema schema.json template.gxl
```

### Programmatic Validation (Planned)

```go
import "github.com/ryo-arima/goxcel/pkg/validator"

result := validator.Validate("template.gxl")
if !result.IsValid {
    for _, err := range result.Errors {
        fmt.Println(err.Message)
    }
}
```

---

## Error Messages

### Good Error Messages

Error messages should include:
- **What**: Description of the problem
- **Where**: File location (line/column if possible)
- **Why**: Explanation of the rule violated
- **How**: Suggestion for fixing

**Example:**
```
Error at line 15, column 10:
  <Sheet name="Sales/Marketing">
                    ^
Sheet name contains invalid character '/'.
Sheet names cannot contain: \ / ? * [ ] :
Suggestion: Use 'Sales & Marketing' or 'Sales-Marketing'
```

---

## Related Topics

- [File Format](./file-format.md) - File structure requirements
- [Core Tags](./core-tags.md) - Tag syntax and attributes
- [Expressions](./expressions.md) - Expression syntax rules
