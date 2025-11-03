# Expressions

Expressions enable dynamic content in GXL templates through value interpolation and evaluation.

---

## Value Interpolation

### Syntax

```
{{ expression }}
{{ expression:type }}
```

Double curly braces `{{ }}` evaluate expressions and insert the result into the document.

### Where Interpolation Works

1. **Inside Grid cells**
2. **In tag attributes**
3. **Within Excel formulas**

### Type Hints

GXL supports explicit type hints using colon syntax to control how values are written to Excel cells:

```
{{ .value:int }}      # Integer number
{{ .value:float }}    # Floating-point number
{{ .value:number }}   # Numeric value (auto-detect int/float)
{{ .value:bool }}     # Boolean (TRUE/FALSE)
{{ .value:date }}     # Date value
{{ .value:string }}   # String (force text)
```

**Without type hint, GXL automatically infers the cell type:**
- Values starting with `=` → Formula
- `true`/`false` → Boolean
- Numeric patterns → Number
- ISO date format → Date
- Everything else → String

---

## Basic Interpolation

### Simple Variables

```xml
<Grid>
| {{ title }} | {{ date }} | {{ amount }} |
</Grid>
```

**Data:**
```json
{
  "title": "Sales Report",
  "date": "2024-11-03",
  "amount": 1500.00
}
```

**Result:**
```
A1: Sales Report | B1: 2024-11-03 | C1: 1500.00
```

---

## Dot Notation (Object Access)

Access nested object properties using dot notation.

### Syntax

```
{{ object.property }}
{{ object.nested.property }}
```

### Examples

```xml
<Grid>
| {{ user.name }} | {{ user.email }} | {{ user.profile.age }} |
</Grid>
```

**Data:**
```json
{
  "user": {
    "name": "Alice",
    "email": "alice@example.com",
    "profile": {
      "age": 30,
      "city": "New York"
    }
  }
}
```

**Result:**
```
A1: Alice | B1: alice@example.com | C1: 30
```

---

## Array Access

### Index Notation

```
{{ array[0] }}
{{ array[1].property }}
```

### Examples

```xml
<Grid>
| {{ items[0].name }} | {{ items[0].price }} |
| {{ items[1].name }} | {{ items[1].price }} |
| {{ items[2].name }} | {{ items[2].price }} |
</Grid>
```

**Data:**
```json
{
  "items": [
    {"name": "Apple", "price": 1.50},
    {"name": "Banana", "price": 0.75},
    {"name": "Cherry", "price": 3.00}
  ]
}
```

---

## Cell Type Handling

GXL automatically detects and sets appropriate Excel cell types for proper data representation.

### Automatic Type Inference

GXL automatically infers cell types based on value patterns:

```xml
<Grid>
| Type | Example | Result |
| Number | {{ 42 }} | Excel numeric cell |
| Float | {{ 3.14159 }} | Excel numeric cell |
| Boolean | {{ true }} | Excel boolean cell (TRUE) |
| Formula | =SUM(A1:A10) | Excel formula cell |
| Date | {{ "2025-11-03" }} | Excel date cell |
| String | {{ "Hello" }} | Excel text cell |
</Grid>
```

**Inference Rules:**
- Values starting with `=` → Formula type
- `true` or `false` (case-insensitive) → Boolean type
- Numeric patterns (`123`, `45.67`, `-10.5`) → Number type
- ISO date format (`YYYY-MM-DD`) → Date type
- Everything else → String type

### Explicit Type Hints

Use type hints to explicitly control cell types:

```xml
<Grid>
| Description | Auto-detected | Type Hint |
| Integer | {{ .quantity }} | {{ .quantity:int }} |
| Float | {{ .price }} | {{ .price:float }} |
| Boolean | {{ .enabled }} | {{ .enabled:bool }} |
| Date | {{ .timestamp }} | {{ .timestamp:date }} |
| Force String | {{ .zipCode }} | {{ .zipCode:string }} |
</Grid>
```

**Data:**
```json
{
  "quantity": 42,
  "price": 1500.50,
  "enabled": false,
  "timestamp": "2025-11-03T15:30:00",
  "zipCode": "00123"
}
```

**Why use type hints?**
- Force numeric values to be treated as strings (e.g., zip codes, IDs)
- Ensure proper type when auto-detection might be ambiguous
- Control how data is stored in Excel for formulas and calculations

### Literal Values

You can also use literal values with type hints:

