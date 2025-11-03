# Basic Concepts

Understanding these core concepts will help you effectively use goxcel.

## Templates and Data

### Template (GXL File)

A **template** defines the structure of your Excel file. It's an XML file with `.gxl` extension that contains:
- Workbook and sheet definitions
- Table layouts (grids)
- Data binding expressions
- Control structures (loops)
- Positioning instructions

**Example**:
```xml
<Book name="Report">
  <Sheet name="Data">
    <Grid>
      | Name | {{ .userName }} |
    </Grid>
  </Sheet>
</Book>
```

### Data File

**Data** provides the values to inject into the template. Supported formats:
- JSON (`.json`)
- YAML (`.yaml`, `.yml`)

**Example**:
```json
{
  "userName": "John Doe"
}
```

### Generation Process

```
Template (.gxl) + Data (.json/.yaml) → goxcel → Excel File (.xlsx)
```

## GXL Structure

### Hierarchy

```
Book (Workbook)
└── Sheet (Worksheet)
    ├── Grid (Table)
    ├── For (Loop)
    ├── Anchor (Position)
    └── Merge (Cell merge)
```

### Book

The root element representing an Excel workbook:

```xml
<Book name="MyWorkbook">
  <!-- Sheets go here -->
</Book>
```

### Sheet

Represents a worksheet within the workbook:

```xml
<Sheet name="Sheet1">
  <!-- Content goes here -->
</Sheet>
```

You can have multiple sheets:

```xml
<Book name="Report">
  <Sheet name="Summary">
    <!-- ... -->
  </Sheet>
  <Sheet name="Details">
    <!-- ... -->
  </Sheet>
</Book>
```

### Grid

Defines a table using pipe-delimited syntax:

```xml
<Grid>
  | Header1 | Header2 | Header3 |
  | Value1  | Value2  | Value3  |
  | Value4  | Value5  | Value6  |
</Grid>
```

**Rules**:
- Each row starts and ends with `|`
- Cells are separated by `|`
- Whitespace around values is trimmed
- Consistent column count per row recommended

## Data Binding

### Expressions

Use `{{ }}` to inject data:

```xml
<Grid>
  | Name | {{ .name }} |
  | Age  | {{ .age }} |
</Grid>
```

With data:
```json
{
  "name": "Alice",
  "age": 30
}
```

### Nested Access

Use dot notation for nested data:

```xml
{{ .user.profile.email }}
```

With data:
```json
{
  "user": {
    "profile": {
      "email": "alice@example.com"
    }
  }
}
```

### Array Access in Loops

```xml
<For each="item in items">
  <Grid>
    | {{ .item.name }} | {{ .item.value }} |
  </Grid>
</For>
```

## Type System

### Automatic Type Inference

goxcel automatically detects types:

- **Number**: `123`, `45.67`, `-10`
- **Boolean**: `true`, `false`
- **Date**: `2025-11-04`, `2025-11-04T10:30:00`
- **Formula**: `=SUM(A1:A10)`
- **String**: Everything else

### Type Hints

Force a specific type:

```xml
{{ .value:number }}    <!-- Force as number -->
{{ .text:string }}     <!-- Force as string -->
{{ .flag:boolean }}    <!-- Force as boolean -->
{{ .created:date }}    <!-- Force as date -->
```

**Example**:
```xml
<Grid>
  | ID | {{ .id:string }} |     <!-- "001" stays as text -->
  | Amount | {{ .amount:number }} |  <!-- Ensures numeric -->
</Grid>
```

## Control Structures

### For Loops

Iterate over arrays:

```xml
<For each="item in items">
  <!-- This repeats for each item -->
  <Grid>
    | {{ .item.name }} |
  </Grid>
</For>
```

**Data**:
```json
{
  "items": [
    {"name": "Apple"},
    {"name": "Banana"},
    {"name": "Cherry"}
  ]
}
```

**Output**:
```
| Apple  |
| Banana |
| Cherry |
```

### Nested Loops

