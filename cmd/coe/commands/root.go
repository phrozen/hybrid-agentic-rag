package commands

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra command gateway.
var RootCmd = &cobra.Command{
	Use:   "coe",
	Short: "Chronicles of Aethelgard Hybrid RAG CLI",
	Long:  `A high-performance CLI utility for hybrid semantic and keyword search operations on Chronicles of Aethelgard.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Silence all log/slog logs globally across all subcommands
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	},
}

// Execute triggers standard routing of subcommands.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Execution failed: %v\n", err)
		os.Exit(1)
	}
}
