# Package Structure

## Core Packages

### `pkg/controller`
CLI commands and HTTP handlers. Entry points for all operations.

### `pkg/usecase`
Business logic: template rendering, cell type inference, expression evaluation.

### `pkg/repository`
Data access: GXL file parsing, XLSX file writing, file I/O operations.

### `pkg/model`
Data structures: GXL AST nodes, XLSX models, cell types, styles.

### `pkg/config`
Configuration management: logger setup, file paths, options.

### `pkg/util`
Utilities: logger interface, message codes, helper functions.

## Dependencies

```
controller → usecase → repository → model
                ↓
              config
                ↓
              util
```

All packages depend on `model` and `util`. No circular dependencies.
