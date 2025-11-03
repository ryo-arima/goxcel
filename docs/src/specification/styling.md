# Styling

**Status:** Partially implemented (v1.0)

Styling system for cell formatting using markdown syntax and type hints.

---

## Implemented Features (v1.0)

### Markdown-Style Formatting

Inline text formatting using markdown syntax:

```xml
<Grid>
| **Bold Text** | _Italic Text_ | Normal Text |
</Grid>
```

**Supported**:
- `**text**`: Bold
- `_text_`: Italic

**Parsing**: Automatic detection and style application during rendering.

### Cell Type Hints

Explicit type specification for cells:

```xml
<Grid>
| {{ .quantity:int }} | {{ .price:float }} | {{ .active:bool }} |
</Grid>
```

**Supported Types**:
- `:int`, `:float`, `:number` → Number
- `:bool`, `:boolean` → Boolean
- `:date` → Date (ISO 8601)
- `:string` → String (explicit)

**Auto-inference**: Without type hints, goxcel automatically infers types from values.

---

## Planned Features (v1.1+)

### Style Tag

```xml
<Style selector="A1:C1" bold fillColor="#4CAF50" color="#FFFFFF" />
```

### Attributes (Future)

**Font**: `fontFamily`, `fontSize`, `bold`, `italic`, `underline`  
**Color**: `color` (text), `fillColor` (background)  
**Alignment**: `hAlign` (left/center/right), `vAlign` (top/middle/bottom)  
**Borders**: `border`, `borderColor`

---

## Implementation Status

| Feature | v1.0 | v1.1 | v1.2 |
|---------|------|------|------|
| Markdown Bold/Italic | ✅ | ✅ | ✅ |
| Type Hints | ✅ | ✅ | ✅ |
| Auto Type Inference | ✅ | ✅ | ✅ |
| Style Tag | ❌ | 🔄 | ✅ |
| Named Styles | ❌ | ❌ | 🔄 |
| Conditional Formatting | ❌ | ❌ | ❌ |

**Legend**: ✅ Implemented | 🔄 Planned | ❌ Not Planned

---

## Related

- [Core Tags](./core-tags.md)
- [Expressions](./expressions.md)

