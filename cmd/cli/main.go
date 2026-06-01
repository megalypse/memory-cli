package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type SampleModel struct{}

func (s SampleModel) Init() tea.Cmd {
	panic("implement me")
}

func (s SampleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	panic("implement me")
}

func (s SampleModel) View() string {
	panic("implement me")
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Start Memory CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(SampleModel{})

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
