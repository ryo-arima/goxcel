package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// getTemplatesDir returns the path to the templates directory
func getTemplatesDir() (string, error) {
	// Get executable directory
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exe)
	
	// Try relative to executable first (for installed binary)
	templatesDir := filepath.Join(exeDir, ".etc", "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		return templatesDir, nil
	}
	
	// Try relative to current working directory (for development)
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	templatesDir = filepath.Join(wd, ".etc", "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		return templatesDir, nil
	}
	
	return "", fmt.Errorf("templates directory not found")
}

// InitGetCmd initializes the 'get' command
func InitGetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get resources (templates, etc.)",
		Long:  "Get various resources like templates for goxcel",
	}

	getCmd.AddCommand(initGetTemplatesCmd())
	return getCmd
}

// initGetTemplatesCmd initializes the 'get templates' subcommand
func initGetTemplatesCmd() *cobra.Command {
	var outputDir string

	var getTemplatesCmd = &cobra.Command{
		Use:   "templates [template-name]",
		Short: "Get template files",
		Long:  "Copy template files (base.gxl and base.yaml) to the current or specified directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetTemplates(args, outputDir)
		},
	}

	getTemplatesCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for template files")

	return getTemplatesCmd
}

// runGetTemplates executes the get templates command
func runGetTemplates(args []string, outputDir string) error {
	var templateName string
	if len(args) > 0 {
		templateName = args[0]
	} else {
		// List available templates
		return listTemplates()
	}

	// Copy template files to output directory
	return copyTemplate(templateName, outputDir)
}

// listTemplates lists all available templates
func listTemplates() error {
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to locate templates directory: %w", err)
	}

	fmt.Println("Available templates:")
	fmt.Println()

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}

	fmt.Println()
	fmt.Println("Usage: goxcel get templates <template-name> [-o output-dir]")
	return nil
}

// copyTemplate copies a specific template to the output directory
func copyTemplate(templateName, outputDir string) error {
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to locate templates directory: %w", err)
	}

	templatePath := filepath.Join(templatesDir, templateName)

	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy files from template
	files := []string{"base.gxl", "base.yaml"}
	for _, filename := range files {
		srcPath := filepath.Join(templatePath, filename)
		dstPath := filepath.Join(outputDir, filename)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			fmt.Printf("Warning: %s not found in template, skipping\n", filename)
			continue
		}

		// Read source file
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filename, err)
		}

		// Write to output directory
		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		fmt.Printf("Created: %s\n", dstPath)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src io.Reader, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
