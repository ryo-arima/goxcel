# Table Structure

**Status:** Implemented (v1.1)

The Table structure provides a declarative way to create tabular data with support for loops and dynamic content. Unlike Grid (which uses pipe-delimited syntax), Table uses explicit `<Row>` and `<Col>` tags.

---

## Syntax

```xml
<Table>
  <Row>
    <Col>Cell 1</Col>
    <Col>Cell 2</Col>
  </Row>
  <Row>
    <Col>Cell 3</Col>
    <Col>Cell 4</Col>
  </Row>
</Table>
```

---

## Components

### `<Table>`

Container for rows and columns.

**Attributes:** None

**Children:** One or more `<Row>` tags

### `<Row>`

Represents a table row. Supports optional loop iteration.

**Attributes:**
- `each` (optional): Loop syntax `"varName in dataPath"` for iterating over data

**Children:** One or more `<Col>` tags

### `<Col>`

Represents a table cell. Supports optional loop iteration.

**Attributes:**
- `each` (optional): Loop syntax `"varName in dataPath"` for iterating over data

**Content:** Text, expressions, or nested tags

---

## Basic Example

### Static Table

**Template:**
```xml
<Table>
  <Row>
    <Col>Product</Col>
    <Col>Price</Col>
    <Col>Stock</Col>
  </Row>
  <Row>
    <Col>Widget A</Col>
    <Col>100</Col>
    <Col>50</Col>
  </Row>
  <Row>
    <Col>Widget B</Col>
    <Col>200</Col>
    <Col>30</Col>
  </Row>
</Table>
```

**Output:**
```
| Product  | Price | Stock |
| Widget A | 100   | 50    |
| Widget B | 200   | 30    |
```

---

## Row Loops (Vertical Iteration)

Use `each` attribute on `<Row>` to iterate vertically over array data.

**Template:**
```xml
<Table>
  <Row>
    <Col>Product</Col>
    <Col>Price</Col>
    <Col>Stock</Col>
  </Row>
  <Row each="product in products">
    <Col>{{ product.name }}</Col>
    <Col>{{ product.price }}</Col>
    <Col>{{ product.stock }}</Col>
  </Row>
</Table>
```

**Data:**
```json
{
  "products": [
    {"name": "Widget A", "price": 10, "stock": 100},
    {"name": "Widget B", "price": 20, "stock": 200},
    {"name": "Widget C", "price": 30, "stock": 300}
  ]
}
```

**Output:**
```
| Product  | Price | Stock |
| Widget A | 10    | 100   |
| Widget B | 20    | 200   |
| Widget C | 30    | 300   |
```

---

## Column Loops (Horizontal Iteration)

Use `each` attribute on `<Col>` to iterate horizontally within a row.

**Template:**
```xml
<Table>
  <Row>
    <Col>Product</Col>
    <Col each="month in months">{{ month.name }}</Col>
  </Row>
  <Row>
    <Col>Sales</Col>
    <Col each="month in months">{{ month.value }}</Col>
  </Row>
</Table>
```

**Data:**
```json
{
  "months": [
    {"name": "Jan", "value": 1000},
    {"name": "Feb", "value": 1500},
    {"name": "Mar", "value": 2000}
  ]
}
```

**Output:**
```
| Product | Jan  | Feb  | Mar  |
| Sales   | 1000 | 1500 | 2000 |
```

---

## Nested Loops

Combine row and column loops for dynamic tables.

**Template:**
```xml
<Table>
  <Row>
    <Col>Product</Col>
    <Col each="quarter in quarters">Q{{ quarter }}</Col>
  </Row>
  <Row each="product in products">
    <Col>{{ product.name }}</Col>
    <Col each="quarter in quarters">{{ product.sales[quarter] }}</Col>
  </Row>
</Table>
```

**Data:**
```json
{
  "quarters": [1, 2, 3, 4],
  "products": [
    {
      "name": "Widget A",
      "sales": {
        "1": 100,
        "2": 120,
        "3": 110,
        "4": 130
      }
    },
    {
      "name": "Widget B",
      "sales": {
        "1": 200,
        "2": 220,
        "3": 210,
        "4": 230
      }
    }
  ]
}
```

