# Repository API

## GxlRepository

GXL template file operations.

```go
type GxlRepository interface {
    ReadGxl() (*model.GXL, error)
}

// Usage
repo := gxlrepo.NewGxlRepository(config)
gxl, err := repo.ReadGxl()
```

## XlsxRepository

Excel file write operations.

```go
type XlsxRepository interface {
    WriteXlsx(book *model.Book, outputPath string) error
}

// Usage
repo := xlsxrepo.NewXlsxRepository(config)
err := repo.WriteXlsx(book, "output.xlsx")
```

## File Operations

Both repositories handle:
- File path validation
- Error wrapping with context
- Logging with message codes
- Resource cleanup
