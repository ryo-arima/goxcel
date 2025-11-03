# Error Handling

## Error Categories

### Parse Errors
Invalid GXL syntax or structure.

**Examples**:
- Unclosed tags
- Invalid attributes
- Malformed grid syntax

### Render Errors
Issues during template rendering.

**Examples**:
- Undefined variable path
- Invalid loop syntax
- Type conversion failures

### Write Errors
Problems writing Excel files.

**Examples**:
- Invalid cell reference
- File write permission denied
- Insufficient disk space

## Error Format

```go
error: "operation: specific problem: underlying cause"
```

**Example**:
```
error: "render: invalid anchor ref \"XYZ\": invalid cell reference format"
```

## Best Practices

1. Check error messages for MCode references
2. Use `--dry-run` to validate templates
3. Enable DEBUG logging for detailed traces
4. Validate input data structure matches template expectations