```xml
<For each="dept in departments">
  <Grid>
    | **{{ .dept.name }}** |
  </Grid>
  
  <For each="emp in dept.employees">
    <Grid>
      | {{ .emp.name }} | {{ .emp.role }} |
    </Grid>
  </For>
</For>
```

## Positioning

### Sequential (Default)

Content flows from current position:

```xml
<Grid>
  | Row 1 |
</Grid>
<Grid>
  | Row 2 |  <!-- Appears below Row 1 -->
</Grid>
```

### Absolute with Anchor

Set specific cell position:

```xml
<Anchor ref="A1" />
<Grid>
  | Title |
</Grid>

<Anchor ref="A10" />
<Grid>
  | Data starts here |
</Grid>
```

### Grid with Ref

Position a grid at specific location:

```xml
<Grid ref="E5">
  | Summary |
</Grid>
```

## Styling

### Markdown Syntax

Apply formatting within cells:

```xml
<Grid>
  | **Bold text** |
  | _Italic text_ |
  | **_Bold and italic_** |
</Grid>
```

**Supported**:
- `**text**` → Bold
- `_text_` → Italic

### Formulas

Excel formulas work directly:

```xml
<Grid>
  | 10 | 20 | =A1+B1 |
  | =SUM(A1:A10) |
</Grid>
```

### Cell Merging

Merge cells after defining them:

```xml
<Grid>
  | Large Title |
</Grid>
<Merge range="A1:C1" />
```

## Context Stack

When using loops, data context changes:

```xml
<!-- Root context -->
{{ .rootValue }}

<For each="item in items">
  <!-- Item context (can still access root) -->
  {{ .item.name }}
  {{ .rootValue }}  <!-- Still accessible -->
  
  <For each="sub in item.subs">
    <!-- Sub context (can access item and root) -->
    {{ .sub.value }}
    {{ .item.name }}
    {{ .rootValue }}
  </For>
</For>
```

**Context Stack** (innermost to outermost):
1. Current loop variable (`.sub`)
2. Parent loop variable (`.item`)
3. Root data (`.rootValue`)

## Best Practices

### Template Design

1. **Keep grids simple**: One table per Grid tag
2. **Use anchors sparingly**: Sequential flow is easier to maintain
3. **Name meaningfully**: Clear sheet and workbook names
4. **Comment complex sections**: Use XML comments `<!-- -->`

### Data Structure

1. **Match template paths**: Ensure JSON structure matches template expressions
2. **Use arrays for loops**: Structure data to match For loops
3. **Consistent types**: Use same type for similar values
4. **Avoid deep nesting**: Keep data structure reasonably flat

### Type Management

1. **Use type hints for IDs**: Force strings for numeric IDs like `001`
2. **Explicit numbers**: Use `:number` for calculations
3. **Date formats**: Use ISO 8601 format: `YYYY-MM-DD`
4. **Boolean clarity**: Use `true`/`false` not `"true"`/`"false"`

## Common Patterns

### Header with Data Rows

```xml
<Grid>
  | **Name** | **Email** | **Status** |
</Grid>
<For each="user in users">
  <Grid>
    | {{ .user.name }} | {{ .user.email }} | {{ .user.status }} |
  </Grid>
</For>
```

### Summary Section

```xml
<Anchor ref="A1" />
<Grid>
  | **Report Summary** |
</Grid>
<Merge range="A1:C1" />

<Grid>
  | Generated | {{ .date }} |
  | Total Records | {{ .count:number }} |
</Grid>
```

### Multi-Sheet Report

```xml
<Book name="MonthlyReport">
  <Sheet name="Summary">
    <Grid>
      | **Total Sales** | {{ .total:number }} |
    </Grid>
  </Sheet>
  
  <Sheet name="Details">
    <For each="item in items">
      <Grid>
        | {{ .item.name }} | {{ .item.amount:number }} |
      </Grid>
    </For>
  </Sheet>
</Book>
```

## Next Steps

- [Core Tags Reference](../specification/core-tags.md) - Complete tag documentation
- [Control Structures](../specification/control-structures.md) - Loops and conditionals
- [Expressions](../specification/expressions.md) - Data binding details
- [Examples](../specification/examples.md) - Real-world examples
