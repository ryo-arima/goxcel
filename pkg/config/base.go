package config

import "github.com/ryo-arima/goxcel/pkg/util"

// BaseConfig is a placeholder configuration root. Extend as needed.
type BaseConfig struct {
	FilePath string      // Path to the .gxl template file
	Logger   util.Logger // Logger instance
}

// NewBaseConfig returns a default config instance.
func NewBaseConfig() BaseConfig {
	logger := util.NewLogger(util.LoggerConfig{
		Component:    "goxcel",
		Service:      "default",
		Level:        "INFO",
		Structured:   false,
		EnableCaller: false,
		Output:       "stdout",
	})
	return BaseConfig{Logger: logger}
}

// NewBaseConfigWithFile returns a config instance with the specified file path.
func NewBaseConfigWithFile(filePath string) BaseConfig {
	logger := util.NewLogger(util.LoggerConfig{
		Component:    "goxcel",
		Service:      "cli",
		Level:        "INFO",
		Structured:   false,
		EnableCaller: false,
		Output:       "stdout",
	})
	return BaseConfig{FilePath: filePath, Logger: logger}
}
