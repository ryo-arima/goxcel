package controller

import (
	"errors"
	"fmt"
	"os"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/usecase"
	"github.com/spf13/cobra"
)

// InitFormatCmd creates the 'format' subcommand which pretty-prints a .gxl file via usecase -> repository.
func InitFormatCmd() *cobra.Command {
	var (
		inPlace bool
		output  string
	)

	cmd := &cobra.Command{
		Use:   "format <template.gxl>",
		Short: "Format a .gxl template (pretty-print XML)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inPlace && output != "" {
				return errors.New("cannot use --write and --output together")
			}
			path := args[0]

			conf := config.NewBaseConfigWithFile(path)
			u := usecase.NewFormatUsecase(conf)
			formatted, err := u.Format(path)
			if err != nil {
				return fmt.Errorf("format: %w", err)
			}

			if inPlace {
				if err := os.WriteFile(path, formatted, 0644); err != nil {
					return fmt.Errorf("write file: %w", err)
				}
				return nil
			}
			if output != "" {
				if err := os.WriteFile(output, formatted, 0644); err != nil {
					return fmt.Errorf("write output: %w", err)
				}
				return nil
			}
			_, err = os.Stdout.Write(formatted)
			return err
		},
	}

	cmd.Flags().BoolVarP(&inPlace, "write", "w", false, "write result to (source) file instead of stdout")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to file path")
	return cmd
}
