package command

import (
	"fmt"
	"os"

	"github.com/ryo-arima/goxcel/pkg/controller"
	"github.com/spf13/cobra"
)

// NewRootCmd builds the root CLI command and wires subcommands.
func NewRootCmd() *cobra.Command {
	var root = &cobra.Command{
		Use:   "goxcel",
		Short: "Generate .xlsx from .gxl templates",
		Long:  "goxcel is a CLI to render .gxl templates with data into Excel .xlsx files.",
	}

	// Subcommands
	root.AddCommand(controller.InitGenerateCmd())
	return root
}

// Execute runs the CLI.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
