# Quick Start

This guide walks you through creating your first Excel file with goxcel in 5 minutes.

## Step 1: Create a Template

Create a file named `hello.gxl`:

```xml
<Book name="HelloWorkbook">
  <Sheet name="Greeting">
    <Grid>
      | **Message** | **Value** |
      | Hello | {{ .name }} |
      | Generated | {{ .date }} |
    </Grid>
  </Sheet>
</Book>
```

**What this does**:
- Creates a workbook named "HelloWorkbook"
- Creates a sheet named "Greeting"
- Defines a 2-column table with headers in bold
- Uses `{{ .name }}` and `{{ .date }}` to inject data

## Step 2: Create Data File

Create `data.json`:

```json
{
  "name": "World",
  "date": "2025-11-04"
}
```

## Step 3: Generate Excel File

Run goxcel:

```bash
goxcel generate \
  --template hello.gxl \
  --data data.json \
  --output hello.xlsx
```

**Output**:
```
[INFO] Starting generate command
[INFO] GXL file parsed successfully
[INFO] Template rendered successfully
[INFO] XLSX file written successfully
```

## Step 4: Open the File

Open `hello.xlsx` in Excel or LibreOffice:

```
| Message   | Value      |
|-----------|------------|
| Hello     | World      |
| Generated | 2025-11-04 |
```

Headers will be **bold**.

## Example 2: Using Table Structure

Create `products.gxl`:

```xml
<Book>
  <Sheet name="Products">
    <Table>
      <Row>
        <Col>**Item**</Col>
        <Col>**Quantity**</Col>
        <Col>**Price**</Col>
      </Row>
      
      <Row each="item in items">
        <Col>{{ item.name }}</Col>
        <Col>{{ item.quantity:number }}</Col>
        <Col>{{ item.price:number }}</Col>
      </Row>
      
      <Row>
        <Col>**Total**</Col>
        <Col></Col>
        <Col>{{ total:number }}</Col>
      </Row>
    </Table>
  </Sheet>
</Book>
```

Create `products-data.json`:

```json
{
  "items": [
    {"name": "Widget A", "quantity": 5, "price": 10.00},
    {"name": "Widget B", "quantity": 3, "price": 25.50}
  ],
  "total": 102.50
}
```

Generate:

```bash
goxcel generate -t products.gxl -d products-data.json -o products.xlsx
```

**Output**: A clean table with header row, data rows (one per item), and total row.

**Key Differences from Grid**:
- No pipe `|` delimiters needed
- `<Row each="...">` iterates vertically (creates new rows downward)
- `<Col each="...">` iterates horizontally (creates new columns rightward)
- Cleaner syntax for structured tabular data

---

## Example 3: Using Template Import

Create `common/header.gxl`:

```xml
<Book>
  <Sheet name="Header">
    <Grid>
    | **{{ title }}** |
    | Generated: {{ date }} |
    </Grid>
  </Sheet>
</Book>
```

Create `report.gxl`:

```xml
<Book>
  <Import src="common/header.gxl" />
  
  <Sheet name="Data">
    <Table>
      <Row>
        <Col>**Item**</Col>
        <Col>**Value**</Col>
      </Row>
      <Row each="item in items">
        <Col>{{ item.name }}</Col>
        <Col>{{ item.value:number }}</Col>
      </Row>
    </Table>
  </Sheet>
</Book>
```

Create `report-data.json`:

```json
{
  "title": "Monthly Report",
  "date": "2026-02-08",
  "items": [
    {"name": "Revenue", "value": 50000},
    {"name": "Expenses", "value": 30000}
  ]
}
```

Generate:

```bash
goxcel generate -t report.gxl -d report-data.json -o report.xlsx
```

**Output**: Two sheets - "Header" (from imported template) and "Data" (from main template).

**Key Features**:
- Reuse common templates with `<Import>`
- Sheets appear in definition order
- Share data context across imported templates

---

## Example 4: Using For Loops

Create `invoice.gxl`:

```xml
<Book name="Invoice">
  <Sheet name="Items">
    <Grid>
      | **Item** | **Quantity** | **Price** |
    </Grid>
    
    <For each="item in items">
      <Grid>
        | {{ .item.name }} | {{ .item.quantity:number }} | {{ .item.price:number }} |
      </Grid>
    </For>
    
    <Grid>
      | **Total** | | {{ .total:number }} |
    </Grid>
  </Sheet>
</Book>
```

