# Components

Components are special tags that insert rich content like images, shapes, charts, and pivot tables into worksheets.

**Current Status:** Component declarations are implemented in GXL v0.1 and goxcel v1.0 creates **placeholders** for these components. Full rendering is planned for future versions.

---

## Image

Inserts an image at a specific cell location.

### Syntax

```xml
<Image 
  ref="CellReference" 
  src="path/to/image.png" 
  width="120" 
  height="60" 
/>
```

### Attributes

#### `ref` (required)
- **Type**: String (cell reference)
- **Description**: Top-left anchor cell for the image
- **Format**: A1 notation (e.g., `B3`, `AA10`)

#### `src` (required)
- **Type**: String
- **Description**: Path to image file or resource identifier
- **Formats**: 
  - Relative path: `assets/logo.png`
  - Absolute path: `/path/to/image.png`
  - URL (future): `https://example.com/logo.png`
  - Resource key (future): `@logo`

#### `width` (optional)
- **Type**: Integer
- **Description**: Image width in pixels
- **Default**: Original image width

#### `height` (optional)
- **Type**: Integer
- **Description**: Image height in pixels  
- **Default**: Original image height

### Supported Formats

**v1.0 (Placeholder):**
- Any format (not validated)

**Planned (v1.2):**
- PNG (`.png`)
- JPEG (`.jpg`, `.jpeg`)
- GIF (`.gif`)
- BMP (`.bmp`)
- SVG (`.svg`) - via rasterization

### Examples

**Basic image:**
```xml
<Image ref="A1" src="company-logo.png" />
```

**With dimensions:**
```xml
<Image 
  ref="B3" 
  src="assets/product-photo.jpg" 
  width="200" 
  height="150" 
/>
```

**Multiple images:**
```xml
<Sheet name="Products">
  <Grid>
  | Product | Image | Description |
  </Grid>
  
  <For each="product in products">
    <Grid>
    | {{ product.name }} | | {{ product.description }} |
    </Grid>
    
    <Image 
      ref="B{{ loop.number + 1 }}" 
      src="{{ product.imagePath }}" 
      width="100" 
      height="100" 
    />
  </For>
</Sheet>
```

### Behavior

**v1.0:** Creates a placeholder text cell with image path
**v1.2:** Embeds actual image into workbook

### Best Practices

1. **Use relative paths** for portability
2. **Specify dimensions** to control layout
3. **Optimize images** before embedding (reduce file size)
4. **Test path resolution** with different working directories

---

## Shape

Inserts a shape (rectangle, arrow, etc.) with optional text.

### Syntax

```xml
<Shape 
  ref="CellReference" 
  kind="rectangle" 
  text="Label" 
  width="150" 
  height="50" 
/>
```

### Attributes

#### `ref` (required)
- **Type**: String (cell reference)
- **Description**: Top-left anchor cell

#### `kind` (required)
- **Type**: String
- **Description**: Shape type
- **Values**:
  - `rectangle` - Rectangle
  - `rounded` - Rounded rectangle
  - `ellipse` - Circle/ellipse
  - `arrow` - Arrow
  - `line` - Straight line
  - `star` - Star shape
  - `triangle` - Triangle
  - `diamond` - Diamond

#### `text` (optional)
- **Type**: String
- **Description**: Text content inside shape
- **Default**: Empty

#### `width` (optional)
- **Type**: Integer
- **Description**: Shape width in pixels
- **Default**: 100

#### `height` (optional)
- **Type**: Integer
- **Description**: Shape height in pixels
- **Default**: 50

#### `style` (optional)
- **Type**: String
- **Description**: Named style preset
- **Examples**: `banner`, `callout`, `warning`, `success`

### Examples

**Simple shape:**
```xml
<Shape ref="D3" kind="rectangle" text="Important" />
```

**Callout banner:**
```xml
<Shape 
  ref="A1" 
  kind="rounded" 
  text="URGENT: Read This" 
  width="200" 
  height="60" 
  style="warning" 
/>
```

**Workflow arrows:**
```xml
<Shape ref="B5" kind="rectangle" text="Step 1" width="120" height="40" />
<Shape ref="D5" kind="arrow" width="40" height="10" />
<Shape ref="F5" kind="rectangle" text="Step 2" width="120" height="40" />
<Shape ref="H5" kind="arrow" width="40" height="10" />
<Shape ref="J5" kind="rectangle" text="Step 3" width="120" height="40" />
```

### Behavior

**v1.0:** Creates placeholder text cell
**v1.2:** Renders actual shape with formatting

---

## Chart

Creates a chart visualization from data ranges.

### Syntax

```xml
<Chart 
  ref="CellReference" 
  type="column" 
  dataRange="A1:C10" 
  title="Chart Title" 
  width="500" 
  height="300" 
/>
```