**Output:**
```
| Product  | Q1  | Q2  | Q3  | Q4  |
| Widget A | 100 | 120 | 110 | 130 |
| Widget B | 200 | 220 | 210 | 230 |
```

---

## Loop Variables

Inside loops, special variables are available:

| Variable | Description | Example |
|----------|-------------|---------|
| `loop.index` | Zero-based index | `0`, `1`, `2` |
| `loop.number` | One-based number | `1`, `2`, `3` |

**Example:**
```xml
<Table>
  <Row>
    <Col>#</Col>
    <Col>Item</Col>
  </Row>
  <Row each="item in items">
    <Col>{{ loop.number }}</Col>
    <Col>{{ item.name }}</Col>
  </Row>
</Table>
```

---

## Comparison: Table vs Grid

| Feature | Table | Grid |
|---------|-------|------|
| Syntax | Explicit tags | Pipe-delimited |
| Verbosity | More verbose | Concise |
| Loop support | Row & Col loops | For tag only |
| Readability | Structured | Compact |
| Best for | Dynamic tables | Simple static tables |

**Grid Example:**
```xml
<Grid>
| Product | Price |
| ------- | ----- |
| Widget  | 100   |
</Grid>
```

**Table Equivalent:**
```xml
<Table>
  <Row>
    <Col>Product</Col>
    <Col>Price</Col>
  </Row>
  <Row>
    <Col>Widget</Col>
    <Col>100</Col>
  </Row>
</Table>
```

---

## Best Practices

### 1. Use Grid for Simple Tables

For static tables without loops, Grid syntax is more concise:

**Prefer:**
```xml
<Grid>
| Name | Value |
| ---- | ----- |
| A    | 1     |
| B    | 2     |
</Grid>
```

**Over:**
```xml
<Table>
  <Row><Col>Name</Col><Col>Value</Col></Row>
  <Row><Col>A</Col><Col>1</Col></Row>
  <Row><Col>B</Col><Col>2</Col></Row>
</Table>
```

### 2. Use Table for Dynamic Content

When you need row or column loops, Table is clearer:

**Good:**
```xml
<Table>
  <Row each="item in items">
    <Col>{{ item.name }}</Col>
    <Col>{{ item.value }}</Col>
  </Row>
</Table>
```

### 3. Header Rows First

Place header rows before loop rows:

**Good:**
```xml
<Table>
  <Row>
    <Col>Header 1</Col>
    <Col>Header 2</Col>
  </Row>
  <Row each="item in items">
    <Col>{{ item.col1 }}</Col>
    <Col>{{ item.col2 }}</Col>
  </Row>
</Table>
```

### 4. Consistent Column Count

Ensure all rows have the same number of columns:

**Bad:**
```xml
<Table>
  <Row>
    <Col>A</Col>
    <Col>B</Col>
  </Row>
  <Row>
    <Col>C</Col>  <!-- Missing column -->
  </Row>
</Table>
```

**Good:**
```xml
<Table>
  <Row>
    <Col>A</Col>
    <Col>B</Col>
  </Row>
  <Row>
    <Col>C</Col>
    <Col>D</Col>
  </Row>
</Table>
```

---

## Error Handling

| Error Condition | Behavior |
|----------------|----------|
| Empty Table | No cells generated |
| Row without Col | Row skipped |
| Invalid loop syntax | Error with message |
| Missing data path | Empty cells |

---

## Implementation Notes

### Architecture

- **Model Layer**: `TableTag`, `TableRowTag`, `TableColTag` structs
- **Parser**: Parses `<Table>`, `<Row>`, `<Col>` tags
- **Renderer**: `handleTable()`, `handleTableRow()`, `renderTableCols()`

### Loop Processing

1. Row loop: Iterates vertically, creates new rows
2. Col loop: Iterates horizontally within a row
3. Nested loops: Row loop outer, Col loop inner

---

## Future Enhancements

- Cell spanning: `<Col span="2">` for merged cells
- Column styling: `<Col class="highlight">`
- Table-level attributes: borders, padding
- Conditional rows: `<Row if="condition">`
- Header/body sections: `<THead>`, `<TBody>`

---

## Related Topics

- [Grid Tag](./core-tags.md#grid) - Pipe-delimited table syntax
- [For Loops](./control-structures.md#for-loop) - General loop syntax
- [Expressions](./expressions.md) - Data binding and interpolation
