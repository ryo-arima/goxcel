# Complete Examples

This document provides **complete, working examples** of GXL templates with sample data and expected output.

---

## Example 1: Simple Invoice

### Template (`invoice.gxl`)

```xml
<Book>
  <Sheet name="Invoice">
    <Grid>
    | Invoice #{{invoiceNumber}} |
    | Date: {{date}} |
    | Customer: {{customer.name}} |
    </Grid>
    
    <Grid>
    | | | |
    </Grid>
    
    <Grid>
    | Item | Quantity | Price | Total |
    </Grid>
    
    <For src="items">
      <Grid>
      | {{name}} | {{quantity}} | ${{price}} | =B{{_startRow}}*C{{_startRow}} |
      </Grid>
    </For>
    
    <Grid>
    | | | **Grand Total:** | =SUM(D7:D{{_endRow}}) |
    </Grid>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "invoiceNumber": "INV-2024-001",
  "date": "2024-01-15",
  "customer": {
    "name": "Acme Corp"
  },
  "items": [
    {"name": "Consulting", "quantity": 10, "price": 150.00},
    {"name": "Development", "quantity": 40, "price": 200.00}
  ]
}
```

### Expected Output

**Sheet: Invoice**

| | A | B | C | D |
|-|---|---|---|---|
| 1 | Invoice #INV-2024-001 | | | |
| 2 | Date: 2024-01-15 | | | |
| 3 | Customer: Acme Corp | | | |
| 4 | | | | |
| 5 | Item | Quantity | Price | Total |
| 6 | Consulting | 10 | $150.00 | =B6*C6 → 1500 |
| 7 | Development | 40 | $200.00 | =B7*C7 → 8000 |
| 8 | | | **Grand Total:** | =SUM(D6:D7) → 9500 |

---

## Example 2: Sales Report with Multiple Sheets

### Template (`sales-report.gxl`)

```xml
<Book>
  <Sheet name="Summary">
    <Grid>
    | **Sales Report** |
    | Period: {{period}} |
    </Grid>
    
    <Grid>
    | | |
    </Grid>
    
    <Grid>
    | Region | Total Sales |
    </Grid>
    
    <For src="regions">
      <Grid>
      | {{name}} | ${{totalSales}} |
      </Grid>
    </For>
  </Sheet>
  
  <Sheet name="Details">
    <For src="regions">
      <Grid>
      | **Region: {{name}}** |
      </Grid>
      
      <Grid>
      | Product | Units | Revenue |
      </Grid>
      
      <For src="products">
        <Grid>
        | {{name}} | {{units}} | ${{revenue}} |
        </Grid>
      </For>
      
      <Grid>
      | | |
      </Grid>
    </For>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "period": "Q1 2024",
  "regions": [
    {
      "name": "North",
      "totalSales": 150000,
      "products": [
        {"name": "Widget A", "units": 1000, "revenue": 50000},
        {"name": "Widget B", "units": 2000, "revenue": 100000}
      ]
    },
    {
      "name": "South",
      "totalSales": 200000,
      "products": [
        {"name": "Widget A", "units": 1500, "revenue": 75000},
        {"name": "Widget C", "units": 2500, "revenue": 125000}
      ]
    }
  ]
}
```

### Expected Output

**Sheet: Summary**

| | A | B |
|-|---|---|
| 1 | **Sales Report** | |
| 2 | Period: Q1 2024 | |
| 3 | | |
| 4 | Region | Total Sales |
| 5 | North | $150000 |
| 6 | South | $200000 |

**Sheet: Details**

| | A | B | C |
|-|---|---|---|
| 1 | **Region: North** | | |
| 2 | Product | Units | Revenue |
| 3 | Widget A | 1000 | $50000 |
| 4 | Widget B | 2000 | $100000 |
| 5 | | | |
| 6 | **Region: South** | | |
| 7 | Product | Units | Revenue |
| 8 | Widget A | 1500 | $75000 |
| 9 | Widget C | 2500 | $125000 |
| 10 | | | |

---

## Example 3: Employee Directory with Anchored Logo

### Template (`directory.gxl`)