### Attributes

#### `ref` (required)
- **Type**: String (cell reference)
- **Description**: Top-left anchor cell for chart

#### `type` (required)
- **Type**: String
- **Description**: Chart type
- **Values**:
  - `column` - Vertical bar chart
  - `bar` - Horizontal bar chart
  - `line` - Line chart
  - `pie` - Pie chart
  - `scatter` - Scatter plot
  - `area` - Area chart
  - `doughnut` - Doughnut chart
  - `radar` - Radar chart
  - `combo` - Combination chart

#### `dataRange` (required)
- **Type**: String
- **Description**: Source data range in A1 notation
- **Format**: `StartCell:EndCell`
- **Supports interpolation**: `A1:C{{ rowCount }}`

#### `title` (optional)
- **Type**: String
- **Description**: Chart title
- **Default**: No title

#### `width` (optional)
- **Type**: Integer
- **Description**: Chart width in pixels
- **Default**: 480

#### `height` (optional)
- **Type**: Integer
- **Description**: Chart height in pixels
- **Default**: 288

### Advanced Attributes (Planned v1.2+)

```xml
<Chart
  ref="A10"
  type="column"
  dataRange="A1:C10"
  title="Sales by Region"
  xAxisTitle="Region"
  yAxisTitle="Revenue ($)"
  legend="bottom"
  colors="#4CAF50,#2196F3,#FF9800"
  stacked="true"
/>
```

### Examples

**Basic column chart:**
```xml
<Grid>
| Month | Revenue | Target |
| Jan | 10000 | 12000 |
| Feb | 15000 | 12000 |
| Mar | 13000 | 12000 |
</Grid>

<Chart 
  ref="E1" 
  type="column" 
  dataRange="A1:C4" 
  title="Monthly Performance" 
/>
```

**Dynamic data range:**
```xml
<Grid>
| Category | Sales |
</Grid>

<For each="item in items">
<Grid>
| {{ item.category }} | {{ item.sales }} |
</Grid>
</For>

<Chart 
  ref="D1" 
  type="pie" 
  dataRange="A1:B{{ items.length + 1 }}" 
  title="Sales by Category" 
  width="400" 
  height="400" 
/>
```

**Multiple charts:**
```xml
<Sheet name="Dashboard">
  <!-- Data -->
  <Grid>
  | Product | Q1 | Q2 | Q3 | Q4 |
  | Product A | 100 | 120 | 110 | 130 |
  | Product B | 80 | 90 | 95 | 100 |
  </Grid>
  
  <!-- Chart 1: Column -->
  <Chart 
    ref="A5" 
    type="column" 
    dataRange="A1:E3" 
    title="Quarterly Sales" 
  />
  
  <!-- Chart 2: Line -->
  <Chart 
    ref="A20" 
    type="line" 
    dataRange="A1:E3" 
    title="Sales Trend" 
  />
</Sheet>
```

### Behavior

**v1.0:** Creates placeholder text cell with chart description
**v1.2:** Generates actual Excel chart object

---

## Pivot Table

Creates a pivot table from source data.

### Syntax

```xml
<Pivot 
  ref="CellReference" 
  sourceRange="A1:D100" 
  rows="Category" 
  columns="Month" 
  values="SUM:Sales" 
/>
```

### Attributes

#### `ref` (required)
- **Type**: String (cell reference)
- **Description**: Top-left cell for pivot table

#### `sourceRange` (required)
- **Type**: String
- **Description**: Source data range in A1 notation
- **Must include**: Header row with field names

#### `rows` (optional)
- **Type**: String (comma-separated)
- **Description**: Fields to use as row labels
- **Example**: `"Category,Product"`

#### `columns` (optional)
- **Type**: String (comma-separated)
- **Description**: Fields to use as column labels
- **Example**: `"Year,Month"`

#### `values` (required)
- **Type**: String (comma-separated)
- **Description**: Aggregate functions and fields
- **Format**: `FUNCTION:FieldName`
- **Functions**: `SUM`, `COUNT`, `AVERAGE`, `MAX`, `MIN`, `PRODUCT`, `STDDEV`, `VAR`
- **Examples**: `"SUM:Sales"`, `"COUNT:Orders,SUM:Revenue"`

#### `filters` (optional)
- **Type**: String (comma-separated)
- **Description**: Fields to use as filters
- **Example**: `"Region,Department"`

### Examples

**Basic pivot table:**
```xml
<Grid>
| Product | Category | Region | Sales |
| Widget A | Electronics | North | 1000 |
| Widget B | Electronics | South | 1500 |
| Gadget A | Toys | North | 800 |
| Gadget B | Toys | South | 1200 |
</Grid>

<Pivot 
  ref="F1" 
  sourceRange="A1:D5" 
  rows="Category" 
  columns="Region" 
  values="SUM:Sales" 
/>
```

