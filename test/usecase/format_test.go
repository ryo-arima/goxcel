package usecase_test

import (
	"path/filepath"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/usecase"
)

func TestFormatUsecase_Format_Success(t *testing.T) {
	conf := config.NewBaseConfig()
	u := usecase.NewFormatUsecase(conf)
	path := filepath.Join("..", ".testdata", "minimal.gxl")
	b, err := u.Format(path)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if len(b) == 0 {
		t.Fatalf("expected non-empty formatted bytes")
	}
	if string(b[:min(5, len(b))]) != "<?xml" {
		t.Fatalf("expected xml header prefix, got: %q", string(b[:min(5, len(b))]))
	}
}

func TestFormatUsecase_Format_ErrOnEmpty(t *testing.T) {
	u := usecase.NewFormatUsecase(config.NewBaseConfig())
	if _, err := u.Format(""); err == nil {
		t.Fatalf("expected error on empty template path")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
