package views

import "sync"

var GetRoot = sync.OnceValue(func() *Router {
	return &Router{}
})
