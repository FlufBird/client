package main

import (
	"github.com/FlufBird/client/packages/global/variables"

	. "github.com/FlufBird/client/packages/global/functions/general"
	"github.com/FlufBird/client/packages/global/functions/logging"

	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"net/http"

	"github.com/juju/fslock"

	"github.com/Jeffail/gabs/v2"

	"github.com/sqweek/dialog"
)

func setGlobalVariables() {
	var data string
	var server string

	displayDataRetrievalError := func() {displayCriticalErrorDialog("Couldn't retrieve data.")}

	variables.Development = true // dont forget to change this in production 😉

	switch variables.Development {
		case true:
			data = "../data"
			variables.Resources = "../resources"

			server = "http://localhost:31822"
		case false:
			data = "data"
			variables.Resources = "resources"

			server = "https://flufbird-api.deta.dev"
	}

	variables.TemporaryDirectory = os.TempDir()

	roamingAppDataDirectory, roamingAppDataDirectoryError := os.UserConfigDir()

	if roamingAppDataDirectoryError != nil {
		logging.Critical("Global Variables Setter", "Couldn't get Roaming AppData directory: %s", roamingAppDataDirectoryError)

		displayCriticalErrorDialog("Couldn't get Roaming AppData directory.")
	}

	variables.RoamingAppDataDirectory = roamingAppDataDirectory

	variables.DataDirectory = fmt.Sprintf("%s/flufbird", variables.RoamingAppDataDirectory)

	variables.ClientVersion = "1.0.0-a.1"

	variables.ApplicationData = fmt.Sprintf("%s/application", data)
	variables.UserData = fmt.Sprintf("%s/user", data)

	apiVersion := "1"

	variables.Api = fmt.Sprintf("%s/v%s", server, apiVersion)

	languages, languagesError := gabs.ParseJSONFile(fmt.Sprintf("%s/languages.json", variables.ApplicationData))

	if languagesError != nil {
		logging.Critical("Global Variables Setter", "Couldn't retrieve languages list: %s", languagesError)

		displayDataRetrievalError()
	}

	variables.Languages = languages

	generalUserData, generalUserDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/general.json", variables.UserData))

	if generalUserDataError != nil {
		logging.Critical("Global Variables Setter", "Couldn't retrieve general user data: %s", generalUserDataError)

		displayDataRetrievalError()
	}

	variables.GeneralUserData = generalUserData

	language, languageError := gabs.ParseJSONFile(fmt.Sprintf("%s/%s.json",
		fmt.Sprintf("%s/languages", variables.Resources),
		variables.GeneralUserData.Path("language").Data().(string),
	))

	if languagesError != nil {
		logging.Critical("Global Variables Setter", "Couldn't retrieve language data: %s", languageError)

		displayDataRetrievalError()
	}

	variables.Language = language

	variables.HttpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			MaxIdleConns: 0,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func checkInstances() {
	lock := fslock.New(fmt.Sprintf("%s/%s", variables.TemporaryDirectory, "flufbird_one_instance_lock"))
	_error := lock.TryLock()

	if _error != nil {
		logging.Information("Instances Checker", "Another instance is already running. | Lock Error: %s", _error.Error())

		displayInformationsDialog(GetLanguageData("general.onlyOneInstance.title"), GetLanguageData("general.onlyOneInstance.message"))

		os.Exit(0)
	}

	logging.Information("Instances Checker", "File locked.")
}

func checkUpdates(currentVersion string, route string) (bool, string, error) {
	response, requestError := variables.HttpClient.Get(fmt.Sprintf("%s/latest_version", route))

	if requestError != nil {
		return false, "", requestError
	}

	defer response.Body.Close()

	body, readError := io.ReadAll(response.Body)

	if readError != nil {
		return false, "", readError
	}

	data, parseError := gabs.ParseJSON(body)

	if parseError != nil {
		return false, "", parseError
	}

	if !data.Path("successful").Data().(bool) {
		return false, "", fmt.Errorf(data.Path("error").Data().(string))
	}

	latestVersion := data.Path("data.latestVersion").Data().(string)

	return currentVersion != latestVersion, latestVersion, nil
}

func updateChecker(currentVersion string, route string) {
	logging.Information("Update Checker (Backend)", "Checking for updates at %s", route)

	for {
		time.Sleep(time.Minute)

		newUpdateAvailable, latestVersion, _error := checkUpdates(currentVersion, route)

		if _error != nil {
			logging.Error("Update Checker (Backend)", "Unable to check for updates: %s", _error)
		}

		if newUpdateAvailable {
			logging.Information("Update Checker (Backend)", "New update available (v%s).", latestVersion)

			displayUpdateDialog()

			break
		}
	}
}

func displayDialog(title string, message string) *dialog.MsgBuilder {
	return dialog.Message(message).Title(title)
}

func displayInformationsDialog(title string, message string) {
	displayDialog(fmt.Sprintf("%s - FlufBird", title), message).Info()
}

func displayCriticalErrorDialog(message string) {
	displayDialog("Critical Error - FlufBird", message).Error()

	os.Exit(1)
}

func startBackend() {
	setGlobalVariables()

	logging.Information("General", "DEVELOPMENT BUILD | OS: %s | Architecture: %s", runtime.GOOS, runtime.GOARCH)

	if !variables.Development {
		checkInstances()
	}

	buildFrontend()
}