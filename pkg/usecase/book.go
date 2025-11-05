package usecase

import (
	"context"
	"errors"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	"github.com/ryo-arima/goxcel/pkg/util"
)

// BookUsecase handles book-level rendering operations
type BookUsecase interface {
	Render(ctx context.Context, gxl *model.GXL, data any) (*model.Book, error)
}

// bookUsecase is the default (unexported) implementation of BookUsecase
type bookUsecase struct {
	conf   config.BaseConfig
	logger util.Logger
}

// NewBookUsecase creates a new BookUsecase with config
func NewBookUsecase(conf config.BaseConfig) BookUsecase {
	return &bookUsecase{conf: conf, logger: conf.Logger}
}

// Render renders the GXL template into a Book
func (rcv *bookUsecase) Render(ctx context.Context, gxl *model.GXL, data any) (*model.Book, error) {
	if gxl == nil {
		return nil, errors.New("book usecase: gxl template is nil")
	}

	book := model.NewBook()

	// Normalize data to map[string]any for consistent access
	normalizedData := rcv.normalizeData(data)

	// Render each sheet using an internal renderer (no same-layer dependency)
	for _, sheetTag := range gxl.Sheets {
		renderer := newSheetRenderer(rcv.conf)
		sheet, err := renderer.RenderSheet(ctx, &sheetTag, normalizedData)
		if err != nil {
			return nil, err
		}
		book.AddSheet(sheet)
	}

	return book, nil
}

// normalizeData converts any data type to map[string]any
func (rcv *bookUsecase) normalizeData(data any) map[string]any {
	if m, ok := data.(map[string]any); ok {
		return m
	}
	return map[string]any{"data": data}
}
