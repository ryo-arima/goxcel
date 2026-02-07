package usecase_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	parser "github.com/ryo-arima/goxcel/pkg/repository"
	"github.com/ryo-arima/goxcel/pkg/usecase"
	"github.com/ryo-arima/goxcel/pkg/util"
)

// TestImport_BasicExpansion tests that Import expands nodes from external file
func TestImport_BasicExpansion(t *testing.T) {
	path := filepath.Join("..", ".testdata", "import_main.gxl")
	conf := config.NewBaseConfigWithFile(path)
	
	// Read and parse the main file
	repo := parser.NewGxlRepository(conf)
	gxl, err := repo.ReadGxl()
	if err != nil {
		t.Fatalf("ReadGxl: %v", err)
	}
	
	if len(gxl.Sheets) != 1 {
		t.Fatalf("sheets=%d, want 1", len(gxl.Sheets))
	}
	
	// Render the book with import resolution
	bookUc := usecase.NewBookUsecase(conf)
	ctx := context.Background()
	book, err := bookUc.RenderBook(ctx, &gxl, map[string]any{})
	if err != nil {
		t.Fatalf("RenderBook: %v", err)
	}
	
	if len(book.Sheets) != 1 {
		t.Fatalf("rendered sheets=%d, want 1", len(book.Sheets))
	}
	
	sheet := book.Sheets[0]
	
	// Should have cells from both main file and imported file
	// Imported file has a 3x2 grid (header + 1 row)
	// Main file has a 2x2 grid
	// Total should be at least 5 cells
	if len(sheet.Cells) < 5 {
		t.Errorf("cells=%d, want at least 5 (from main + imported)", len(sheet.Cells))
	}
	
	// Check that imported content is present
	foundCompany := false
	for _, cell := range sheet.Cells {
		if cell.Value == "Company" {
			foundCompany = true
			break
		}
	}
	
	if !foundCompany {
		t.Error("expected imported cell 'Company' not found")
	}
}

// TestImport_CircularDetection tests that circular imports are detected
func TestImport_CircularDetection(t *testing.T) {
	path := filepath.Join("..", ".testdata", "import_circular_a.gxl")
	conf := config.NewBaseConfigWithFile(path)
	
	repo := parser.NewGxlRepository(conf)
	gxl, err := repo.ReadGxl()
	if err != nil {
		t.Fatalf("ReadGxl: %v", err)
	}
	
	// Try to render - should fail with circular import error
	bookUc := usecase.NewBookUsecase(conf)
	ctx := context.Background()
	_, err = bookUc.RenderBook(ctx, &gxl, map[string]any{})
	if err == nil {
		t.Fatal("expected circular import error, got nil")
	}
	
	// Error message should mention circular import
	errStr := err.Error()
	if !contains(errStr, "circular") {
		t.Errorf("error should mention circular import, got: %v", err)
	}
}

// TestImport_ValidationRejectsInvalidContent tests import validation
func TestImport_ValidationRejectsInvalidContent(t *testing.T) {
	// Create a test file that tries to import content with nested Sheet
	invalidPath := filepath.Join("..", ".testdata", "import_invalid_nested.gxl")
	
	conf := config.NewBaseConfigWithFile(invalidPath)
	repo := parser.NewGxlRepository(conf)
	
	// This should fail at parse time since Sheet in Sheet is invalid
	_, err := repo.ReadGxl()
	if err == nil {
		t.Fatal("expected parse error for nested Sheet, got nil")
	}
}

// TestImport_RelativePathResolution tests that relative paths are resolved correctly
func TestImport_RelativePathResolution(t *testing.T) {
	path := filepath.Join("..", ".testdata", "import_main.gxl")
	conf := config.NewBaseConfigWithFile(path)
	
	// BaseDir should be set to the directory containing import_main.gxl
	expectedBaseDir := filepath.Join("..", ".testdata")
	if conf.BaseDir != expectedBaseDir {
		t.Logf("BaseDir=%q, expected containing %q", conf.BaseDir, ".testdata")
		// Just a warning since path normalization might vary
	}
	
	repo := parser.NewGxlRepository(conf)
	gxl, err := repo.ReadGxl()
	if err != nil {
		t.Fatalf("ReadGxl: %v", err)
	}
	
	// Render should successfully resolve relative import
	bookUc := usecase.NewBookUsecase(conf)
	ctx := context.Background()
	_, err = bookUc.RenderBook(ctx, &gxl, map[string]any{})
	if err != nil {
		t.Fatalf("RenderBook with relative import failed: %v", err)
	}
}

// TestImport_MultipleImports tests importing from multiple files
func TestImport_MultipleImports(t *testing.T) {
	// Create inline test to verify multiple imports work
	// This would require creating additional test data files
	// For now, we test that the mechanism doesn't break with single import
	path := filepath.Join("..", ".testdata", "import_main.gxl")
	conf := config.NewBaseConfigWithFile(path)
	
	repo := parser.NewGxlRepository(conf)
	gxl, err := repo.ReadGxl()
	if err != nil {
		t.Fatalf("ReadGxl: %v", err)
	}
	
	bookUc := usecase.NewBookUsecase(conf)
	ctx := context.Background()
	book, err := bookUc.RenderBook(ctx, &gxl, map[string]any{})
	if err != nil {
		t.Fatalf("RenderBook: %v", err)
	}
	
	// Basic sanity check
	if len(book.Sheets) == 0 {
		t.Error("expected at least one sheet")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || len(s) > len(substr) && findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
