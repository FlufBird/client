package main

import (
	. "github.com/FlufBird/client/packages/global/functions/general"
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
	logging.Information("Frontend", "Application started up.")

	application.context = context
}

func (application *Application) onDomReady(_ context.Context) {
	logging.Information("Frontend", "Application's DOM is ready.")

	runtime.WindowCenter(application.context)
	runtime.WindowShow(application.context)
}

func (application *Application) GetLanguageData_(key string) string {
	return GetLanguageData(key).(string)
}

func displayUpdateDialog() {}