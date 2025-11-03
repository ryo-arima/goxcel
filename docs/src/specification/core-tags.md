# Core Tags

Core tags are the fundamental building blocks of GXL templates. They define the structure of workbooks, sheets, and cell content.

---

## Book

The root element that defines an Excel workbook.

### Syntax

```xml
<Book>
  <!-- sheets and content -->
</Book>
```

### Attributes

Currently, `<Book>` has no required attributes. Future versions may support:
- `title`: Workbook title
- `author`: Document author
- `created`: Creation date
- `modified`: Last modification date

### Rules

1. **Must be root element**: Every GXL file must have exactly one `<Book>` element
2. **Must contain sheets**: At least one `<Sheet>` element is required
3. **No content outside**: No text or tags allowed outside `<Book>`

### Examples

**Minimal workbook:**
```xml
<Book>
  <Sheet name="Sheet1">
    <Grid>
    | Hello | World |
    </Grid>
  </Sheet>
</Book>
```

**Multiple sheets:**
```xml
<Book>
  <Sheet name="Sales">
    <Grid>
    | Product | Revenue |
    </Grid>
  </Sheet>
  
  <Sheet name="Expenses">
    <Grid>
    | Category | Amount |
    </Grid>
  </Sheet>
</Book>
```

### Future Features (Planned)

```xml
<Book
  title="Annual Report 2024"
  author="Finance Department"
  created="2024-01-01"
>
  <Metadata>
    <Property name="department" value="Finance" />
    <Property name="classification" value="Internal" />
  </Metadata>
  
  <Sheet name="Summary">
    ...
  </Sheet>
</Book>
```

---

## Sheet

Defines a worksheet within the workbook.

### Syntax

```xml
<Sheet name="SheetName">
  <!-- content -->
</Sheet>
```

### Attributes

#### `name` (required)
- **Type**: String
- **Description**: The name of the worksheet as it appears in Excel
- **Constraints**:
  - Must be unique within the workbook
  - Maximum 31 characters (Excel limitation)
  - Cannot contain: `\ / ? * [ ]` or `:`
  - Cannot be empty
  - Leading/trailing spaces are trimmed

#### `col_width` (optional)
- **Type**: Length
- **Description**: Default column width for the entire sheet
- **Units**: `ch` (characters, default if no unit), `cm`, `mm`, `in`, `pt`, `px`
- **Examples**: `"8.43"`, `"1cm"`, `"72px"`

#### `row_height` (optional)
- **Type**: Length
- **Description**: Default row height for the entire sheet
- **Units**: `pt` (points, default if no unit), `cm`, `mm`, `in`, `px`
- **Examples**: `"15"`, `"1cm"`, `"20px"`

Note: For backward compatibility, `row_heigh` is accepted as an alias of `row_height`.

### Rules

1. **Unique names**: No two sheets can have the same name
2. **At least one sheet**: A workbook must contain at least one sheet
3. **Order matters**: Sheets appear in the order defined
4. **Case-sensitive**: `"Sales"` and `"sales"` are different sheets

### Examples

**Basic sheet:**
```xml
<Sheet name="Sales Data" col_width="1cm" row_height="1cm">
  <Grid>
  | Date | Amount |
  | 2024-01-01 | 1000 |
  </Grid>
</Sheet>
```

**Multiple sheets with different purposes:**
```xml
<Book>
  <!-- Data sheet -->
  <Sheet name="Raw Data">
    <Grid>
    | ID | Name | Value |
    </Grid>
    <For each="row in data">
    <Grid>
    | {{ row.id }} | {{ row.name }} | {{ row.value }} |
    </Grid>
    </For>
  </Sheet>
  
  <!-- Summary sheet -->
  <Sheet name="Summary">
    <Grid>
    | Total Records | =COUNTA('Raw Data'!A:A)-1 |
    | Total Value | =SUM('Raw Data'!C:C) |
    </Grid>
  </Sheet>
  
  <!-- Charts sheet -->
  <Sheet name="Visualizations">
    <Chart 
      ref="A1" 
      type="column" 
      dataRange="'Raw Data'!A1:C100" 
      title="Data Overview"
    />
  </Sheet>
</Book>
```

### Naming Best Practices

**Good names:**
- Descriptive: `Sales 2024`, `Employee List`, `Profit & Loss`
- Concise: Keep under 20 characters when possible
- Clear purpose: Name indicates content