```xml
<Grid>
| String Literal | {{ "Hello World" }} |
| Number Literal | {{ 42 }} |
| Boolean Literal | {{ true }} |
| With Type Hint | {{ "123":string }} |
</Grid>
```

### Mixed Content

When multiple expressions appear in a single cell, the result is always a string:

```xml
<Grid>
| Description | Type |
| Price: {{ .price }} yen | String (mixed content) |
| Total: {{ .quantity }} items | String (mixed content) |
| {{ .amount }} | Number (single expression) |
</Grid>
```

---

## Array Length

```xml
<Grid>
| Total Items | {{ items.length }} |
</Grid>
```

---

## Attribute Interpolation

Expressions can be used in tag attributes.

### Chart with Dynamic Range

```xml
<Chart 
  ref="A10" 
  type="column" 
  dataRange="A1:C{{ rowCount }}" 
  title="Sales for {{ year }}"
/>
```

### Merge with Dynamic Range

```xml
<Merge range="A1:{{ lastColumn }}1" />
```

### Anchor with Computed Position

```xml
<Anchor ref="{{ startColumn }}{{ startRow }}" />
```

---

## Formula Interpolation

Use expressions within Excel formulas.

### Dynamic Cell References

```xml
<Grid>
| Total | =SUM(B2:B{{ rowCount + 1 }}) |
</Grid>
```

### Dynamic Ranges

```xml
<Grid>
| Average | =AVERAGE(A{{ startRow }}:A{{ endRow }}) |
| Maximum | =MAX(A{{ startRow }}:A{{ endRow }}) |
| Minimum | =MIN(A{{ startRow }}:A{{ endRow }}) |
</Grid>
```

### With Loop Variables

```xml
<For each="item in items">
<Grid>
| {{ item.name }} | {{ item.qty }} | {{ item.price }} | =B{{ loop.number + 1 }}*C{{ loop.number + 1 }} |
</Grid>
</For>
```

---

## Type Coercion

Expression results are automatically converted to appropriate types.

### Type Inference

| Expression Result | Excel Cell Type | Example |
|-------------------|----------------|---------|
| Number | Number | `123`, `45.67`, `-10` |
| String | Text | `"Hello"`, `"ABC123"` |
| Boolean | Boolean | `true`, `false` |
| Date (ISO 8601) | Date | `"2024-11-03"` |
| Null/Undefined | Empty | `null` |

### Examples

```xml
<Grid>
| {{ 100 }} | {{ "Text" }} | {{ true }} | {{ "2024-11-03" }} |
</Grid>
```

**Result:**
- A1: Number 100
- B1: Text "Text"  
- C1: Boolean TRUE
- D1: Date (formatted based on Excel settings)

---

## String Concatenation

### Using Plus Operator (Planned)

```xml
<Grid>
| {{ firstName + " " + lastName }} |
</Grid>
```

### Template Literals (Future)

```xml
<Grid>
| {{ `Full name: ${firstName} ${lastName}` }} |
</Grid>
```

### Current Workaround

Pre-concatenate in data:

```json
{
  "fullName": "Alice Smith",
  "firstName": "Alice",
  "lastName": "Smith"
}
```

```xml
<Grid>
| {{ fullName }} |
</Grid>
```

---

## Arithmetic Operators (Planned)

**Status:** Planned for v1.2

### Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `{{ a + b }}` |
| `-` | Subtraction | `{{ a - b }}` |
| `*` | Multiplication | `{{ a * b }}` |
| `/` | Division | `{{ a / b }}` |
| `%` | Modulo | `{{ a % b }}` |
| `^` | Exponentiation | `{{ a ^ b }}` |

### Examples

```xml
<Grid>
| Subtotal | {{ subtotal }} |
| Tax (10%) | {{ subtotal * 0.1 }} |
| Total | {{ subtotal + (subtotal * 0.1) }} |
</Grid>
```

### Precedence

Standard mathematical precedence:
1. Parentheses `()`
2. Exponentiation `^`
3. Multiplication `*`, Division `/`, Modulo `%`
4. Addition `+`, Subtraction `-`

---

## Comparison Operators (Planned)

**Status:** Planned for v1.1 (for use with `<If>`)

### Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal | `{{ a == b }}` |
| `!=` | Not equal | `{{ a != b }}` |
| `<` | Less than | `{{ a < b }}` |
| `>` | Greater than | `{{ a > b }}` |
| `<=` | Less than or equal | `{{ a <= b }}` |
| `>=` | Greater than or equal | `{{ a >= b }}` |