```xml
<Book>
  <Sheet name="Employees">
    <Anchor cell="E1">
      <Image src="company-logo.png" width="100" height="50" />
    </Anchor>
    
    <Grid>
    | **Employee Directory** |
    | As of {{date}} |
    </Grid>
    
    <Grid>
    | | |
    </Grid>
    
    <Grid>
    | ID | Name | Department | Email |
    </Grid>
    
    <For src="employees">
      <Grid>
      | {{id}} | {{firstName}} {{lastName}} | {{department}} | {{email}} |
      </Grid>
    </For>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "date": "2024-01-20",
  "employees": [
    {
      "id": "E001",
      "firstName": "Alice",
      "lastName": "Johnson",
      "department": "Engineering",
      "email": "alice@example.com"
    },
    {
      "id": "E002",
      "firstName": "Bob",
      "lastName": "Smith",
      "department": "Sales",
      "email": "bob@example.com"
    },
    {
      "id": "E003",
      "firstName": "Carol",
      "lastName": "Williams",
      "department": "Marketing",
      "email": "carol@example.com"
    }
  ]
}
```

### Expected Output

**Sheet: Employees**

| | A | B | C | D | E |
|-|---|---|---|---|---|
| 1 | **Employee Directory** | | | | [Logo] |
| 2 | As of 2024-01-20 | | | | |
| 3 | | | | | |
| 4 | ID | Name | Department | Email |
| 5 | E001 | Alice Johnson | Engineering | alice@example.com |
| 6 | E002 | Bob Smith | Sales | bob@example.com |
| 7 | E003 | Carol Williams | Marketing | carol@example.com |

---

## Example 4: Nested Categories

### Template (`catalog.gxl`)

```xml
<Book>
  <Sheet name="Catalog">
    <Grid>
    | **Product Catalog** |
    </Grid>
    
    <Grid>
    | | |
    </Grid>
    
    <For src="categories">
      <Grid>
      | **{{name}}** |
      </Grid>
      
      <Grid>
      | SKU | Product | Price |
      </Grid>
      
      <For src="products">
        <Grid>
        | {{sku}} | {{name}} | ${{price}} |
        </Grid>
      </For>
      
      <Grid>
      | | |
      </Grid>
    </For>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "categories": [
    {
      "name": "Electronics",
      "products": [
        {"sku": "E001", "name": "Laptop", "price": 999.99},
        {"sku": "E002", "name": "Mouse", "price": 29.99}
      ]
    },
    {
      "name": "Books",
      "products": [
        {"sku": "B001", "name": "Programming Book", "price": 49.99},
        {"sku": "B002", "name": "Novel", "price": 19.99}
      ]
    }
  ]
}
```

### Expected Output

**Sheet: Catalog**

| | A | B | C |
|-|---|---|---|
| 1 | **Product Catalog** | | |
| 2 | | | |
| 3 | **Electronics** | | |
| 4 | SKU | Product | Price |
| 5 | E001 | Laptop | $999.99 |
| 6 | E002 | Mouse | $29.99 |
| 7 | | | |
| 8 | **Books** | | |
| 9 | SKU | Product | Price |
| 10 | B001 | Programming Book | $49.99 |
| 11 | B002 | Novel | $19.99 |
| 12 | | | |

---

## Example 5: Loop Variables and Formulas

### Template (`inventory.gxl`)

```xml
<Book>
  <Sheet name="Inventory">
    <Grid>
    | # | Product | Quantity | Unit Price | Value |
    </Grid>
    
    <For src="items">
      <Grid>
      | {{_number}} | {{product}} | {{quantity}} | ${{unitPrice}} | =C{{_startRow}}*D{{_startRow}} |
      </Grid>
    </For>
    
    <Grid>
    | | | | **Total:** | =SUM(E2:E{{_endRow}}) |
    </Grid>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "items": [
    {"product": "Widget A", "quantity": 100, "unitPrice": 10.00},
    {"product": "Widget B", "quantity": 50, "unitPrice": 25.00},
    {"product": "Widget C", "quantity": 200, "unitPrice": 5.00}
  ]
}
```

### Expected Output

**Sheet: Inventory**

| | A | B | C | D | E |
|-|---|---|---|---|---|
| 1 | # | Product | Quantity | Unit Price | Value |
| 2 | 1 | Widget A | 100 | $10.00 | =C2*D2 → 1000 |
| 3 | 2 | Widget B | 50 | $25.00 | =C3*D3 → 1250 |
| 4 | 3 | Widget C | 200 | $5.00 | =C4*D4 → 1000 |
| 5 | | | | **Total:** | =SUM(E2:E4) → 3250 |

---

## Example 6: Cell Merging

### Template (`banner.gxl`)

```xml
<Book>
  <Sheet name="Report">
    <Grid>
    | Annual Report 2024 | | | |
    </Grid>
    <Merge range="A1:D1" />
    
    <Grid>
    | | | | |
    </Grid>
    
    <Grid>
    | Quarter | Revenue | Expenses | Profit |
    </Grid>
    
    <For src="quarters">
      <Grid>
      | {{name}} | ${{revenue}} | ${{expenses}} | =B{{_startRow}}-C{{_startRow}} |
      </Grid>
    </For>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "quarters": [
    {"name": "Q1", "revenue": 100000, "expenses": 70000},
    {"name": "Q2", "revenue": 120000, "expenses": 80000},
    {"name": "Q3", "revenue": 110000, "expenses": 75000},
    {"name": "Q4", "revenue": 130000, "expenses": 85000}
  ]
}
```

