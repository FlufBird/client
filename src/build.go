package main

import (
	"github.com/FlufBird/client/packages/global/variables"
	"github.com/FlufBird/client/packages/global/functions/logging"

	"embed"
	"fmt"
	"os"

	"net/http"

	"github.com/wailsapp/wails/v2"
	wailsOptions "github.com/wailsapp/wails/v2/pkg/options"
	assetServerOptions "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var frontend embed.FS

func assetServerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	file := request.URL.Path

	if file == "/favicon.ico" {
		writer.WriteHeader(http.StatusNotImplemented)

		return
	}

	prefix := ""

	if variables.Development {
		prefix = ".."
	}

	data, _error := os.ReadFile(fmt.Sprintf("%s%s", prefix, file))

	if _error != nil {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	writer.Write(data)
	writer.WriteHeader(http.StatusOK)
}

func buildFrontend() {
	logging.Information("Frontend (Build)", "Building frontend...")

	application := createApplication()

	_error := wails.Run(&wailsOptions.App{ // TODO: windows, linux
		Width: 300,
		Height: 400,

		MinWidth: 300,
		MinHeight: 400,

		Frameless: true,

		StartHidden: true,

		AssetServer: &assetServerOptions.Options{
			Assets:  frontend,
			Handler: http.HandlerFunc(assetServerHandler),
		},

		OnStartup: application.onStartup,
		OnDomReady: application.onDomReady,

		Bind: []interface{}{application},
	})

	if _error != nil {
		logging.Critical("Frontend", "Unable to build frontend: %s", _error)
	}
}