**Avoid:**
- Generic names: `Sheet1`, `Sheet2`
- Special characters: `Sales/Expenses` (use `Sales & Expenses`)
- Too long: `This is a very long sheet name that exceeds the limit`

### Sheet References in Formulas

When referencing cells from other sheets, use single quotes if sheet name contains spaces:

```xml
<Sheet name="Summary">
  <Grid>
  | Total from Sales | =SUM('Sales Data'!B:B) |
  </Grid>
</Sheet>
```

---

## Grid

Defines a grid of cells using pipe-delimited rows.

### Syntax

**Basic usage:**
```xml
<Grid>
| Cell A1 | Cell B1 | Cell C1 |
| Cell A2 | Cell B2 | Cell C2 |
| Cell A3 | Cell B3 | Cell C3 |
</Grid>
```

**With absolute positioning (v1.0+):**
```xml
<Grid ref="D5">
| Cell D5 | Cell E5 |
| Cell D6 | Cell E6 |
</Grid>
```

**With style attributes (v1.x+):**
```xml
<Grid font="Arial" font_size="12" text_color="#333333" fill_color="#FFFFCC">
| Header 1 | Header 2 |
| Data 1   | Data 2   |
</Grid>
```

### Attributes

#### `ref` (optional, v1.0+)
- **Type**: String (A1 notation)
- **Description**: Absolute starting position for the grid
- **Default**: Current cursor position
- **Examples**: `"A1"`, `"B5"`, `"D10"`
- **Behavior**: When specified, the grid is placed at the absolute position without affecting the cursor position

#### Style attributes (optional, v1.x+)
- `font` / `font_name`: Font family for all cells in the grid (e.g., `Arial`)
- `font_size` / `text_size`: Font size in points (integer)
- `font_color` / `text_color`: Font color in RGB hex; `#` optional (e.g., `#FF0000` or `FF0000`)
- `fill_color` / `color`: Background fill color in RGB hex; `#` optional
 - `border` / `border_style`: Border style for the grid's cells. Supported: `thin`, `medium`, `thick`, `dashed`, `dotted`, `double`
 - `border_color`: Border color in RGB hex; `#` optional
 - `border_sides`: Comma-separated sides to apply (default `all`). Options: `all`, `top`, `right`, `bottom`, `left`

These defaults apply to every cell produced by the Grid unless overridden by per-cell formatting (e.g., markdown `**bold**`).

**Examples (borders):**
```xml
<Grid border="thin" border_color="#999999">
| A | B |
| 1 | 2 |
</Grid>

<Grid border="dashed" border_color="#FF0000" border_sides="top,bottom">
| Header 1 | Header 2 |
| Data 1   | Data 2   |
</Grid>
```

**Example - Grid with ref:**
```xml
<!-- Sequential grids -->
<Grid>
| Header 1 | Header 2 |
</Grid>

<Grid>
| Data 1 | Data 2 |
</Grid>

<!-- Absolute position at E1 (doesn't affect cursor) -->
<Grid ref="E1">
| Side Note |
</Grid>

<!-- Continues from A3 (after the first two grids) -->
<Grid>
| Row 3 | More Data |
</Grid>
```

**Result:**
| | A | B | C | D | E |
|-|---|---|---|---|---|
| 1 | Header 1 | Header 2 | | | Side Note |
| 2 | Data 1 | Data 2 | | | |
| 3 | Row 3 | More Data | | | |


### Structure

- **Rows**: Each line within `<Grid>` represents one row
- **Columns**: Cells are delimited by `|` (pipe character)
- **Optional pipes**: Leading and trailing pipes are optional
- **Whitespace**: Trimmed around cell content

### Cell Content Types

#### 1. Literal Values

```xml
<Grid>
| Plain text | 123 | 45.67 | true |
</Grid>
```

Cell types are inferred:
- Numbers: `123`, `45.67`, `-10.5`
- Strings: `Hello`, `Product Name`
- Booleans: `true`, `false`
- Dates: `2024-01-01` (ISO format)

#### 2. Formulas

Cells starting with `=` are Excel formulas:

```xml
<Grid>
| Product | Price | Quantity | Total |
| Widget | 10.50 | 5 | =B2*C2 |
| Gadget | 25.00 | 3 | =B3*C3 |
| | | Grand Total | =SUM(D2:D3) |
</Grid>
```

Supported formula features:
- Cell references: `A1`, `B2`, `$A$1`
- Ranges: `A1:A10`, `B2:D5`
- Functions: `SUM()`, `AVERAGE()`, `IF()`, etc.
- Operators: `+`, `-`, `*`, `/`, `^`
- Sheet references: `'Sheet1'!A1`

