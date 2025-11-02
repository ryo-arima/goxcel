# Rendering Semantics

This document describes how GXL templates are **processed** and **rendered** into Excel workbooks.

---

## Rendering Overview

**Rendering** is the process of transforming a `.gxl` template + JSON data into an Excel workbook.

### Key Concepts

1. **Template Parsing:** Parse `.gxl` text into an internal structure
2. **Data Binding:** Merge JSON context into expressions
3. **Loop Expansion:** Generate repeated content from arrays
4. **Grid Placement:** Convert pipe-delimited grids into cells
5. **Component Insertion:** Add images, charts, etc.
6. **Excel Generation:** Write final `.xlsx` file

---

## Execution Phases

### 1. Parse Phase

**Input:** `.gxl` text file  
**Output:** Abstract Syntax Tree (AST)

**Process:**
- Read file line by line
- Identify tags: `<Book>`, `<Sheet>`, `<Grid>`, `<For>`, etc.
- Parse attributes: `name="Sheet1"`, `src="data.items"`, etc.
- Build hierarchical tree structure
- Validate syntax (tag nesting, attribute format)

**Errors:**
- Unclosed tags
- Invalid tag names
- Missing required attributes

---

### 2. Render Phase

**Input:** AST + JSON data context  
**Output:** Expanded cell data + components

**Process:**
- Walk AST tree depth-first
- Evaluate expressions with data context
- Expand loops (generate rows)
- Place grids starting at current cursor position
- Record component placements
- Track cursor position after each element

**Result:** In-memory representation of all cells and objects with final positions.

---

### 3. Write Phase

**Input:** Expanded cell data + components  
**Output:** `.xlsx` file

**Process:**
- Create Excel workbook object
- Create sheets with specified names
- Write cell values at computed positions
- Apply formatting (if styled)
- Insert components at recorded positions
- Save workbook to file

**Output:** Binary Excel file ready to open in Excel/LibreOffice.

---

## Cursor Positioning

The **cursor** determines where the next content will be placed.

### Initial Position

Each sheet starts with cursor at `A1`.

### Grid Placement

Grids are placed **starting at the current cursor position**.

**Example:**

```xml
<Grid>
| Name | Age |
| Alice | 30 |
</Grid>
```

- If cursor is at `A1`, grid fills `A1:B2`
- After grid, cursor moves to `A3` (next row after grid)

---

### Cursor Movement Rules

| Element | Cursor Behavior |
|---------|-----------------|
| `<Grid>` | Moves to first column of next row after grid |
| `<For>` | Moves to next row after all loop iterations |
| `<Anchor>` | **Does not move cursor** (absolute positioning) |
| `<Merge>` | No movement (operates on existing cells) |
| `<Image>`, `<Chart>` | No movement (absolute positioning) |

---

### Example: Sequential Grids

**Template:**

```xml
<Sheet name="Report">
  <Grid>
  | Header 1 | Header 2 |
  </Grid>
  
  <Grid>
  | Data 1 | Data 2 |
  | Data 3 | Data 4 |
  </Grid>
</Sheet>
```

**Rendering:**

1. First grid placed at `A1:B1`, cursor moves to `A2`
2. Second grid placed at `A2:B3`, cursor moves to `A4`

**Result:**

| | A | B |
|-|---|---|
| 1 | Header 1 | Header 2 |
| 2 | Data 1 | Data 2 |
| 3 | Data 3 | Data 4 |

---

### Example: Anchor Positioning

**Template:**

```xml
<Sheet name="Report">
  <Grid>
  | Title |
  </Grid>
  
  <Anchor cell="D1">
    <Grid>
    | Side Note |
    </Grid>
  </Anchor>
  
  <Grid>
  | Next Row |
  </Grid>
</Sheet>
```

**Rendering:**

1. First grid at `A1`, cursor moves to `A2`
2. Anchor places grid at `D1` (cursor stays at `A2`)
3. Third grid at `A2`, cursor moves to `A3`

**Result:**

| | A | B | C | D |
|-|---|---|---|---|
| 1 | Title | | | Side Note |
| 2 | Next Row | | | |

---

## Loop Expansion

Loops generate multiple rows by iterating over arrays.

### Basic Loop

**Template:**

```xml
<For src="users">
  <Grid>
  | {{name}} | {{age}} |
  </Grid>
</For>
```

**Data:**

```json
{
  "users": [
    {"name": "Alice", "age": 30},
    {"name": "Bob", "age": 25}
  ]
}
```

**Rendering:**

