# Logging System

## Message Codes

All log messages include a message code (MCode) for easy lookup and troubleshooting.

### Format
```
[LEVEL] MCode: Message {context}
```

### Categories

- **CI-**: Controller/CLI operations
- **UC-**: Usecase/Cell operations
- **US-**: Usecase/Sheet operations
- **UB-**: Usecase/Book operations
- **RP-**: Repository operations
- **GXL-**: GXL parsing
- **FS-**: File system operations

### Log Levels

- **DEBUG**: Detailed trace information
- **INFO**: General informational messages
- **WARN**: Warning messages
- **ERROR**: Error messages

### Example

```
[INFO] CI1: Starting generate command {template: report.gxl, output: out.xlsx}
[DEBUG] UCE1: Expanding mustache template: {{ item.name }}
[ERROR] RP2: Failed to read GXL template {error: file not found}
```

See implementation in `pkg/util/logger.go`.
