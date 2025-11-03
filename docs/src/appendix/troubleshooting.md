# Troubleshooting

## Template Errors

### XML Parsing Failed

**Symptom**: `failed to parse GXL XML` error

**Solutions**:
1. Validate XML syntax
2. Check all tags are properly closed
3. Ensure no special characters in attribute values
4. Verify file encoding is UTF-8

**Example**:
```xml
<!-- Bad -->
<Sheet name="Sheet 1">

<!-- Good -->
<Sheet name="Sheet1">
</Sheet>
```

### Invalid Grid Syntax

**Symptom**: Grid not rendering correctly

**Solutions**:
1. Use `|` as cell delimiter
2. Each row must start and end with `|`
3. Consistent column count per row

**Example**:
```xml
<!-- Bad -->
<Grid>
Header1 | Header2
Value1 | Value2
</Grid>

<!-- Good -->
<Grid>
| Header1 | Header2 |
| Value1 | Value2 |
</Grid>
```

### Anchor Reference Error

**Symptom**: `invalid anchor ref` error

**Solutions**:
1. Use valid Excel cell references (A1, B2, etc.)
2. Column letters must be uppercase
3. Row numbers must be positive

**Example**:
```xml
<!-- Bad -->
<Anchor ref="a1" />

<!-- Good -->
<Anchor ref="A1" />
```

## Data Errors

### Data Not Displaying

**Symptom**: Template expressions show as empty or literal

**Solutions**:
1. Verify JSON/YAML syntax
2. Check data path matches template expression
3. Ensure data file is loaded correctly

**Example**:
```json
// data.json
{
  "user": {
    "name": "John"
  }
}
```

```xml
<!-- Template -->
<Grid>
| Name |
| {{ .user.name }} |
</Grid>
```

### Type Inference Issues

**Symptom**: Numbers displayed as text

**Solutions**:
1. Use type hints: `{{ .value:number }}`
2. Ensure numeric values in JSON are not quoted
3. Check date format is ISO 8601

**Example**:
```json
// Bad
{
  "price": "100"
}

// Good
{
  "price": 100
}
```

### Loop Not Iterating

**Symptom**: `<For>` loop produces no output

**Solutions**:
1. Verify data path points to an array
2. Check loop syntax: `each="item in items"`
3. Ensure array is not empty

**Example**:
```json
{
  "items": [
    {"name": "Item 1"},
    {"name": "Item 2"}
  ]
}
```

```xml
<For each="item in items">
  <Grid>
  | {{ .item.name }} |
  </Grid>
</For>
```

## Output Errors

### Excel File Won't Open

**Symptom**: Generated `.xlsx` file is corrupted

**Solutions**:
1. Check for write permissions
2. Ensure output path exists
3. Verify sufficient disk space
4. Close Excel if file is already open

### Missing Styles

**Symptom**: Bold/italic formatting not applied

**Solutions**:
1. Use markdown syntax correctly: `**bold**`, `_italic_`
2. Ensure no extra spaces around markers
3. Check GXL version supports styling

**Example**:
```xml
<!-- Bad -->
<Grid>
| ** bold ** |
</Grid>

<!-- Good -->
<Grid>
| **bold** |
</Grid>
```

### Merged Cells Not Working

**Symptom**: `<Merge>` tag has no effect

**Solutions**:
1. Use valid range: `A1:B2`
2. Ensure cells exist before merging
3. Check for overlapping merge ranges

## Performance Issues

### Slow Generation

**Symptom**: Template takes long time to generate

**Solutions**:
1. Reduce loop iterations if possible
2. Simplify complex expressions
3. Use batch operations where applicable
4. Check system resources (CPU, memory)

### Memory Errors

**Symptom**: Out of memory errors

**Solutions**:
1. Process large datasets in chunks
2. Reduce template complexity
3. Increase available memory
4. Optimize data structure

## CLI Issues

### Command Not Found

**Symptom**: `goxcel: command not found`

**Solutions**:
1. Ensure goxcel is installed: `go install github.com/ryo-arima/goxcel/cmd/goxcel@latest`
2. Check `$GOPATH/bin` is in `$PATH`
3. Verify installation: `which goxcel`

### Invalid Arguments

**Symptom**: CLI command fails

**Solutions**:
1. Check required flags: `--template`, `--data`, `--output`
2. Verify file paths are correct
3. Use absolute paths if relative paths fail
4. Check file permissions

**Example**:
```bash
# Full command with all flags
goxcel generate \
  --template /path/to/template.gxl \
  --data /path/to/data.json \
  --output /path/to/output.xlsx
```

## Logging and Debugging

### Enable Debug Logging

Set log level in code:
```go
logger := util.NewLogger(util.LoggerConfig{
    Level: "DEBUG",
    // ... other config
})
```

### Check Log Messages

Look for message codes:
- `GXL-P1/P2`: GXL parsing
- `U-R1/R2`: Rendering
- `R-W1/W2`: Writing

### Common Error Codes

- `RP2`: Failed to read GXL file
- `UR2`: Failed to render template
- `RW2`: Failed to write XLSX file
- `FSR2`: Failed to read data file

## Getting Help

If issues persist:

1. Check GitHub issues: https://github.com/ryo-arima/goxcel/issues
2. Review specification docs
3. Create a minimal reproduction case
4. Report bug with logs and sample files