Create `invoice-data.json`:

```json
{
  "items": [
    {"name": "Widget A", "quantity": 5, "price": 10.00},
    {"name": "Widget B", "quantity": 3, "price": 25.50},
    {"name": "Widget C", "quantity": 2, "price": 15.75}
  ],
  "total": 157.00
}
```

Generate:

```bash
goxcel generate \
  --template invoice.gxl \
  --data invoice-data.json \
  --output invoice.xlsx
```

**Result**:
```
| Item     | Quantity | Price |
|----------|----------|-------|
| Widget A | 5        | 10.00 |
| Widget B | 3        | 25.50 |
| Widget C | 2        | 15.75 |
| Total    |          | 157.00|
```

## Example 3: Positioning with Anchors

Create `positioned.gxl`:

```xml
<Book name="Report">
  <Sheet name="Dashboard">
    <Anchor ref="A1" />
    <Grid>
      | **Title** |
      | Sales Report |
    </Grid>
    
    <Anchor ref="A5" />
    <Grid>
      | **Region** | **Sales** |
      | North | {{ .north:number }} |
      | South | {{ .south:number }} |
    </Grid>
    
    <Anchor ref="E5" />
    <Grid>
      | **Summary** |
      | Total: {{ .total:number }} |
    </Grid>
  </Sheet>
</Book>
```

Create `report-data.json`:

```json
{
  "north": 15000,
  "south": 23000,
  "total": 38000
}
```

This creates content at specific cell positions (A1, A5, E5).

## Common Features

### Cell Formatting

```xml
<Grid>
  | **Bold Text** | _Italic Text_ |
</Grid>
```

### Type Hints

```xml
<Grid>
  | {{ .text:string }} |
  | {{ .count:number }} |
  | {{ .active:boolean }} |
  | {{ .created:date }} |
</Grid>
```

### Formulas

```xml
<Grid>
  | Value 1 | Value 2 | Sum |
  | 10 | 20 | =A2+B2 |
</Grid>
```

### Cell Merging

```xml
<Grid>
  | Title |
</Grid>
<Merge range="A1:C1" />
```

## Dry Run Mode

Preview without creating a file:

```bash
goxcel generate \
  --template hello.gxl \
  --data data.json \
  --dry-run
```

Output shows parsed structure and cell data.

## Using YAML Data

goxcel also supports YAML:

```yaml
# data.yaml
name: World
date: 2025-11-04
items:
  - name: Item 1
    value: 100
  - name: Item 2
    value: 200
```

```bash
goxcel generate \
  --template template.gxl \
  --data data.yaml \
  --output output.xlsx
```

## Next Steps

- [Basic Concepts](./concepts.md) - Understand GXL fundamentals
- [Core Tags](../specification/core-tags.md) - Complete tag reference
- [Examples](../specification/examples.md) - More complex examples
- [Troubleshooting](../appendix/troubleshooting.md) - Common issues

## Optional: Format Your Template

Use the built-in formatter to keep your `.gxl` templates readable and consistent:

```bash
goxcel format template.gxl                # print to stdout
goxcel format -w template.gxl             # overwrite in place
goxcel format -o formatted.gxl template.gxl
```

What it does:
- Pretty-prints tags with indentation
- Inlines empty tags as a single line: `<Merge range="A1:C1"> </Merge>`
- Removes double blank lines outside content
- Aligns `|` columns inside `<Grid>` so tables are easy to read
- Preserves comments and significant text

## Quick Reference

### Template Structure
```xml
<Book name="WorkbookName">
  <Sheet name="SheetName">
    <!-- Content here -->
  </Sheet>
</Book>
```

### Grid Syntax
```xml
<Grid>
  | Header1 | Header2 |
  | Value1  | Value2  |
</Grid>
```

### Data Binding
```xml
{{ .path.to.value }}
{{ .value:type }}
```

### Control Structures
```xml
<For each="item in items">
  <!-- Repeated content -->
</For>
```

### Positioning
```xml
<Anchor ref="A1" />
<Merge range="A1:B2" />
```
