# Control Structures

Control structures enable dynamic template behavior based on data. They allow loops, conditionals, and other programmatic patterns.

---

## For Loop

Iterates over array data, expanding rows downward for each element.

### Syntax

```xml
<For each="variableName in dataPath">
  <!-- content repeated for each element -->
</For>
```

### Attributes

#### `each` (required)
- **Type**: String
- **Format**: `<variable> in <path>`
- **Description**: Defines the loop variable and data source

**Components:**
- `<variable>`: Local variable name for the current iteration
- `<path>`: Dot-notation path to array in data context

### How It Works

1. Parser evaluates `<path>` to get an array from data context
2. For each element in array:
   - Creates local scope with `<variable>` bound to current element
   - Renders content within `<For>` tags
   - Advances cursor downward by number of rows generated
3. After loop completes, cursor is positioned after all generated content

### Basic Example

**Template:**
```xml
<Grid>
| Name | Email |
</Grid>

<For each="user in users">
<Grid>
| {{ user.name }} | {{ user.email }} |
</Grid>
</For>
```

**Data:**
```json
{
  "users": [
    {"name": "Alice", "email": "alice@example.com"},
    {"name": "Bob", "email": "bob@example.com"},
    {"name": "Charlie", "email": "charlie@example.com"}
  ]
}
```

**Output:**
```
A1: Name       | B1: Email
A2: Alice      | B2: alice@example.com
A3: Bob        | B3: bob@example.com
A4: Charlie    | B4: charlie@example.com
```

---

## Loop Variables

Built-in variables available within `<For>` loops.

### `loop.index`
- **Type**: Integer
- **Description**: Zero-based iteration index
- **Range**: 0 to (array.length - 1)

### `loop.number`
- **Type**: Integer
- **Description**: One-based iteration number
- **Range**: 1 to array.length

### `loop.startRow`
- **Type**: Integer
- **Description**: Starting row number for current iteration (absolute)
- **Use case**: Building cell references in formulas

### `loop.endRow`
- **Type**: Integer
- **Description**: Ending row number for current iteration (absolute)
- **Use case**: Creating dynamic ranges

### Example with Loop Variables

```xml
<Grid>
| # | Item | Quantity | Price | Total |
</Grid>

<For each="item in items">
<Grid>
| {{ loop.number }} | {{ item.name }} | {{ item.qty }} | {{ item.price }} | =C{{ loop.number + 1 }}*D{{ loop.number + 1 }} |
</Grid>
</For>

<Grid>
| | | | Total: | =SUM(E2:E{{ items.length + 1 }}) |
</Grid>
```

---

## Nested Loops

Loops can be nested to handle hierarchical data structures.

### Basic Nested Loop

```xml
<Grid>
| Category | Product | Price |
</Grid>

<For each="category in categories">
  <Grid>
  | {{ category.name }} | | |
  </Grid>
  
  <For each="product in category.products">
  <Grid>
  | | {{ product.name }} | {{ product.price }} |
  </Grid>
  </For>
</For>
```

**Data:**
```json
{
  "categories": [
    {
      "name": "Electronics",
      "products": [
        {"name": "Laptop", "price": 1200},
        {"name": "Mouse", "price": 25}
      ]
    },
    {
      "name": "Books",
      "products": [
        {"name": "Novel", "price": 15},
        {"name": "Textbook", "price": 80}
      ]
    }
  ]
}
```

### Nested Loop Variables

In nested loops, each level has its own `loop` variable:

```xml
<For each="category in categories">
  <Grid>
  | Category {{ loop.number }}: {{ category.name }} |
  </Grid>
  
  <For each="item in category.items">
  <Grid>
  | Item {{ loop.number }} in category {{ loop.parent.number }}: {{ item.name }} |
  </Grid>
  </For>
</For>
```

**Note:** `loop.parent` access is planned for future versions.

---

## Advanced For Loop Patterns

### With Formulas

```xml
<Grid>
| Product | Q1 | Q2 | Q3 | Q4 | Total |
</Grid>

<For each="product in products">
<Grid>
| {{ product.name }} | {{ product.q1 }} | {{ product.q2 }} | {{ product.q3 }} | {{ product.q4 }} | =SUM(B{{ loop.number + 1 }}:E{{ loop.number + 1 }}) |
</Grid>
</For>

<Grid>
| Total | =SUM(B2:B{{ products.length + 1 }}) | =SUM(C2:C{{ products.length + 1 }}) | =SUM(D2:D{{ products.length + 1 }}) | =SUM(E2:E{{ products.length + 1 }}) | =SUM(F2:F{{ products.length + 1 }}) |
</Grid>
```

