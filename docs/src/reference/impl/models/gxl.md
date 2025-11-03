# GXL Models

## Core Structures

### GXL
Root template structure.
```go
type GXL struct {
    HeaderTag HeaderTag
    BookTag   BookTag
    Sheets    []SheetTag
}
```

### SheetTag
Individual worksheet definition.
```go
type SheetTag struct {
    Name  string
    Nodes []any  // GridTag, ForTag, AnchorTag, etc.
}
```

### GridTag
Table with rows and cells.
```go
type GridTag struct {
    Ref     string         // Optional: absolute position
    Content string         // Raw pipe-delimited text
    Rows    []GridRowTag   // Parsed rows
}
```

### ForTag
Loop iteration.
```go
type ForTag struct {
    Each string    // "item in items"
    Body []any     // Child nodes
}
```

### AnchorTag
Position anchor.
```go
type AnchorTag struct {
    Ref string  // Cell reference like "A1"
}
```

## Node Hierarchy

```
GXL
└── SheetTag[]
    └── Nodes[]
        ├── GridTag
        ├── ForTag
        │   └── Body[]
        ├── AnchorTag
        ├── MergeTag
        └── ComponentTags (Image, Chart, etc.)
```
