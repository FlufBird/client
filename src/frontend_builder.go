package main

import (
	"github.com/FlufBird/client/packages/global/functions/logging"
	"github.com/FlufBird/client/packages/global/variables"

	"embed"
	"fmt"
	"os"
	"strings"

	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	assetServerOptions "github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed frontend/dist
var assets embed.FS

func assetServerHandler(writer http.ResponseWriter, request *http.Request) { // how to access local file? ~ here you go 😊:
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	asset := strings.TrimPrefix(request.URL.Path, "/")

	if asset == "favicon.ico" {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	prefix := ""

	if strings.HasPrefix(asset, "resources") || strings.HasPrefix(asset, "data") {
		if variables.Development {
			prefix = ".."
		}
	} else if asset == "wailsjs/go/main/Application" {
		if variables.Development {
			prefix = "frontend"
		}

		asset += ".js"
	}

	data, _error := os.ReadFile(fmt.Sprintf("%s/%s", prefix, asset))

	if _error != nil {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	writer.Write(data)

	if strings.HasSuffix(asset, ".js") {
		writer.Header().Set("Content-Type", "text/javascript")
	}

	writer.WriteHeader(http.StatusOK)
}

func buildFrontend() {
	logging.Information("Frontend", "Building frontend...")

	application := createApplication()

	const width int = 800;
	const height int = 600;

	_error := wails.Run(&options.App{ // TODO: icon
		Title: "FlufBird",

		Width: width,
		Height: height,

		MinWidth: width,
		MinHeight: height,

		StartHidden: true,

		Windows: &windows.Options{
			WebviewUserDataPath: fmt.Sprintf("%s/webview_flufbird", variables.RoamingAppDataDirectory),
		},
		Linux: &linux.Options{
			Icon: []byte{}, // TODO: icon
		},

		AssetServer: &assetServerOptions.Options{
			Assets: assets,
			Handler: http.HandlerFunc(assetServerHandler),
		},

		OnStartup: application.onStartup,

		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},

		Bind: []interface{}{application},
	})

	if _error != nil {
		logging.Critical("Frontend", "Unable to build frontend: %s", _error)
	}
}