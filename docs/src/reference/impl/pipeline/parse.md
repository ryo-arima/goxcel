# Parse Phase

## Overview

Converts .gxl XML file into an in-memory AST (Abstract Syntax Tree).

## Process

1. **Read File**: Load .gxl template file
2. **XML Parse**: Use Go's `encoding/xml` decoder
3. **Tag Recognition**: Identify Sheet, Grid, For, Anchor, etc.
4. **Attribute Extraction**: Parse tag attributes
5. **Grid Parsing**: Convert pipe-delimited text to cell arrays
6. **AST Construction**: Build hierarchical node structure
7. **Validation**: Check syntax and constraints

## Input/Output

**Input**: `.gxl` file (XML text)  
**Output**: `model.GXL` struct (AST)

## Key Functions

- `repository.ReadGxl()`: Main entry point
- `parseSheetTag()`: Parse sheet elements
- `parseNodeTag()`: Parse child nodes
- `gridRowsFromText()`: Parse grid content

## Example

```xml
<Grid>| A | B |</Grid>
```
↓
```go
GridTag{
  Rows: []GridRowTag{
    {Cells: []string{"A", "B"}}
  }
}
```
