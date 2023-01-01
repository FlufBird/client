package main

import (
	"github.com/FlufBird/client/packages/global/functions/logging"

	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Application struct {
	context context.Context
}

func createApplication() *Application {
	return &Application{}
}

func (application *Application) onStartup(context context.Context) {
	logging.Information("Frontend (Startup)", "Application started up.")

	application.context = context
}

func (application *Application) onDomReady(_ context.Context) {
	logging.Information("Frontend (DOM Ready)", "Application's DOM is ready.")

	runtime.WindowShow(application.context)
}

func displayUpdateDialog() {}