1. Enter loop with `users` array (2 items)
2. **First iteration:** `name = "Alice"`, `age = 30`
   - Grid placed at `A1:B1`
   - Cursor moves to `A2`
3. **Second iteration:** `name = "Bob"`, `age = 25`
   - Grid placed at `A2:B2`
   - Cursor moves to `A3`

**Result:**

| | A | B |
|-|---|---|
| 1 | Alice | 30 |
| 2 | Bob | 25 |

---

### Loop Variables

Inside loops, special variables are available:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{_index}}` | Zero-based index | `0`, `1`, `2`, ... |
| `{{_number}}` | One-based number | `1`, `2`, `3`, ... |
| `{{_startRow}}` | First row of iteration | `1`, `2`, `3`, ... |
| `{{_endRow}}` | Last row of iteration | `1`, `2`, `3`, ... |

**Example:**

```xml
<For src="items">
  <Grid>
  | {{_number}} | {{name}} |
  </Grid>
</For>
```

**Data:**

```json
{
  "items": [
    {"name": "Apple"},
    {"name": "Banana"}
  ]
}
```

**Result:**

| | A | B |
|-|---|---|
| 1 | 1 | Apple |
| 2 | 2 | Banana |

---

### Multi-Row Loop Bodies

Loops can have multiple grids per iteration.

**Template:**

```xml
<For src="sections">
  <Grid>
  | **{{title}}** |
  </Grid>
  <Grid>
  | {{content}} |
  </Grid>
</For>
```

**Data:**

```json
{
  "sections": [
    {"title": "Intro", "content": "Welcome"},
    {"title": "Body", "content": "Main content"}
  ]
}
```

**Rendering:**

1. **First iteration:**
   - First grid at `A1` (title "Intro")
   - Second grid at `A2` (content "Welcome")
   - Cursor at `A3`
2. **Second iteration:**
   - First grid at `A3` (title "Body")
   - Second grid at `A4` (content "Main content")
   - Cursor at `A5`

**Result:**

| | A |
|-|---|
| 1 | **Intro** |
| 2 | Welcome |
| 3 | **Body** |
| 4 | Main content |

---

### Nested Loops

Loops can be nested to handle hierarchical data.

**Template:**

```xml
<For src="categories">
  <Grid>
  | **{{name}}** |
  </Grid>
  <For src="items">
    <Grid>
    | - {{title}} |
    </Grid>
  </For>
</For>
```

**Data:**

```json
{
  "categories": [
    {
      "name": "Fruits",
      "items": [
        {"title": "Apple"},
        {"title": "Banana"}
      ]
    },
    {
      "name": "Vegetables",
      "items": [
        {"title": "Carrot"}
      ]
    }
  ]
}
```

**Rendering:**

1. Outer loop iteration 1 (Fruits):
   - Grid at `A1`: "**Fruits**"
   - Inner loop iteration 1: Grid at `A2`: "- Apple"
   - Inner loop iteration 2: Grid at `A3`: "- Banana"
   - Cursor at `A4`
2. Outer loop iteration 2 (Vegetables):
   - Grid at `A4`: "**Vegetables**"
   - Inner loop iteration 1: Grid at `A5`: "- Carrot"
   - Cursor at `A6`

**Result:**

| | A |
|-|---|
| 1 | **Fruits** |
| 2 | - Apple |
| 3 | - Banana |
| 4 | **Vegetables** |
| 5 | - Carrot |

---

## Expression Evaluation

Expressions like `{{name}}` are evaluated during the render phase.

### Evaluation Process

1. **Parse expression:** Extract variable path (`name`, `user.name`, `items[0].price`)
2. **Lookup in context:** Traverse JSON data to find value
3. **Type coercion:** Convert to string for cell output
4. **Error handling:** If path not found, use empty string or error marker

---

### Lookup Order (Nested Loops)

When loops are nested, variables are looked up from **innermost to outermost scope**.

**Example:**

```xml
<For src="departments">
  <Grid>
  | Department: {{name}} |
  </Grid>
  <For src="employees">
    <Grid>
    | - {{name}} (Dept: {{name}}) |
    </Grid>
  </For>
