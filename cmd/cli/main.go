package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/views"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Start Memory CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(views.GetRootInstance())

		_, err := p.Run()
		return err
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