### Examples with If (Future)

```xml
<If cond="price > 100">
  <Grid>| Premium Product |</Grid>
</If>

<If cond="status == 'active'">
  <Grid>| Active |</Grid>
<Else>
  <Grid>| Inactive |</Grid>
</Else>
</If>
```

---

## Logical Operators (Planned)

**Status:** Planned for v1.1

### Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `&&` | Logical AND | `{{ a && b }}` |
| `\|\|` | Logical OR | `{{ a \|\| b }}` |
| `!` | Logical NOT | `{{ !a }}` |

### Examples

```xml
<If cond="isActive && isPremium">
  <Grid>| Premium Active User |</Grid>
</If>

<If cond="outOfStock || discontinued">
  <Grid>| Not Available |</Grid>
</If>

<If cond="!isDeleted">
  <Grid>| {{ item.name }} |</Grid>
</If>
```

---

## Built-in Functions (Planned)

**Status:** Planned for v1.2+

### String Functions

| Function | Description | Example |
|----------|-------------|---------|
| `len(str)` | String length | `{{ len(name) }}` |
| `upper(str)` | Uppercase | `{{ upper(name) }}` |
| `lower(str)` | Lowercase | `{{ lower(name) }}` |
| `trim(str)` | Trim whitespace | `{{ trim(input) }}` |
| `substr(str, start, len)` | Substring | `{{ substr(text, 0, 10) }}` |

### Array Functions

| Function | Description | Example |
|----------|-------------|---------|
| `len(array)` | Array length | `{{ len(items) }}` |
| `sum(array)` | Sum of numbers | `{{ sum(prices) }}` |
| `avg(array)` | Average | `{{ avg(scores) }}` |
| `min(array)` | Minimum value | `{{ min(values) }}` |
| `max(array)` | Maximum value | `{{ max(values) }}` |

### Math Functions

| Function | Description | Example |
|----------|-------------|---------|
| `round(num)` | Round to integer | `{{ round(3.7) }}` |
| `round(num, decimals)` | Round to decimals | `{{ round(3.14159, 2) }}` |
| `floor(num)` | Round down | `{{ floor(3.7) }}` |
| `ceil(num)` | Round up | `{{ ceil(3.2) }}` |
| `abs(num)` | Absolute value | `{{ abs(-5) }}` |

### Date Functions

| Function | Description | Example |
|----------|-------------|---------|
| `now()` | Current date/time | `{{ now() }}` |
| `today()` | Current date | `{{ today() }}` |
| `year(date)` | Extract year | `{{ year(orderDate) }}` |
| `month(date)` | Extract month | `{{ month(orderDate) }}` |
| `day(date)` | Extract day | `{{ day(orderDate) }}` |

### Examples

```xml
<Grid>
| Product Name | {{ upper(product.name) }} |
| Total Items | {{ len(items) }} |
| Average Price | {{ round(avg(prices), 2) }} |
| Report Date | {{ today() }} |
</Grid>
```

---

## Conditional Expressions (Ternary)

**Status:** Planned for v1.2

### Syntax

```
{{ condition ? valueIfTrue : valueIfFalse }}
```

### Examples

```xml
<Grid>
| Status | {{ isActive ? "Active" : "Inactive" }} |
| Price | {{ inStock ? price : "N/A" }} |
| Discount | {{ isPremium ? "20%" : "5%" }} |
</Grid>
```

---

## Null Coalescing

**Status:** Planned for v1.2

### Syntax

```
{{ value ?? defaultValue }}
```

### Examples

```xml
<Grid>
| Name | {{ user.name ?? "Anonymous" }} |
| Email | {{ user.email ?? "No email provided" }} |
| Phone | {{ user.phone ?? "N/A" }} |
</Grid>
```

---

## Escaping Special Characters

### Escaping Braces

To include literal `{{` or `}}` in output:

```xml
<Grid>
| Template syntax uses {{ "{{" }} and {{ "}}" }} |
</Grid>
```

**Result:**
```
A1: Template syntax uses {{ and }}
```

### Escaping Pipes

To include literal `|` in grid cells:

```xml
<Grid>
| Column A {{ "|" }} Column B |
</Grid>
```

Or use expression:

```xml
<Grid>
| {{ "Value | with | pipes" }} |
</Grid>
```

