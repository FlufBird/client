package main

import (
	"github.com/FlufBird/client/packages/global/functions/logging"

	"embed"

	"github.com/wailsapp/wails/v2"
	wailsOptions "github.com/wailsapp/wails/v2/pkg/options"
	assetServerOptions "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var frontend embed.FS

func buildFrontend() {
	logging.Information("Frontend (Build)", "Building frontend...")

	application := createApplication()

	_error := wails.Run(&wailsOptions.App{ // TODO: windows, linux
		MinWidth: 800,
		MinHeight: 600,

		Frameless: true,

		StartHidden: true,

		WindowStartState: 1,

		AssetServer: &assetServerOptions.Options{Assets: frontend},

		OnStartup: application.onStartup,
		OnDomReady: application.onDomReady,

		Bind: []interface{}{application},
	})

	if _error != nil {
		logging.Critical("Frontend (Build)", "Unable to build frontend: %s", _error)
	}
}