package config_test

import (
	"testing"

	"github.com/ryo-arima/goxcel/pkg/config"
)

func TestNewBaseConfigWithFile(t *testing.T) {
	c := config.NewBaseConfigWithFile("/tmp/example.gxl")
	if c.FilePath != "/tmp/example.gxl" {
		t.Fatalf("FilePath = %q, want /tmp/example.gxl", c.FilePath)
	}
	if c.Logger == nil {
		t.Fatal("Logger is nil")
	}
}
