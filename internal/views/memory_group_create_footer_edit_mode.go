package views

import "github.com/megalypse/go-cli-components/clicomponents"

func newMemoryGroupCreateFooterEditMode() *footer {
	return &footer{
		base: base{},
		Options: &clicomponents.CursorList{
			Items: []string{"(esc) Quit Edit Mode", "(tab) Switch Input"},
		},
	}
}
