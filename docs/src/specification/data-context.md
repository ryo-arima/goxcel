# Data Context

The data context is the JSON data structure passed to a GXL template during rendering. It provides the values for all expressions and control structures.

---

## Overview

### What is Data Context?

The data context is a JSON object that contains:
- **Variables**: Simple values (strings, numbers, booleans)
- **Objects**: Nested structures with properties
- **Arrays**: Lists of items for iteration

### How It Works

1. **Prepare data**: Create JSON data structure
2. **Pass to renderer**: Provide data when rendering template
3. **Access in template**: Use expressions `{{ }}` to access data
4. **Render output**: Template is filled with data values

---

## Data Structure

### Simple Values

```json
{
  "title": "Sales Report",
  "date": "2024-11-03",
  "amount": 1500.00,
  "isPaid": true
}
```

**Access in template:**
```xml
<Grid>
| {{ title }} | {{ date }} | {{ amount }} | {{ isPaid }} |
</Grid>
```

### Nested Objects

```json
{
  "company": {
    "name": "Acme Corp",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "zip": "10001"
    }
  }
}
```

**Access in template:**
```xml
<Grid>
| {{ company.name }} |
| {{ company.address.street }} |
| {{ company.address.city }}, {{ company.address.zip }} |
</Grid>
```

### Arrays

```json
{
  "items": [
    {"name": "Widget A", "price": 10.00},
    {"name": "Widget B", "price": 25.00},
    {"name": "Widget C", "price": 15.00}
  ]
}
```

**Access in template:**
```xml
<Grid>
| Product | Price |
</Grid>

<For each="item in items">
<Grid>
| {{ item.name }} | ${{ item.price }} |
</Grid>
</For>
```

---

## Data Types

### Supported Types

| JSON Type | Excel Type | Example |
|-----------|------------|---------|
| String | Text | `"Hello World"` |
| Number | Number | `123`, `45.67` |
| Boolean | Boolean | `true`, `false` |
| Null | Empty | `null` |
| Array | (Iterable) | `[1, 2, 3]` |
| Object | (Structure) | `{"key": "value"}` |

### Type Conversion

#### Strings
```json
{"text": "Hello", "code": "ABC123"}
```
→ Excel text cells

#### Numbers
```json
{"integer": 123, "decimal": 45.67, "negative": -10.5}
```
→ Excel number cells

#### Booleans
```json
{"isActive": true, "isDeleted": false}
```
→ Excel boolean cells (TRUE/FALSE)

#### Dates
```json
{"date": "2024-11-03", "datetime": "2024-11-03T14:30:00Z"}
```
→ Excel date cells (formatted based on locale)

#### Null
```json
{"emptyField": null}
```
→ Empty Excel cell

---

## Common Patterns

### Invoice Data

```json
{
  "invoice": {
    "number": "INV-2024-001",
    "date": "2024-11-03",
    "dueDate": "2024-12-03"
  },
  "customer": {
    "name": "Acme Corp",
    "email": "billing@acme.com",
    "address": "123 Main St, New York, NY 10001"
  },
  "items": [
    {"description": "Consulting Services", "hours": 40, "rate": 150.00},
    {"description": "Development Work", "hours": 80, "rate": 200.00}
  ],
  "subtotal": 22000.00,
  "tax": 2200.00,
  "total": 24200.00
}
```

### Report Data

```json
{
  "report": {
    "title": "Monthly Sales Report",
    "period": "November 2024",
    "generatedAt": "2024-11-03T10:00:00Z"
  },
  "summary": {
    "totalSales": 150000.00,
    "totalOrders": 450,
    "averageOrder": 333.33
  },
  "salesByRegion": [
    {"region": "North", "sales": 50000.00, "orders": 150},
    {"region": "South", "sales": 45000.00, "orders": 135},
    {"region": "East", "sales": 30000.00, "orders": 90},
    {"region": "West": 25000.00, "orders": 75}
  ]
}
```

