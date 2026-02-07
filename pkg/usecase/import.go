package usecase

import (
	"path/filepath"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	parser "github.com/ryo-arima/goxcel/pkg/repository"
	"github.com/ryo-arima/goxcel/pkg/util"
)

// isAbsolutePath checks if a path is absolute
func isAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// joinPath joins path components
func joinPath(base, rel string) string {
	return filepath.Join(base, rel)
}

// normalizePath returns a clean, absolute path
func normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// Clean the path to remove . and .. elements
	return filepath.Clean(absPath), nil
}

// extractBaseDir extracts the directory from a file path
func extractBaseDir(filePath string) string {
	if filePath == "" {
		return "."
	}
	dir := filepath.Dir(filePath)
	if dir == "" {
		return "."
	}
	return dir
}

// readGxlFile reads a GXL file using the repository layer
func readGxlFile(filePath string, logger util.Logger) (model.GXL, error) {
	conf := config.BaseConfig{
		FilePath: filePath,
		Logger:   logger,
		BaseDir:  extractBaseDir(filePath),
	}
	repo := parser.NewGxlRepository(conf)
	return repo.ReadGxl()
}
