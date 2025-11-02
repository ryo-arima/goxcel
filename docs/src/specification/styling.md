# Styling

**Status:** Planned for v1.1+

The styling system allows applying visual formatting to cells and ranges.

---

## Overview

GXL will support styling through:
1. **Inline styles**: Attributes on `<Style>` tags
2. **Named styles**: Reusable style definitions
3. **Style classes**: CSS-like class system

---

## Style Tag

### Syntax

```xml
<Style selector="A1" name="header" />
```

### With Inline Properties

```xml
<Style 
  selector="A1:C1" 
  fontFamily="Arial" 
  fontSize="14" 
  bold 
  color="#333333" 
  fillColor="#FFF8E1" 
  hAlign="center" 
/>
```

---

## Attributes

### `selector` (required)
- Cell or range in A1 notation
- Examples: `A1`, `B2:D5`

### Font Properties
- `fontFamily`: Font name (`"Arial"`, `"Calibri"`, etc.)
- `fontSize`: Size in points (`10`, `12`, `14`)
- `bold`: Boolean flag
- `italic`: Boolean flag
- `underline`: Boolean flag

### Color Properties
- `color`: Text color (`#RRGGBB`)
- `fillColor`: Background color (`#RRGGBB`)

### Alignment Properties
- `hAlign`: `left`, `center`, `right`
- `vAlign`: `top`, `middle`, `bottom`

### Border Properties (Planned)
- `border`: Border style
- `borderColor`: Border color

---

## Examples

**Header row:**
```xml
<Grid>
| Product | Price | Stock |
</Grid>
<Style selector="A1:C1" bold fillColor="#4CAF50" color="#FFFFFF" hAlign="center" />
```

**Alternating rows (future):**
```xml
<Style selector="A2:C10" class="zebra-stripes" />
```

---

## Implementation Status

**v1.0:** Not implemented  
**v1.1:** Planned (basic styling)  
**v1.2:** Advanced styling

---

## Related

- [Core Tags](./core-tags.md)
- [Examples](./examples.md)