### Hierarchical Data

```json
{
  "departments": [
    {
      "name": "Engineering",
      "employees": [
        {"name": "Alice", "role": "Engineer", "salary": 100000},
        {"name": "Bob", "role": "Senior Engineer", "salary": 120000}
      ]
    },
    {
      "name": "Sales",
      "employees": [
        {"name": "Charlie", "role": "Sales Rep", "salary": 80000},
        {"name": "Diana", "role": "Sales Manager", "salary": 110000}
      ]
    }
  ]
}
```

---

## Best Practices

### 1. Use Clear Key Names

**Good:**
```json
{
  "customerName": "Acme Corp",
  "invoiceNumber": "INV-2024-001",
  "orderTotal": 5000.00
}
```

**Avoid:**
```json
{
  "cn": "Acme Corp",
  "inv": "INV-2024-001",
  "tot": 5000.00
}
```

### 2. Pre-format Data

Instead of complex template logic, prepare data:

**Better:**
```json
{
  "price": 100.00,
  "priceFormatted": "$100.00",
  "priceWithTax": 110.00,
  "discount": "20%"
}
```

### 3. Handle Missing Values

Provide default values instead of null/undefined:

**Good:**
```json
{
  "user": {
    "name": "Alice",
    "email": "alice@example.com",
    "phone": ""  // Empty string instead of null
  }
}
```

### 4. Flatten When Possible

Deeply nested structures are hard to work with:

**Instead of:**
```json
{
  "data": {
    "customer": {
      "profile": {
        "personal": {
          "name": {
            "first": "Alice",
            "last": "Smith"
          }
        }
      }
    }
  }
}
```

**Use:**
```json
{
  "customerFirstName": "Alice",
  "customerLastName": "Smith",
  "customerFullName": "Alice Smith"
}
```

### 5. Include Metadata

```json
{
  "_meta": {
    "version": "1.0",
    "generated": "2024-11-03T10:00:00Z",
    "source": "api-v2"
  },
  "data": {
    ...
  }
}
```

---

## Data Preparation

### From Database

```go
// Example in Go
rows, _ := db.Query("SELECT name, price FROM products")
defer rows.Close()

products := []map[string]interface{}{}
for rows.Next() {
    var name string
    var price float64
    rows.Scan(&name, &price)
    products = append(products, map[string]interface{}{
        "name": name,
        "price": price,
    })
}

data := map[string]interface{}{
    "products": products,
    "total": len(products),
}
```

### From API

```javascript
// Example in JavaScript
const response = await fetch('/api/sales');
const apiData = await response.json();

const templateData = {
    report: {
        title: "Sales Report",
        date: new Date().toISOString().split('T')[0]
    },
    sales: apiData.results.map(item => ({
        product: item.product_name,
        quantity: item.qty,
        revenue: item.total_amount
    }))
};
```

### From Files

```python
# Example in Python
import json
import csv

# Load CSV
with open('data.csv') as f:
    reader = csv.DictReader(f)
    items = list(reader)

# Prepare context
data = {
    "title": "Data Export",
    "date": "2024-11-03",
    "items": items,
    "count": len(items)
}

# Save as JSON
with open('context.json', 'w') as f:
    json.dump(data, f, indent=2)
```

---

## Validation

### Schema Validation (Recommended)

Use JSON Schema to validate data before rendering:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["invoice", "customer", "items"],
  "properties": {
    "invoice": {
      "type": "object",
      "required": ["number", "date"],
      "properties": {
        "number": {"type": "string"},
        "date": {"type": "string", "format": "date"}
      }
    },
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["description", "amount"],
        "properties": {
          "description": {"type": "string"},
          "amount": {"type": "number"}
        }
      }
    }
  }
}
```

---

## Related Topics

- [Expressions](./expressions.md) - How to access data in templates
- [Control Structures](./control-structures.md) - Iterating over arrays
- [Examples](./examples.md) - Complete examples with data
