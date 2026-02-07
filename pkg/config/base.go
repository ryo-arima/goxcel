package config

import "github.com/ryo-arima/goxcel/pkg/util"

// BaseConfig is a placeholder configuration root. Extend as needed.
type BaseConfig struct {
	FilePath string      // Path to the .gxl template file
	Logger   util.Logger // Logger instance
	BaseDir  string      // Base directory for resolving relative imports
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
	return BaseConfig{Logger: logger, BaseDir: "."}
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
	return BaseConfig{FilePath: filePath, Logger: logger, BaseDir: extractBaseDirFromPath(filePath)}
}

// extractBaseDirFromPath extracts the directory from a file path
func extractBaseDirFromPath(filePath string) string {
	if filePath == "" {
		return "."
	}
	// Extract base directory for relative import resolution
	lastSlash := -1
	for i, c := range filePath {
		if c == '/' || c == '\\' {
			lastSlash = i
		}
	}
	if lastSlash == -1 {
		return "."
	}
	return filePath[:lastSlash]
}
