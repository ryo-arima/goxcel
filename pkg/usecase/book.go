package usecase

import (
	"context"
	"errors"

	"github.com/ryo-arima/goxcel/pkg/model"
)

// Renderer converts a GXL and input data into a model.Book workbook (backward compatibility)
type Renderer interface {
	Render(ctx context.Context, t *model.GXL, data any) (*model.Book, error)
}

// DefaultRenderer is a backward-compatible alias for DefaultBookUsecase
type DefaultRenderer struct {
	bookUsecase BookUsecase
}

// Render renders the template (backward compatibility)
func (r DefaultRenderer) Render(ctx context.Context, t *model.GXL, data any) (*model.Book, error) {
	if r.bookUsecase == nil {
		r.bookUsecase = NewDefaultBookUsecase()
	}
	return r.bookUsecase.Render(ctx, t, data)
}

// BookUsecase handles book-level rendering operations
type BookUsecase interface {
	Render(ctx context.Context, gxl *model.GXL, data any) (*model.Book, error)
}

// DefaultBookUsecase is the default implementation of BookUsecase
type DefaultBookUsecase struct {
	sheetUsecase SheetUsecase
}

// NewDefaultBookUsecase creates a new DefaultBookUsecase
func NewDefaultBookUsecase() *DefaultBookUsecase {
	return &DefaultBookUsecase{
		sheetUsecase: NewDefaultSheetUsecase(),
	}
}

// Render renders the GXL template into a Book
func (u *DefaultBookUsecase) Render(ctx context.Context, gxl *model.GXL, data any) (*model.Book, error) {
	if gxl == nil {
		return nil, errors.New("book usecase: gxl template is nil")
	}

	book := model.NewBook()

	// Normalize data to map[string]any for consistent access
	normalizedData := u.normalizeData(data)

	// Render each sheet
	for _, sheetTag := range gxl.Sheets {
		sheet, err := u.sheetUsecase.RenderSheet(ctx, &sheetTag, normalizedData)
		if err != nil {
			return nil, err
		}
		book.AddSheet(sheet)
	}

	return book, nil
}

// normalizeData converts any data type to map[string]any
func (u *DefaultBookUsecase) normalizeData(data any) map[string]any {
	if m, ok := data.(map[string]any); ok {
		return m
	}
	return map[string]any{"data": data}
}
