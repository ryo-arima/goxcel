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

	// Initialize import context for circular detection
	importCtx := &importContext{
		visitedFiles: make(map[string]bool),
		importDepth:  0,
		baseDir:      rcv.conf.BaseDir,
	}

	// Process imports at book level (creates new sheets)
	for _, importTag := range gxl.Imports {
		importedSheets, err := rcv.resolveAndRenderImports(ctx, importTag, normalizedData, importCtx)
		if err != nil {
			return nil, err
		}
		for _, importedSheet := range importedSheets {
			book.AddSheet(importedSheet)
		}
	}

	// Render each sheet defined in the main file
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

// resolveAndRenderImports loads an external .gxl file and renders the specified sheet
func (rcv *bookUsecase) resolveAndRenderImports(ctx context.Context, importTag model.ImportTag, data map[string]any, importCtx *importContext) ([]*model.Sheet, error) {
	// Check import depth limit
	if importCtx.importDepth >= maxImportDepth {
		return nil, errors.New("import depth limit exceeded (max 10)")
	}

	// Resolve file path
	filePath := importTag.Src
	if !isAbsolutePath(importTag.Src) && importCtx.baseDir != "" {
		filePath = joinPath(importCtx.baseDir, importTag.Src)
	}

	// Normalize path for circular detection
	normalizedPath, err := normalizePath(filePath)
	if err != nil {
		return nil, err
	}

	// Check for circular import
	if importCtx.visitedFiles[normalizedPath] {
		return nil, errors.New("circular import detected: " + normalizedPath)
	}

	// Mark as visited
	importCtx.visitedFiles[normalizedPath] = true
	importCtx.importDepth++
	defer func() {
		delete(importCtx.visitedFiles, normalizedPath)
		importCtx.importDepth--
	}()

	// Load and parse the imported .gxl file
	rcv.logger.DEBUG(util.UBR1, "Loading imported file for sheet creation", map[string]interface{}{
		"file": normalizedPath,
		"sheet": importTag.Sheet,
	})

	importedGxl, err := readGxlFile(normalizedPath, rcv.logger)
	if err != nil {
		return nil, err
	}

	// Find the specified sheet
	var targetSheetTag *model.SheetTag
	for i := range importedGxl.Sheets {
		if importedGxl.Sheets[i].Name == importTag.Sheet {
			targetSheetTag = &importedGxl.Sheets[i]
			break
		}
	}

	if targetSheetTag == nil {
		return nil, errors.New("sheet \"" + importTag.Sheet + "\" not found in " + normalizedPath)
	}

	// Render the imported sheet
	renderer := newSheetRenderer(rcv.conf)
	sheet, err := renderer.RenderSheet(ctx, targetSheetTag, data)
	if err != nil {
		return nil, err
	}

	return []*model.Sheet{sheet}, nil
}
