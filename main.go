package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"ggami-go/api"
	"ggami-go/internal/builder"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	app := NewApp()
	pm := builder.NewProjectManager()
	handler := api.NewHandler(pm)

	err := wails.Run(&options.App{
		Title:  "까미GO빌더",
		Width:  1400,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: handler,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
