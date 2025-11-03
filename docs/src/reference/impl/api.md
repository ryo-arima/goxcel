# API Documentation

## Programmatic Usage

### Basic Example

```go
import "github.com/ryo-arima/goxcel/pkg/controller"

func main() {
    // Generate Excel from template
    err := controller.Generate(
        "template.gxl",
        "data.json", 
        "output.xlsx",
    )
}
```

### With Custom Config

```go
import (
    "github.com/ryo-arima/goxcel/pkg/config"
    "github.com/ryo-arima/goxcel/pkg/controller"
)

func main() {
    cfg := config.NewBaseConfig()
    cfg.Logger.SetLevel("DEBUG")
    
    ctrl := controller.NewController(cfg)
    err := ctrl.Generate("template.gxl", "data.json", "output.xlsx")
}
```

## API Components

- [Usecase API](./api/usecase.md)
- [Repository API](./api/repository.md)
- [Config API](./api/config.md)
- [Controller API](./api/controller.md)
