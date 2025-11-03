# Architecture Overview

## Layer Structure

```
Controller (CLI/API)
    ↓
Usecase (Business Logic)
    ↓
Repository (Data Access)
    ↓
Model (Data Structures)
```

## Core Components

- **Controller**: CLI commands, HTTP handlers
- **Usecase**: Template rendering, cell processing, sheet generation
- **Repository**: File I/O (GXL parsing, XLSX writing)
- **Model**: GXL AST, XLSX structures, Cell types

## Design Principles

- Clean Architecture: Dependencies point inward
- Single Responsibility: Each package has one purpose
- Dependency Injection: Logger and config injected
- Immutability: Models are read-only after creation

See [Packages](./packages.md) for detailed package structure.
