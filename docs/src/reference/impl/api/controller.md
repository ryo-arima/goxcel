# Controller API

## Generate Function

Main entry point for template generation.

```go
func Generate(templatePath, dataPath, outputPath string) error

// Usage
err := controller.Generate(
    "template.gxl",
    "data.json",
    "output.xlsx",
)
```

## With Custom Config

```go
cfg := config.NewBaseConfig()
cfg.Logger.SetLevel("DEBUG")

ctrl := controller.NewController(cfg)
err := ctrl.Generate("template.gxl", "data.json", "output.xlsx")
```

## Dry Run

```go
// Preview without writing file
err := controller.GenerateDryRun("template.gxl", "data.json")
```

## Error Handling

```go
if err != nil {
    // Check error type
    if errors.Is(err, os.ErrNotExist) {
        // File not found
    }
    // Log with context
    log.Printf("Generation failed: %v", err)
}
```
