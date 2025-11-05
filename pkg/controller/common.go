package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		dataPath     string
		outputPath   string
		dryRun       bool
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a .gxl template to .xlsx",
		Long:  "Generate a .gxl template with optional JSON or YAML data into an Excel .xlsx file.",
		Run: func(cmd *cobra.Command, args []string) {
			if templatePath == "" && len(args) > 0 {
				templatePath = args[0]
			}
			if strings.TrimSpace(templatePath) == "" {
				log.Fatal("template path is required (pass as arg or --template)")
			}
			if err := runGenerate(templatePath, dataPath, outputPath, dryRun); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&templatePath, "template", "t", "", "path to .gxl template file")
	cmd.Flags().StringVarP(&dataPath, "data", "d", "", "path to JSON or YAML data file (optional)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output .xlsx file path (optional; if empty with --dry-run prints summary)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "do not write .xlsx; print a summary instead")
	return cmd
}

func runGenerate(templatePath, dataPath, outputPath string, dryRun bool) error {
	// Create config with file path
	conf := config.NewBaseConfigWithFile(templatePath)
	conf.Logger.DEBUG(util.CI1, "Starting generate command", map[string]interface{}{"template": templatePath, "data": dataPath, "output": outputPath, "dry_run": dryRun})

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
		conf.Logger.DEBUG(util.CC1, "Dry run mode - printing summary")
		printBookSummary(book)
		return nil
	}

	// Write XLSX file
	conf.Logger.DEBUG(util.RW1, "Writing XLSX file", map[string]interface{}{"output": outputPath})
	if err := gxlrepo.WriteBookToFile(book, outputPath); err != nil {
		conf.Logger.ERROR(util.RW2, "Failed to write XLSX file")
		return fmt.Errorf("write xlsx: %w", err)
	}

	conf.Logger.INFO(util.CC1, "Successfully generated XLSX file")
	fmt.Printf("Successfully generated: %s\n", outputPath)
	return nil
}

func printBookSummary(b *model.Book) {
	fmt.Printf("Workbook: %d sheets\n", len(b.Sheets))
	for _, s := range b.Sheets {
		fmt.Printf("- Sheet %q: %d cells, %d merges, %d images, %d shapes, %d charts, %d pivots\n",
			s.Name, len(s.Cells), len(s.Merges), len(s.Images), len(s.Shapes), len(s.Charts), len(s.Pivots))
	}
}
