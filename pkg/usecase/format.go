package usecase

import (
	"fmt"

	"github.com/ryo-arima/goxcel/pkg/config"
	gxlrepo "github.com/ryo-arima/goxcel/pkg/repository"
)

// FormatUsecase coordinates formatting of a .gxl template.
// It returns the pretty-printed bytes and does not write to disk.
type FormatUsecase interface {
	Format(templatePath string) ([]byte, error)
}

// DefaultFormatUsecase is the default implementation.
type DefaultFormatUsecase struct {
	conf config.BaseConfig
}

// NewFormatUsecase creates a new format use case with config.
func NewFormatUsecase(conf config.BaseConfig) FormatUsecase {
	return &DefaultFormatUsecase{conf: conf}
}

// Format pretty-prints the provided template file and returns the formatted bytes.
func (u *DefaultFormatUsecase) Format(templatePath string) ([]byte, error) {
	if templatePath == "" && u.conf.FilePath == "" {
		return nil, fmt.Errorf("template path is required")
	}

	// Ensure repository has the file path (prefer explicit arg)
	conf := u.conf
	if templatePath != "" {
		conf.FilePath = templatePath
	}

	repo := gxlrepo.NewGxlRepository(conf)
	b, err := repo.FormatGxl()
	if err != nil {
		return nil, err
	}
	return b, nil
}
