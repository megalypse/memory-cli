package views

import "github.com/megalypse/go-cli-components/clicomponents"

func newMemoryGroupCreateFooterEditMode() *footer {
	return &footer{
		base: base{},
		Options: &clicomponents.CursorListVertical{
			Items: []string{"(esc) Quit Edit Mode", "(tab) Switch Input"},
		},
	}
}
