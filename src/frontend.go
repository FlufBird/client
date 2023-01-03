package main

import (
	"github.com/FlufBird/client/packages/global/variables"

	. "github.com/FlufBird/client/packages/global/functions/general"
	"github.com/FlufBird/client/packages/global/functions/logging"

	"fmt"
	"context"
	"time"

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

func (application *Application) onDomReady(_ context.Context) { // FIXME: ensure EventsOnce is executed before EventsEmit
	logging.Information("Frontend", "Frontend's DOM is ready.")

	runtime.EventsEmit(application.context, "domReady")

	runtime.WindowCenter(application.context)
	runtime.WindowShow(application.context)

	contentLoaded := false

	runtime.EventsOnce(application.context, "contentLoaded", func(_ ...interface{}) {contentLoaded = true})

	for {
		if contentLoaded {
			logging.Information("Frontend", "Frontend's content is loaded.")

			break
		}

		time.Sleep(2 * time.Second)
	}

	newUpdateAvailable, checkUpdatesError := checkUpdates(variables.ClientVersion, fmt.Sprintf("%s/update", variables.Api))

	if checkUpdatesError != nil {
		logging.Information("Update Checker (Frontend)", "Unable to check for updates: %s", checkUpdatesError)

		runtime.EventsEmit(application.context, "startupUpdateCheckerError")
	}

	if newUpdateAvailable {
		logging.Information("Update Checker (Frontend)", "New update available.")

		runtime.EventsEmit(application.context, "startupUpdateCheckerUpdateAvailable", []interface{}{"oldversion", "newversion"})

		// TODO: ask user & redirect, stop app if accepted
	}
}

func (application *Application) GetLanguageData_(key string) string {
	return GetLanguageData(key)
}

func displayUpdateDialog() {}