---

## Error Handling

### Undefined Variables

If variable doesn't exist in data context:
- **Behavior**: Renders empty string
- **Warning**: Implementation should log warning
- **No error**: Graceful degradation

### Invalid Paths

```xml
<Grid>
| {{ user.nonexistent.property }} |
</Grid>
```

If path is invalid:
- **Behavior**: Empty cell
- **Warning**: Logged if possible

### Type Errors

```xml
<Grid>
| {{ "string" + 123 }} |  <!-- Type mismatch -->
</Grid>
```

Behavior depends on implementation:
- May coerce to string: `"string123"`
- May return empty
- May throw error

---

## Best Practices

### 1. Keep Expressions Simple

**Good:**
```xml
<Grid>
| {{ user.name }} | {{ user.email }} |
</Grid>
```

**Avoid:**
```xml
<Grid>
| {{ user.profile.personal.names.first + " " + user.profile.personal.names.last }} |
</Grid>
```

**Better:**
Pre-compute in data:
```json
{
  "user": {
    "fullName": "Alice Smith",
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```

### 2. Use Descriptive Data Keys

**Good:**
```json
{
  "invoiceNumber": "INV-2024-001",
  "customerName": "Acme Corp",
  "orderTotal": 5000.00
}
```

**Avoid:**
```json
{
  "n": "INV-2024-001",
  "c": "Acme Corp",
  "t": 5000.00
}
```

### 3. Handle Missing Data

Prepare data to avoid undefined values:

```json
{
  "user": {
    "name": "Alice",
    "email": "alice@example.com",
    "phone": ""  // Empty string instead of null/undefined
  }
}
```

### 4. Pre-format Complex Values

**Instead of:**
```xml
<Grid>
| {{ price * 1.1 }} |  <!-- Calculate in template -->
</Grid>
```

**Do:**
```json
{
  "price": 100,
  "priceWithTax": 110
}
```

```xml
<Grid>
| {{ priceWithTax }} |
</Grid>
```

### 5. Document Expected Data Structure

```xml
<!--
  Required data structure:
  {
    "company": {
      "name": string,
      "address": string
    },
    "invoice": {
      "number": string,
      "date": string (ISO 8601),
      "items": [
        {
          "name": string,
          "qty": number,
          "price": number
        }
      ],
      "total": number
    }
  }
-->

<Book>
  <Sheet name="Invoice">
    <Grid>
    | {{ company.name }} |
    </Grid>
    ...
  </Sheet>
</Book>
```

---

## Expression Evaluation Order

1. **Parse template**: Extract expressions
2. **Evaluate expressions**: Resolve against data context
3. **Type conversion**: Convert to appropriate Excel types
4. **Cell generation**: Insert values into cells

---

## Performance Considerations

### Expression Complexity

Simple expressions are fast:
```xml
{{ user.name }}  <!-- Fast: direct property access -->
```

Complex expressions may be slower:
```xml
{{ sum(filter(items, item => item.active).map(item => item.price)) }}  <!-- Slower -->
```

**Recommendation:** Pre-compute complex values in data preparation step.

### Large Arrays

Accessing arrays in expressions:
```xml
{{ items[999].name }}  <!-- Fine for small arrays -->
```

For very large arrays (> 10,000 elements):
- Pre-filter data before passing to template
- Avoid iteration in expressions

---

## Related Topics

- [Data Context](./data-context.md) - How data is structured and accessed
- [Control Structures](./control-structures.md) - Using expressions in loops and conditionals
- [Core Tags](./core-tags.md) - Where expressions can be used

---

## Implementation Status

| Feature | Status | Version |
|---------|--------|---------|
| Basic interpolation `{{ var }}` | ✅ Implemented | v1.0 |
| Dot notation `{{ obj.prop }}` | ✅ Implemented | v1.0 |
| Array access `{{ arr[0] }}` | ✅ Implemented | v1.0 |
| Attribute interpolation | ✅ Implemented | v1.0 |
| Formula interpolation | ✅ Implemented | v1.0 |
| Arithmetic operators | ⏳ Planned | v1.2 |
| Comparison operators | ⏳ Planned | v1.1 |
| Logical operators | ⏳ Planned | v1.1 |
| Built-in functions | ⏳ Planned | v1.2 |
| Ternary operator | ⏳ Planned | v1.2 |
| Null coalescing | ⏳ Planned | v1.2 |
