# CLI Reference

## Commands

### generate
Generate Excel file from template and data.

```bash
goxcel generate --template FILE --data FILE --output FILE [--dry-run]
```

**Options**:
- `--template, -t`: GXL template file path (required)
- `--data, -d`: JSON/YAML data file path (optional)
- `--output, -o`: Output XLSX file path (required)
- `--dry-run`: Show summary without writing file

**Example**:
```bash
goxcel generate -t report.gxl -d data.json -o report.xlsx
```

### version
Show goxcel version.

```bash
goxcel version
```

### help
Show help information.

```bash
goxcel help [command]
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Invalid arguments
- `3`: File not found
- `4`: Parse error
- `5`: Render error
