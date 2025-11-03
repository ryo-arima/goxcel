# Config API

## BaseConfig

Main configuration structure for goxcel.

```go
type BaseConfig struct {
    Logger      LoggerInterface
    FilePath    string
    OutputPath  string
    DataPath    string
    DryRun      bool
}

// Create default config
cfg := config.NewBaseConfig()

// Create with file path
cfg := config.NewBaseConfigWithFile("template.gxl")
```

## Logger Configuration

```go
// Set log level
cfg.Logger.SetLevel("DEBUG")  // DEBUG, INFO, WARN, ERROR

// Custom logger
cfg.Logger = myCustomLogger
```

## Usage

```go
import "github.com/ryo-arima/goxcel/pkg/config"

cfg := config.NewBaseConfig()
cfg.FilePath = "input.gxl"
cfg.OutputPath = "output.xlsx"
cfg.DataPath = "data.json"
```

## BaseConfig

Central configuration for goxcel operations.

```go
type BaseConfig struct {
    Logger   LoggerInterface
    FilePath string
}

// Create default config
cfg := config.NewBaseConfig()

// Create with file path
cfg := config.NewBaseConfigWithFile("template.gxl")

// Access logger
cfg.Logger.INFO(util.CI1, "Starting operation", nil)
```

## Logger Configuration

```go
// Set log level
cfg.Logger.SetLevel("DEBUG") // DEBUG, INFO, WARN, ERROR

// Custom logger implementation
type CustomLogger struct{}
func (l *CustomLogger) DEBUG(mcode util.MCode, msg string, ctx map[string]interface{}) {}
func (l *CustomLogger) INFO(mcode util.MCode, msg string, ctx map[string]interface{}) {}
// ... implement other methods

cfg.Logger = &CustomLogger{}
```

## Usage Pattern

```go
cfg := config.NewBaseConfig()
repo := gxlrepo.NewGxlRepository(cfg)
ctrl := controller.NewController(cfg)
```
