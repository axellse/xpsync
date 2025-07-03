package main

import (
	"embed"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	page := "/"
	width := 390
	height := 235

	serverErr := StartServer()
	if serverErr != nil {
		page = "/error.html"
		width = 313
		height = 125
	}

	err := wails.Run(&options.App{
		Title:  "LAN File Transfer Utility (XPSync)",
		Width:  width,
		Height: height,

		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Middleware: func(next http.Handler) http.Handler {
				return RouterMiddleWare{
					Page: page,
					NextHandler: next,
				}
			},
		},
		Frameless: true,
		BackgroundColour: &options.RGBA{R: 236, G: 233, B: 216, A: 0},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