#### 3. Interpolated Values

Use `{{ }}` for dynamic values:

```xml
<Grid>
| {{ user.name }} | {{ user.email }} | {{ user.age }} |
</Grid>
```

#### 4. Mixed (Formulas with Interpolation)

Combine formulas and expressions:

```xml
<Grid>
| Total | =SUM(A1:A{{ rowCount }}) |
| Average | =AVERAGE(B1:B{{ rowCount }}) |
</Grid>
```

### Empty Cells

Multiple consecutive pipes create empty cells:

```xml
<Grid>
| A | | C |  <!-- B is empty -->
| | B | |     <!-- A and C are empty -->
</Grid>
```

### Positioning

Grid cells are placed relative to the **current cursor position** (unless `ref` is specified):
- Starts at `A1` by default
- Advances after each `<Grid>` block
- Can be reset with `<Anchor>`

**Sequential positioning:**
```xml
<!-- Starts at A1 -->
<Grid>
| Row 1 |
</Grid>

<!-- Continues at A2 -->
<Grid>
| Row 2 |
</Grid>
```

**With Anchor:**
```xml
<!-- Reset to E1 -->
<Anchor ref="E1" />
<Grid>
| Over here |
</Grid>
```

**With Grid ref attribute (v1.0+):**
```xml
<!-- Normal sequential -->
<Grid>
| Row 1 |
</Grid>

<!-- Absolute position (cursor stays at A2) -->
<Grid ref="E1">
| Absolute |
</Grid>

<!-- Continues at A2 -->
<Grid>
| Row 2 |
</Grid>
```

**Comparison: Anchor vs Grid ref**

| Feature | `<Anchor>` | `<Grid ref="">` |
|---------|------------|-----------------|
| Scope | Affects all subsequent content | Only affects that grid |
| Cursor Movement | Moves cursor permanently | Doesn't affect cursor |
| Use Case | Change layout flow | Place independent content |


### Multi-column Grids

```xml
<Grid>
| Name | Email | Phone | Address |
| John | j@example.com | 555-1234 | 123 Main St |
| Jane | jane@example.com | 555-5678 | 456 Oak Ave |
</Grid>
```

### Column Alignment

Pipes don't need to align, but alignment improves readability:

```xml
<!-- Not aligned (valid but hard to read) -->
<Grid>
| Name | Email | Phone |
| John Doe | john@example.com | 555-1234 |
| Jane Smith | jane@example.com | 555-5678 |
</Grid>

<!-- Aligned (recommended) -->
<Grid>
| Name       | Email             | Phone    |
| John Doe   | john@example.com  | 555-1234 |
| Jane Smith | jane@example.com  | 555-5678 |
</Grid>
```

---

## Anchor

Sets the absolute position for subsequent grid placement.

### Syntax

```xml
<Anchor ref="CellReference" />
```

### Attributes

#### `ref` (required)
- **Type**: String
- **Description**: Absolute cell reference in A1 notation
- **Format**: `[Column][Row]` (e.g., `A1`, `Z100`, `AA5`)
- **Valid examples**: `A1`, `B10`, `AA1`, `XFD1048576`

### Purpose

By default, content flows from top-left (A1) downward. Use `<Anchor>` to:
1. Position content at specific locations
2. Create multiple independent sections
3. Layout complex reports

### Cursor Behavior

- **Initial position**: `A1` (if no anchor specified)
- **After Grid**: Cursor advances downward by number of rows
- **After Anchor**: Cursor jumps to specified position

### Examples

**Position content at specific cell:**
```xml
<Anchor ref="A1" />
<Grid>
| Header |
</Grid>

<Anchor ref="A10" />
<Grid>
| Footer |
</Grid>
```

**Create side-by-side sections:**
```xml
<!-- Left section -->
<Anchor ref="A1" />
<Grid>
| Section 1 |
| Data here |
</Grid>

<!-- Right section -->
<Anchor ref="E1" />
<Grid>
| Section 2 |
| Data here |
</Grid>
```

