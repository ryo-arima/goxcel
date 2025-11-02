# Reference Overview

This reference section is divided into two main parts:

## 1. GXL Specification

The **GXL (Grid eXcel Language)** specification defines the template format used to describe Excel workbooks. This is a **format specification**, independent of any particular implementation.

**What you'll find here:**
- Template syntax and grammar
- Tag definitions and attributes
- Expression language semantics
- Data context and variable binding
- Validation rules and constraints
- Examples and best practices

**Audience:**
- Template authors who want to write .gxl files
- Tool developers implementing GXL parsers
- Anyone wanting to understand the GXL format

**Key Resources:**
- [GXL Format Specification](./gxl/specification.md) - Complete formal specification
- [Template Syntax](./gxl/syntax.md) - Practical syntax guide
- [Data Context](./gxl/data-context.md) - How data flows through templates

---

## 2. goxcel Implementation

The **goxcel implementation** documentation describes how the Go library processes GXL templates and generates Excel files. This is specific to the goxcel project.

**What you'll find here:**
- Architecture and design decisions
- Package structure and organization
- Internal data models and AST
- Rendering pipeline stages
- API documentation for programmatic use
- Logging and error handling
- CLI command reference

**Audience:**
- Developers using goxcel as a library
- Contributors working on goxcel codebase
- DevOps engineers deploying goxcel
- Anyone debugging or extending goxcel

**Key Resources:**
- [Architecture Overview](./impl/architecture.md) - High-level design
- [Package Structure](./impl/packages.md) - Code organization
- [Rendering Pipeline](./impl/pipeline.md) - How templates become Excel files
- [API Documentation](./impl/api.md) - Programmatic usage

---

## Navigation Tips

### For Template Authors
Start with:
1. [GXL Format Specification](./gxl/specification.md)
2. [Template Syntax](./gxl/syntax.md)
3. [Value Interpolation](./gxl/interpolation.md)
4. [Control Structures](./gxl/control-structures.md)

### For Library Users
Start with:
1. [Architecture Overview](./impl/architecture.md)
2. [API Documentation](./impl/api.md)
3. [CLI Reference](./impl/cli.md)
4. [Error Handling](./impl/errors.md)

### For Contributors
Start with:
1. [Architecture Overview](./impl/architecture.md)
2. [Package Structure](./impl/packages.md)
3. [Data Models](./impl/models.md)
4. [Rendering Pipeline](./impl/pipeline.md)

---

## Versioning

- **GXL Specification**: Version 0.1 (Draft)
- **goxcel Implementation**: Version 1.0.0

The GXL specification and goxcel implementation are versioned independently. A given version of goxcel supports one or more versions of the GXL specification.

## Compatibility Matrix

| goxcel Version | GXL Spec Version | Status |
|----------------|------------------|--------|
| 1.0.x          | 0.1              | Current |

Future versions may add support for additional GXL specification versions while maintaining backward compatibility.
