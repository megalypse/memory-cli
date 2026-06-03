package factory

import (
	"sync"

	"github.com/megalypse/memory_cli/internal/views"
)

var GetRoot = sync.OnceValue(getRoot)

func GetHeight() int {
	return GetRoot().Height()
}

func GetWidth() int {
	return GetRoot().Width()
}

func getRoot() *views.Root {
	return &views.Root{}
}