</For>
```

**Data:**

```json
{
  "departments": [
    {
      "name": "Engineering",
      "employees": [
        {"name": "Alice"}
      ]
    }
  ]
}
```

**Rendering:**

In the inner loop:
- First `{{name}}` resolves to employee's `name` ("Alice")
- Outer loop's `name` is shadowed

**To access outer scope explicitly (future enhancement):**

```xml
| - {{name}} (Dept: {{..name}}) |
```

---

## Component Rendering

Components like `<Image>`, `<Chart>` are rendered **after** grids.

### Rendering Order

1. **Grids and loops:** Populate all cells
2. **Components:** Insert images, charts at specified positions
3. **Merges:** Apply cell merges after grid placement

---

### Component Positioning

Components use **absolute positioning** and do not affect cursor.

**Example:**

```xml
<Sheet name="Report">
  <Grid>
  | Sales Report |
  </Grid>
  
  <Image src="logo.png" cell="E1" width="100" height="50" />
  
  <Grid>
  | Q1 | Q2 |
  </Grid>
</Sheet>
```

**Rendering:**

1. First grid at `A1`, cursor moves to `A2`
2. Image inserted at `E1` (no cursor movement)
3. Second grid at `A2`, cursor moves to `A3`

**Result:**

| | A | B | ... | E |
|-|---|---|-----|---|
| 1 | Sales Report | | | [Logo] |
| 2 | Q1 | Q2 | | |

---

## Error Handling

### Parse Errors

**Causes:**
- Syntax errors (unclosed tags)
- Invalid attributes

**Behavior:**
- Rendering stops
- Error message with line number

---

### Runtime Errors

**Causes:**
- Undefined variable: `{{missing}}`
- Invalid data type: `{{user.name}}` when `user` is not an object

**Behavior:**
- **v1.0:** Insert empty string or error marker
- **Future:** Configurable (strict mode vs. lenient mode)

---

### Best Practices

1. **Validate data structure** before rendering
2. **Provide default values** in data prep
3. **Test templates** with sample data
4. **Handle missing data** gracefully

---

## Rendering Example: Invoice

**Template:**

```xml
<Book>
  <Sheet name="Invoice">
    <Grid>
    | Invoice #{{invoiceNumber}} |
    | Date: {{date}} |
    </Grid>
    
    <Grid>
    | Item | Quantity | Price | Total |
    </Grid>
    
    <For src="items">
      <Grid>
      | {{name}} | {{quantity}} | {{price}} | =B{{_startRow}}*C{{_startRow}} |
      </Grid>
    </For>
    
    <Grid>
    | | | **Total:** | =SUM(D4:D{{_endRow}}) |
    </Grid>
  </Sheet>
</Book>
```

**Data:**

```json
{
  "invoiceNumber": "INV-001",
  "date": "2024-01-15",
  "items": [
    {"name": "Widget", "quantity": 10, "price": 5.00},
    {"name": "Gadget", "quantity": 5, "price": 12.50}
  ]
}
```

**Rendering Steps:**

1. **Parse phase:** Build AST
2. **Render phase:**
   - First grid at `A1:A2` (header)
   - Second grid at `A3:D3` (table header)
   - Loop starts at row 4:
     - Iteration 1: Grid at `A4:D4` (Widget), `_startRow=4`
     - Iteration 2: Grid at `A5:D5` (Gadget), `_startRow=5`
   - Final grid at `A6:D6` (total)
3. **Write phase:** Generate `.xlsx`

**Result:**

| | A | B | C | D |
|-|---|---|---|---|
| 1 | Invoice #INV-001 | | | |
| 2 | Date: 2024-01-15 | | | |
| 3 | Item | Quantity | Price | Total |
| 4 | Widget | 10 | 5.00 | =B4*C4 |
| 5 | Gadget | 5 | 12.50 | =B5*C5 |
| 6 | | | **Total:** | =SUM(D4:D5) |

---

## Performance Considerations

### Large Data Sets

**Challenge:** Rendering 10,000+ rows can be slow

**Optimization:**
- Stream data instead of loading all into memory
- Use efficient Excel library (excelize)
- Limit formula recalculations

---

### Complex Templates

**Challenge:** Deeply nested loops with many expressions

**Optimization:**
- Pre-process data to flatten structures
- Cache evaluated expressions
- Avoid redundant lookups

---

## Future Enhancements

**v1.1:**
- Conditional rendering (`<If>`)
- Loop filtering (`<For src="items" filter="status=='active'">`)

**v1.2:**
- Partial rendering (update only changed cells)
- Template caching (compile once, render many times)

**v2.0:**
- Incremental rendering (stream large datasets)
- Parallel rendering (multi-sheet concurrency)

---

## Related Topics

- [Core Tags](./core-tags.md) - Tag syntax and behavior
- [Control Structures](./control-structures.md) - Loop details
- [Expressions](./expressions.md) - Expression evaluation
- [Examples](./examples.md) - Complete rendering examples