**Complex layout:**
```xml
<Sheet name="Dashboard">
  <!-- Title at top -->
  <Anchor ref="A1" />
  <Grid>
  | Sales Dashboard |
  </Grid>
  <Merge range="A1:F1" />
  
  <!-- KPIs in row 3 -->
  <Anchor ref="A3" />
  <Grid>
  | Revenue | $1,000,000 |
  | Orders | 500 |
  </Grid>
  
  <Anchor ref="D3" />
  <Grid>
  | Customers | 250 |
  | Conversion | 5% |
  </Grid>
  
  <!-- Chart at row 7 -->
  <Anchor ref="A7" />
  <Chart ref="A7" type="column" dataRange="A3:B4" />
</Sheet>
```

### Best Practices

1. **Use sparingly**: Let content flow naturally when possible
2. **Prefer Grid ref**: For independent content, use `<Grid ref="">` instead of `<Anchor>` to avoid affecting layout flow
3. **Document reasons**: Add comments explaining why specific positioning is needed
4. **Avoid overlaps**: Ensure anchored content doesn't overlap
5. **Test thoroughly**: Verify layout with different data sizes

**When to use Anchor vs Grid ref:**
- **Use `<Anchor>`**: When you want to change the layout flow permanently (e.g., switching from top section to side panel)
- **Use `<Grid ref="">`**: When you want to place independent content (e.g., logo, side notes) without affecting the main flow

---

## Merge

Merges a range of cells into a single cell.

### Syntax

```xml
<Merge range="StartCell:EndCell" />
```

### Attributes

#### `range` (required)
- **Type**: String
- **Description**: Cell range in A1 notation
- **Format**: `StartCell:EndCell` (e.g., `A1:C1`, `B2:D5`)
- **Examples**: `A1:C1` (horizontal), `A1:A3` (vertical), `B2:D4` (rectangular)

### Behavior

- **Content**: Only the top-left cell's content is visible
- **Other cells**: Content in other cells of the range is discarded
- **Formulas**: Can reference merged cells normally
- **Formatting**: Merge applies to entire range

### Examples

**Horizontal merge (title spanning columns):**
```xml
<Grid>
| Sales Report for Q4 2024 | | | |
| Region | Q1 | Q2 | Q3 | Q4 |
</Grid>
<Merge range="A1:E1" />
```

**Vertical merge (row headers):**
```xml
<Grid>
| Category | Product A | 100 |
| | Product B | 150 |
| | Product C | 200 |
</Grid>
<Merge range="A1:A3" />
```

**Rectangular merge:**
```xml
<Grid>
| Large Merged Area | | |
| | | |
| | | |
</Grid>
<Merge range="A1:C3" />
```

**Multiple merges:**
```xml
<Grid>
| Title | | | Date | |
| Section 1 | Data | Data | Section 2 | Data |
</Grid>
<Merge range="A1:C1" />  <!-- Title -->
<Merge range="D1:E1" />  <!-- Date -->
```

### Complex Example: Invoice Header

```xml
<Grid>
| Company Name | | | Invoice #12345 |
| 123 Main Street | | | Date: 2024-11-03 |
| City, State 12345 | | | Due: 2024-12-03 |
</Grid>
<Merge range="A1:C1" />  <!-- Company name -->
<Merge range="A2:C2" />  <!-- Address line 1 -->
<Merge range="A3:C3" />  <!-- Address line 2 -->
```

### Best Practices

1. **Use for headers**: Merge cells for titles and section headers
2. **Preserve alignment**: Consider how merged cells affect layout
3. **Document merges**: Add comments for complex merge patterns
4. **Test formulas**: Ensure formulas work with merged cells

### Common Patterns

**Report title:**
```xml
<Grid>
| Annual Sales Report | | | |
</Grid>
<Merge range="A1:D1" />
```

**Section headers:**
```xml
<Grid>
| Q1 Results | | |
| Jan | Feb | Mar |
</Grid>
<Merge range="A1:C1" />
```

**Grouped data:**
```xml
<Grid>
| Department | | Employee | Hours |
| Engineering | | John | 40 |
| | | Jane | 38 |
| | | Bob | 42 |
</Grid>
<Merge range="A2:B4" />  <!-- Department spans 3 rows -->
```

---

## Summary

Core tags provide the foundation for GXL templates:

- **`<Book>`**: Root element containing all sheets
- **`<Sheet>`**: Individual worksheets with unique names
- **`<Grid>`**: Pipe-delimited cell content
- **`<Anchor>`**: Position content at specific cells
- **`<Merge>`**: Combine cells into single merged cell

---

## Next Steps

- [Control Structures](./control-structures.md) - Learn about loops and conditionals
- [Expressions](./expressions.md) - Dynamic value interpolation
- [Examples](./examples.md) - See complete templates using core tags