**Multiple aggregations:**
```xml
<Pivot 
  ref="A20" 
  sourceRange="A1:E1000" 
  rows="Product,Category" 
  columns="Year" 
  values="SUM:Revenue,COUNT:Orders,AVERAGE:Price" 
  filters="Region,SalesRep" 
/>
```

**Dynamic source range:**
```xml
<For each="row in data">
<Grid>
| {{ row.product }} | {{ row.category }} | {{ row.sales }} |
</Grid>
</For>

<Pivot 
  ref="E1" 
  sourceRange="A1:C{{ data.length + 1 }}" 
  rows="Category" 
  values="SUM:Sales" 
/>
```

### Behavior

**v1.0:** Creates placeholder text cell
**v2.0:** Generates actual Excel pivot table

---

## Component Positioning

### Absolute Positioning

Components use absolute cell references:

```xml
<Image ref="B2" src="logo.png" />
<Chart ref="F2" type="column" dataRange="A1:C10" />
```

### Relative to Cursor

Components don't affect cursor position. They overlay cells without moving subsequent content:

```xml
<Grid>
| Header |
</Grid>

<!-- Image overlays B1:C3, but doesn't move cursor -->
<Image ref="B1" src="image.png" width="200" height="100" />

<!-- Grid continues at A2 -->
<Grid>
| Data |
</Grid>
```

### Avoiding Overlaps

Plan component positions to avoid overlapping:

```xml
<!-- Data occupies A1:D10 -->
<Grid>
| A | B | C | D |
</Grid>
<For each="row in rows">
<Grid>
| {{ row.a }} | {{ row.b }} | {{ row.c }} | {{ row.d }} |
</Grid>
</For>

<!-- Place chart starting at F1 (clear of data) -->
<Chart ref="F1" type="column" dataRange="A1:D10" />
```

---

## Future Enhancements

### Conditional Components (v1.2+)

```xml
<If cond="includeChart">
  <Chart ref="E1" type="column" dataRange="A1:C10" />
</If>
```

### Component Loops (v1.2+)

```xml
<For each="dataset in datasets">
  <Chart 
    ref="A{{ loop.index * 20 + 1 }}" 
    type="line" 
    dataRange="{{ dataset.range }}" 
    title="{{ dataset.title }}" 
  />
</For>
```

### Component Styling (v1.3+)

```xml
<Chart 
  ref="A1" 
  type="column" 
  dataRange="A1:C10"
  colors="#FF5733,#33FF57,#3357FF"
  borderColor="#000000"
  borderWidth="2"
  backgroundColor="#FFFFFF"
/>
```

### Interactive Components (v2.0+)

```xml
<Button 
  ref="A1" 
  text="Click Me" 
  action="macro:refreshData" 
/>

<Slider 
  ref="C1" 
  min="0" 
  max="100" 
  value="50" 
  linkedCell="D1" 
/>
```

---

## Best Practices

### 1. Plan Layout First

Sketch the desired layout before writing template:
```
+--------+--------+--------+
| Data   | Data   | Chart  |
|        |        |        |
+--------+--------+--------+
```

### 2. Use Descriptive Comments

```xml
<!-- Company logo in top-right -->
<Image ref="F1" src="logo.png" width="120" height="60" />

<!-- Sales chart below data table -->
<Chart ref="A15" type="column" dataRange="A1:C12" title="Monthly Sales" />
```

### 3. Test with Real Data

Verify component positions with actual data sizes:
- What if array has 100 items instead of 10?
- Does chart still fit on one page?

### 4. Optimize Resource Files

- **Images**: Use compressed formats (PNG, JPEG)
- **File size**: Keep images under 500KB when possible
- **Dimensions**: Resize images before embedding

### 5. Version Control Resources

Store images and other resources in version control alongside templates:
```
project/
├── templates/
│   └── report.gxl
└── assets/
    ├── logo.png
    ├── icon.png
    └── chart-background.png
```

---

## Implementation Status

| Component | v1.0 (Placeholder) | v1.2 (Rendering) | v2.0 (Advanced) |
|-----------|--------------------|--------------------|------------------|
| Image | ✅ | ⏳ | - |
| Shape | ✅ | ⏳ | - |
| Chart | ✅ | ⏳ | - |
| Pivot Table | ✅ | - | ⏳ |
| Button | - | - | 💭 |
| Slider | - | - | 💭 |

**Legend:**
- ✅ Implemented
- ⏳ Planned
- 💭 Under consideration
- \- Not planned

---

## Related Topics

- [Core Tags](./core-tags.md) - Positioning with Anchor
- [Control Structures](./control-structures.md) - Dynamic component placement
- [Expressions](./expressions.md) - Dynamic component attributes
