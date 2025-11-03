# Usecase API

## CellUsecase

Cell-level operations: expression evaluation, type inference, style parsing.

```go
type CellUsecase interface {
    ExpandMustacheWithType(ctxStack []map[string]any, template string) (string, model.CellType)
    InferCellType(value string) model.CellType
    ResolvePath(ctxStack []map[string]any, path string) any
    ParseMarkdownStyle(text string) (string, *model.CellStyle)
    ParseTypeHint(expr string) (string, model.CellType)
}
```

## SheetUsecase

Sheet-level rendering operations.

```go
type SheetUsecase interface {
    RenderSheet(ctx context.Context, sheetTag *model.SheetTag, data map[string]any) (*model.Sheet, error)
}
```

## BookUsecase

Book-level operations (future).

```go
type BookUsecase interface {
    RenderBook(ctx context.Context, gxl *model.GXL, data any) (*model.Book, error)
}
```
