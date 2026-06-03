package views

import "sync"

var GetRoot = sync.OnceValue(func() *Root {
	return &Root{}
})