### Expected Output

**Sheet: Report**

| | A | B | C | D |
|-|---|---|---|---|
| 1 | Annual Report 2024 (merged across A1:D1) | | | |
| 2 | | | | |
| 3 | Quarter | Revenue | Expenses | Profit |
| 4 | Q1 | $100000 | $70000 | =B4-C4 → 30000 |
| 5 | Q2 | $120000 | $80000 | =B5-C5 → 40000 |
| 6 | Q3 | $110000 | $75000 | =B6-C6 → 35000 |
| 7 | Q4 | $130000 | $85000 | =B7-C7 → 45000 |

---

## Example 7: Multi-Sheet Workbook

### Template (`company-report.gxl`)

```xml
<Book>
  <Sheet name="Overview">
    <Grid>
    | **{{companyName}}** |
    | {{year}} Annual Report |
    </Grid>
    
    <Grid>
    | | |
    </Grid>
    
    <Grid>
    | Total Revenue: | ${{totalRevenue}} |
    | Total Expenses: | ${{totalExpenses}} |
    | Net Profit: | ${{netProfit}} |
    </Grid>
  </Sheet>
  
  <Sheet name="Revenue">
    <Grid>
    | Month | Amount |
    </Grid>
    
    <For src="revenueByMonth">
      <Grid>
      | {{month}} | ${{amount}} |
      </Grid>
    </For>
  </Sheet>
  
  <Sheet name="Expenses">
    <Grid>
    | Category | Amount |
    </Grid>
    
    <For src="expensesByCategory">
      <Grid>
      | {{category}} | ${{amount}} |
      </Grid>
    </For>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "companyName": "TechCorp Inc.",
  "year": 2024,
  "totalRevenue": 500000,
  "totalExpenses": 350000,
  "netProfit": 150000,
  "revenueByMonth": [
    {"month": "Jan", "amount": 40000},
    {"month": "Feb", "amount": 45000},
    {"month": "Mar", "amount": 50000}
  ],
  "expensesByCategory": [
    {"category": "Salaries", "amount": 200000},
    {"category": "Marketing", "amount": 80000},
    {"category": "Operations", "amount": 70000}
  ]
}
```

### Expected Output

**Sheet: Overview**

| | A | B |
|-|---|---|
| 1 | **TechCorp Inc.** | |
| 2 | 2024 Annual Report | |
| 3 | | |
| 4 | Total Revenue: | $500000 |
| 5 | Total Expenses: | $350000 |
| 6 | Net Profit: | $150000 |

**Sheet: Revenue**

| | A | B |
|-|---|---|
| 1 | Month | Amount |
| 2 | Jan | $40000 |
| 3 | Feb | $45000 |
| 4 | Mar | $50000 |

**Sheet: Expenses**

| | A | B |
|-|---|---|
| 1 | Category | Amount |
| 2 | Salaries | $200000 |
| 3 | Marketing | $80000 |
| 4 | Operations | $70000 |

---

## Best Practices Demonstrated

### 1. Clear Structure
- Use blank rows (`<Grid>| | |</Grid>`) for spacing
- Separate sections visually

### 2. Nested Data
- Use nested loops for hierarchical data (Example 4)
- Access nested properties with dot notation

### 3. Loop Variables
- Use `{{_number}}` for row numbering (Example 5)
- Use `{{_startRow}}` in formulas for dynamic references

### 4. Formulas
- Reference cells with row variables: `=B{{_startRow}}`
- Use SUM with `{{_endRow}}` for dynamic ranges

### 5. Anchoring
- Position logos/images independently (Example 3)
- Keep content flow unaffected

### 6. Multi-Sheet Reports
- Organize related data across sheets (Examples 2, 7)
- Summary sheet + detail sheets

---

## Common Patterns

### Pattern 1: Header + Table

```xml
<Grid>
| **{{title}}** |
</Grid>

<Grid>
| | |
</Grid>

<Grid>
| Column1 | Column2 |
</Grid>

<For src="data">
  <Grid>
  | {{field1}} | {{field2}} |
  </Grid>
</For>
```

---

### Pattern 2: Grouped Data

