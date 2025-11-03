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

## Example 2: Using Loops

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