### With Conditional Content

```xml
<For each="order in orders">
<Grid>
| Order #{{ order.id }} | {{ order.date }} | ${{ order.total }} | {{ order.status }} |
</Grid>

<!-- Future: Will be replaced with <If> when implemented -->
<!-- For now, use empty expressions or conditional data preparation -->
</For>
```

### With Multi-Row Content

```xml
<For each="invoice in invoices">
  <!-- Invoice header -->
  <Grid>
  | Invoice #{{ invoice.number }} | | Date: {{ invoice.date }} |
  </Grid>
  <Merge range="A{{ loop.startRow }}:B{{ loop.startRow }}" />
  
  <!-- Invoice items -->
  <For each="item in invoice.items">
  <Grid>
  | {{ item.name }} | {{ item.qty }} | ${{ item.price }} |
  </Grid>
  </For>
  
  <!-- Invoice total -->
  <Grid>
  | | Total: | ${{ invoice.total }} |
  </Grid>
  
  <!-- Spacer row -->
  <Grid>
  | | | |
  </Grid>
</For>
```

---

## If / Else (Conditional Rendering)

Conditionally render content based on boolean expressions.

**Status:** Planned for v1.1 (not yet implemented in goxcel v1.0)

### Syntax

```xml
<If cond="expression">
  <!-- rendered if expression is truthy -->
</If>
```

**With Else:**
```xml
<If cond="expression">
  <!-- rendered if truthy -->
<Else>
  <!-- rendered if falsy -->
</Else>
</If>
```

### Attributes

#### `cond` (required)
- **Type**: String (expression)
- **Description**: Expression evaluated to boolean

**Truthy values:**
- Non-zero numbers: `1`, `-5`, `3.14`
- Non-empty strings: `"hello"`, `"false"`
- Boolean true: `true`
- Non-empty arrays: `[1, 2]`
- Non-null objects: `{"key": "value"}`

**Falsy values:**
- Zero: `0`, `0.0`
- Empty string: `""`
- Boolean false: `false`
- Null: `null`
- Undefined: `undefined`
- Empty array: `[]`

### Examples

**Simple conditional:**
```xml
<If cond="showHeader">
<Grid>
| Company Name | Report Date |
</Grid>
</If>

<Grid>
| Data | Data |
</Grid>
```

**With Else:**
```xml
<If cond="isPremium">
<Grid>
| Premium Customer | Discount: 20% |
</Grid>
<Else>
<Grid>
| Standard Customer | Discount: 5% |
</Grid>
</Else>
</If>
```

**Comparison operators:**
```xml
<If cond="total > 1000">
<Grid>
| Discount Applied | 10% |
</Grid>
</If>

<If cond="status == 'paid'">
<Grid>
| Payment Status | PAID |
</Grid>
<Else>
<Grid>
| Payment Status | PENDING |
</Grid>
</Else>
</If>
```

**With nested paths:**
```xml
<If cond="user.subscription.isPremium">
<Grid>
| Welcome, Premium Member! |
</Grid>
</If>
```

### Combining with For Loops

```xml
<For each="item in items">
  <If cond="item.inStock">
  <Grid>
  | {{ item.name }} | In Stock | ${{ item.price }} |
  </Grid>
  <Else>
  <Grid>
  | {{ item.name }} | Out of Stock | - |
  </Grid>
  </Else>
  </If>
</For>
```

### Nested Conditionals

```xml
<If cond="hasData">
  <If cond="dataType == 'sales'">
  <Grid>
  | Sales Report |
  </Grid>
  <Else>
  <Grid>
  | Other Report |
  </Grid>
  </Else>
  </If>
<Else>
<Grid>
| No Data Available |
</Grid>
</Else>
</If>
```

---

## Switch / Case (Future)

**Status:** Under consideration for v2.0+

Multiple conditional branches based on a value:

```xml
<Switch value="status">
  <Case match="pending">
    <Grid>| Status: Pending |</Grid>
  </Case>
  
  <Case match="approved">
    <Grid>| Status: Approved |</Grid>
  </Case>
  
  <Case match="rejected">
    <Grid>| Status: Rejected |</Grid>
  </Case>
  
  <Default>
    <Grid>| Status: Unknown |</Grid>
  </Default>
</Switch>
```

