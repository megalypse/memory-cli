package factory

import (
	"sync"

	"github.com/megalypse/memory_cli/internal/views"
	"github.com/megalypse/memory_cli/internal/views/mainmenu"
)

var GetRouter = sync.OnceValue(getRouter)

func getRouter() *views.Router {
	return views.NewRouter(mainmenu.NewMainMenu(nil))
}
