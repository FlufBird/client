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

	go setupEvents(application.context)
}

func (application *Application) onDomReady(_ context.Context) {
	logging.Information("Frontend", "Frontend's DOM is ready.")

	runtime.EventsEmit(application.context, "domReady")

	runtime.WindowCenter(application.context)
	runtime.WindowShow(application.context)
}

func setupEvents(context context.Context) {
	contentLoaded := false

	runtime.EventsOnce(context, "contentLoaded", func(_ ...interface{}) {contentLoaded = true})

	for {
		if contentLoaded {
			logging.Information("Frontend", "Frontend's content is loaded.")

			break
		}

		time.Sleep(3 * time.Second)
	}

	newUpdateAvailable, latestVersion, checkUpdatesError := checkUpdates(variables.ClientVersion, fmt.Sprintf("%s/update", variables.Api))

	if checkUpdatesError != nil {
		logging.Error("Update Checker (Frontend)", "Unable to check for updates: %s", checkUpdatesError)

		runtime.EventsEmit(context, "startupUpdateCheckerError")
	}

	if newUpdateAvailable {
		logging.Information("Update Checker (Frontend)", "New update available (v%s).", latestVersion)

		runtime.EventsEmit(context, "startupUpdateCheckerUpdateAvailable", []interface{}{latestVersion})
	} else {
		logging.Information("Update Checker (Frontend)", "Client is up-to-date (v%s).", variables.ClientVersion)

		runtime.EventsEmit(context, "startupUpdateCheckerUpToDate")
	}

	// go updateChecker(clientVersion, fmt.Sprintf("%s/update", variables.Api))
}

func (application *Application) GetLanguageData_(key string) string {
	return GetLanguageData(key)
}

func displayUpdateDialog() {}