---

## While Loop (Future)

**Status:** Under consideration (low priority)

Conditional looping:

```xml
<While cond="index < maxRows">
  <Grid>
  | Row {{ index }} |
  </Grid>
  <!-- Note: Need mechanism to update 'index' -->
</While>
```

**Challenges:**
- Requires mutable state
- Risk of infinite loops
- Complex to implement safely

**Alternative:** Pre-process data to create finite arrays, then use `<For>`

---

## Best Practices

### 1. Keep Loops Simple

**Good:**
```xml
<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.value }} |
</Grid>
</For>
```

**Avoid:**
```xml
<For each="item in items">
  <For each="sub in item.subs">
    <For each="detail in sub.details">
      <!-- Too deeply nested -->
    </For>
  </For>
</For>
```

### 2. Use Descriptive Variable Names

**Good:**
```xml
<For each="employee in employees">
<For each="product in products">
<For each="transaction in transactions">
```

**Avoid:**
```xml
<For each="i in items">
<For each="x in list">
<For each="e in data">
```

### 3. Pre-calculate Complex Logic

Instead of complex conditionals in template, prepare data:

**Better approach:**
```json
{
  "items": [
    {"name": "A", "displayPrice": "$10.00", "showDiscount": true},
    {"name": "B", "displayPrice": "$20.00", "showDiscount": false}
  ]
}
```

```xml
<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.displayPrice }} |
</Grid>
</For>
```

### 4. Document Complex Loops

```xml
<!-- 
  Generate invoice sections
  Each invoice contains:
  - Header row with invoice number
  - Item rows (nested loop)
  - Total row
  - Blank separator
-->
<For each="invoice in invoices">
  <!-- header -->
  <Grid>| Invoice #{{ invoice.number }} |</Grid>
  
  <!-- items -->
  <For each="item in invoice.items">
  <Grid>| {{ item.name }} | ${{ item.price }} |</Grid>
  </For>
  
  <!-- total -->
  <Grid>| Total | ${{ invoice.total }} |</Grid>
  <Grid>| | |</Grid>
</For>
```

### 5. Handle Empty Arrays

Prepare data to always have valid arrays:

```json
{
  "items": []  // Empty array instead of null/undefined
}
```

Or use conditionals (when available):

```xml
<If cond="items.length > 0">
  <For each="item in items">
    <Grid>| {{ item.name }} |</Grid>
  </For>
<Else>
  <Grid>| No items found |</Grid>
</Else>
</If>
```

---

## Error Handling

### Invalid Data Path

If `dataPath` doesn't exist or isn't an array:
- **Behavior**: Loop is skipped (zero iterations)
- **Warning**: Implementation should log warning

### Null/Undefined Values

If data path resolves to `null` or `undefined`:
- **Behavior**: Treated as empty array (zero iterations)
- **No error**: Graceful degradation

### Non-Array Values

If data path resolves to non-array value:
- **Behavior**: Implementation-dependent
  - goxcel v1.0: Treats as single-element array
  - Future: May throw error or skip

---

## Performance Considerations

### Large Datasets

For loops generate rows during rendering. Very large arrays can:
- Increase memory usage
- Slow rendering
- Create huge Excel files

**Recommendations:**
- Limit arrays to reasonable sizes (< 10,000 rows)
- Use pagination for large datasets
- Consider streaming mode (future feature)

### Nested Loops

Each nesting level multiplies row count:
- 100 categories × 50 products = 5,000 rows
- 10 invoices × 20 items = 200 rows

**Watch for:**
- Cartesian products (unintended)
- Deep nesting (>3 levels)

---

## Related Topics

- [Expressions](./expressions.md) - Variable interpolation syntax
- [Data Context](./data-context.md) - How data is structured and accessed
- [Rendering Semantics](./rendering.md) - How loops affect rendering

---

## Implementation Status

| Feature | Status | Version |
|---------|--------|---------|
| `<For>` loops | ✅ Implemented | v1.0 |
| Loop variables (`loop.index`, `loop.number`) | ✅ Implemented | v1.0 |
| `loop.startRow`, `loop.endRow` | ⏳ Planned | v1.1 |
| `<If>` / `<Else>` | ⏳ Planned | v1.1 |
| `<Switch>` / `<Case>` | 💭 Under consideration | v2.0+ |
| `<While>` | 💭 Under consideration | TBD |