```xml
<For src="groups">
  <Grid>
  | **{{groupName}}** |
  </Grid>
  
  <For src="items">
    <Grid>
    | - {{itemName}} |
    </Grid>
  </For>
  
  <Grid>
  | | |
  </Grid>
</For>
```

---

### Pattern 3: Summary + Detail

```xml
<!-- Summary Sheet -->
<Sheet name="Summary">
  <Grid>
  | Category | Total |
  </Grid>
  <For src="categories">
    <Grid>
    | {{name}} | {{total}} |
    </Grid>
  </For>
</Sheet>

<!-- Detail Sheet -->
<Sheet name="Details">
  <For src="categories">
    <Grid>
    | **{{name}}** |
    </Grid>
    <For src="items">
      <Grid>
      | {{item}} | {{value}} |
      </Grid>
    </For>
  </For>
</Sheet>
```

---

## Example 9: Cell Types and Type Hints

### Template (`cell-types.gxl`)

```xml
<Book>
  <Sheet name="TypeDemo">
    <Grid>
    | Data Type | Auto-detected | With Type Hint | Description |
    </Grid>
    
    <Grid>
    | String | Hello World | {{ "Text" }} | Default text |
    | Number | {{ .price }} | {{ .quantity:int }} | Numeric values |
    | Float | {{ .discount }} | {{ .price:float }} | Decimal numbers |
    | Boolean | {{ .active }} | {{ .enabled:bool }} | TRUE/FALSE |
    | Formula | =SUM(B2:B5) | =AVERAGE(C2:C5) | Excel formulas |
    | Date | {{ .timestamp }} | {{ .created:date }} | Date values |
    </Grid>
    
    <Grid>
    | | | | |
    </Grid>
    
    <Grid>
    | **Mixed Content** | Price: {{ .price }} yen | Total: {{ .quantity }} items | Always string |
    </Grid>
    
    <Grid>
    | | | | |
    </Grid>
    
    <Grid>
    | **Force String Type** | {{ .zipCode:string }} | {{ .id:string }} | Preserve leading zeros |
    </Grid>
  </Sheet>
</Book>
```

### Data Context

```json
{
  "price": 1500.50,
  "quantity": 42,
  "discount": 0.15,
  "active": true,
  "enabled": false,
  "timestamp": "2025-11-03T15:30:00",
  "created": "2025-11-03",
  "zipCode": "00123",
  "id": "00456"
}
```

### Expected Output

**Sheet: TypeDemo**

| A | B | C | D |
|---|---|---|---|
| Data Type | Auto-detected | With Type Hint | Description |
| String | Hello World | Text | Default text |
| Number | 1500.5 (number) | 42 (number) | Numeric values |
| Float | 0.15 (number) | 1500.5 (number) | Decimal numbers |
| Boolean | TRUE (boolean) | FALSE (boolean) | TRUE/FALSE |
| Formula | (calculated) | (calculated) | Excel formulas |
| Date | 2025-11-03T15:30:00 | 2025-11-03 | Date values |
| | | | |
| **Mixed Content** | Price: 1500.5 yen | Total: 42 items | Always string |
| | | | |
| **Force String Type** | 00123 (text) | 00456 (text) | Preserve leading zeros |

**Type Inference:**
- Numbers are stored as Excel numeric cells (can be used in formulas)
- Formulas are evaluated by Excel
- Booleans become TRUE/FALSE
- Dates can be formatted with Excel date formats
- Type hints override automatic detection

**Available Type Hints:**
- `:int`, `:float`, `:number` - Numeric types
- `:bool`, `:boolean` - Boolean types
- `:date` - Date types
- `:string` - Force text (preserves leading zeros)

---

## Running Examples

### Using goxcel CLI

```bash
# Render a template
goxcel generate -t invoice.gxl -d data.json -o output.xlsx

# With YAML data
goxcel generate -t invoice.gxl -d data.yaml -o output.xlsx

# Dry run (preview structure)
goxcel generate -t invoice.gxl -d data.json --dry-run
```

### Using goxcel as Library (Go)

```go
package main

import (
    "github.com/ryo-arima/goxcel/pkg/config"
    "github.com/ryo-arima/goxcel/pkg/controller"
)

func main() {
    cfg := config.NewBaseConfig()
    ctrl := controller.NewCommonController(cfg)
    
    err := ctrl.Generate(
        "invoice.gxl",
        "data.json",
        "output.xlsx",
        false, // dry-run
    )
    if err != nil {
        panic(err)
    }
}
```

---

## Related Topics

- [Core Tags](./core-tags.md) - Tag syntax reference
- [Control Structures](./control-structures.md) - Loop details
- [Expressions](./expressions.md) - Data access and interpolation
- [Rendering](./rendering.md) - How templates are processed
