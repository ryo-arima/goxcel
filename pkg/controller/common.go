package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	gxlrepo "github.com/ryo-arima/goxcel/pkg/repository"
	"github.com/ryo-arima/goxcel/pkg/usecase"
	"github.com/ryo-arima/goxcel/pkg/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// InitGenerateCmd creates the 'generate' subcommand which parses a .gxl and generates to .xlsx.
func InitGenerateCmd() *cobra.Command {
	var (
		templatePath string
		templateName string
		dataPath     string
		outputPath   string
		dryRun       bool
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a .gxl template to .xlsx",
		Long:  "Generate a .gxl template with optional JSON or YAML data into an Excel .xlsx file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle template name (from .etc/templates)
			if templateName != "" {
				templatesDir := ".etc/templates"
				templatePath = filepath.Join(templatesDir, templateName, "base.gxl")
				
				// Check if template exists
				if _, err := os.Stat(templatePath); err != nil {
					return fmt.Errorf("template '%s' not found in %s", templateName, templatesDir)
				}
				
				// Auto-load data file if exists
				if dataPath == "" {
					dataFile := filepath.Join(templatesDir, templateName, "base.yaml")
					if _, err := os.Stat(dataFile); err == nil {
						dataPath = dataFile
					}
				}
			} else if templatePath == "" && len(args) > 0 {
				templatePath = args[0]
			}
			
			if strings.TrimSpace(templatePath) == "" {
				return fmt.Errorf("template path is required (pass as arg, --template, or --template-name)")
			}
			if err := RunGenerate(templatePath, dataPath, outputPath, dryRun); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&templatePath, "template", "t", "", "path to .gxl template file")
	cmd.Flags().StringVar(&templateName, "template-name", "", "template name from .etc/templates (e.g., 'b4-landscape')")
	cmd.Flags().StringVarP(&dataPath, "data", "d", "", "path to JSON or YAML data file (optional)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output .xlsx file path (optional; if empty with --dry-run prints summary)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "do not write .xlsx; print a summary instead")
	return cmd
}

// RunGenerate executes the generate command logic
func RunGenerate(templatePath, dataPath, outputPath string, dryRun bool) error {
	// Create config with file path
	conf := config.NewBaseConfigWithFile(templatePath)
	conf.Logger.DEBUG(util.CI1, "Starting generate command", map[string]interface{}{"template": templatePath, "data": dataPath, "output": outputPath, "dry_run": dryRun})

	// Validate template file existence early for clearer error
	if _, statErr := os.Stat(templatePath); statErr != nil {
		conf.Logger.ERROR(util.FSR2, "Template file not found")
		return fmt.Errorf("template not found: %w", statErr)
	}

	// Read and parse template via repository
	repo := gxlrepo.NewGxlRepository(conf)
	gt, err := repo.ReadGxl()
	if err != nil {
		conf.Logger.ERROR(util.RP2, "Failed to read GXL template")
		return fmt.Errorf("read gxl via repository: %w", err)
	}
	conf.Logger.DEBUG(util.GXLP1, "GXL template parsed successfully", map[string]interface{}{"sheets": len(gt.Sheets)})

	// Load data (optional)
	var data any
	if strings.TrimSpace(dataPath) != "" {
		conf.Logger.DEBUG(util.FSR1, "Reading data file", map[string]interface{}{"file": dataPath})
		db, err := os.ReadFile(dataPath)
		if err != nil {
			conf.Logger.ERROR(util.FSR2, "Failed to read data file")
			return fmt.Errorf("read data: %w", err)
		}
		var m map[string]any

		// Determine file format by extension
		ext := strings.ToLower(filepath.Ext(dataPath))
		switch ext {
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(db, &m); err != nil {
				conf.Logger.ERROR(util.FSR2, "Failed to parse data YAML")
				return fmt.Errorf("parse data yaml: %w", err)
			}
			conf.Logger.DEBUG(util.FSR1, "YAML data loaded successfully", nil)
		case ".json":
			if err := json.Unmarshal(db, &m); err != nil {
				conf.Logger.ERROR(util.FSR2, "Failed to parse data JSON")
				return fmt.Errorf("parse data json: %w", err)
			}
			conf.Logger.DEBUG(util.FSR1, "JSON data loaded successfully", nil)
		default:
			// Try JSON first, then YAML
			if err := json.Unmarshal(db, &m); err != nil {
				if err := yaml.Unmarshal(db, &m); err != nil {
					conf.Logger.ERROR(util.FSR2, "Failed to parse data as JSON or YAML")
					return fmt.Errorf("parse data (tried JSON and YAML): %w", err)
				}
				conf.Logger.DEBUG(util.FSR1, "Data loaded successfully as YAML", nil)
			} else {
				conf.Logger.DEBUG(util.FSR1, "Data loaded successfully as JSON", nil)
			}
		}
		data = m
	}

	// Generate
	conf.Logger.DEBUG(util.UR1, "Rendering template")
	bookUsecase := usecase.NewBookUsecase(conf)
	book, err := bookUsecase.Render(context.Background(), &gt, data)
	if err != nil {
		conf.Logger.ERROR(util.UR2, "Failed to render template")
		return fmt.Errorf("generate: %w", err)
	}
	conf.Logger.DEBUG(util.UR1, "Template rendered successfully", map[string]interface{}{"sheets": len(book.Sheets)})

	// Dry run summary or write
	if dryRun || strings.TrimSpace(outputPath) == "" {
		conf.Logger.INFO(util.CC1, "Dry run summary")
		PrintBookSummary(book)
		return nil
	}

	// Write XLSX file
	conf.Logger.DEBUG(util.RW1, "Writing XLSX file", map[string]interface{}{"output": outputPath})
	// Ensure output directory exists
	outDir := filepath.Dir(outputPath)
	if _, derr := os.Stat(outDir); os.IsNotExist(derr) {
		if mkErr := os.MkdirAll(outDir, 0o755); mkErr != nil {
			conf.Logger.ERROR(util.FSR2, "Failed to create output directory")
			return fmt.Errorf("create output directory: %w", mkErr)
		}
		conf.Logger.DEBUG(util.FSM1, "Created output directory", map[string]interface{}{"dir": outDir})
	}

	if err := gxlrepo.WriteBookToFile(book, outputPath); err != nil {
		conf.Logger.ERROR(util.RW2, "Failed to write XLSX file")
		return fmt.Errorf("write xlsx: %w", err)
	}

	conf.Logger.INFO(util.CC1, fmt.Sprintf("Successfully generated: %s", outputPath))
	return nil
}

// PrintBookSummary prints a summary of the book contents
func PrintBookSummary(b *model.Book) {
	logger := util.NewLogger(util.LoggerConfig{
		Component:    "goxcel",
		Service:      "summary",
		Level:        "INFO",
		Structured:   false,
		EnableCaller: false,
		Output:       "stdout",
	})
	logger.INFO(util.CC1, fmt.Sprintf("Workbook: %d sheets", len(b.Sheets)))
	for _, s := range b.Sheets {
		logger.INFO(util.CC1, fmt.Sprintf("Sheet %q: %d cells, %d merges, %d images, %d shapes, %d charts, %d pivots",
			s.Name, len(s.Cells), len(s.Merges), len(s.Images), len(s.Shapes), len(s.Charts), len(s.Pivots)))
	}
}
