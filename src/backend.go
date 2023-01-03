package main

import (
	"github.com/FlufBird/client/packages/global/variables"

	. "github.com/FlufBird/client/packages/global/functions/general"
	"github.com/FlufBird/client/packages/global/functions/logging"

	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

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

			server = "https://flufbird.is-an.app"
	}

	variables.TemporaryDirectory = os.TempDir()

	variables.ClientVersion = "1.0.0-a.1"

	variables.ApplicationData = fmt.Sprintf("%s/application", data)
	variables.UserData = fmt.Sprintf("%s/user", data)

	variables.ApiVersion = "1"

	variables.Api = fmt.Sprintf("%s/api/v%s", server, variables.ApiVersion)

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

func checkInstances(temporaryDirectory string) {
	lock := fslock.New(fmt.Sprintf("%s/%s", temporaryDirectory, "flufbird_one_instance_lock"))
	_error := lock.TryLock()

	if _error != nil {
		logging.Information("Instances Checker", "Another instance is already running. | Lock Error: %s", _error.Error())

		displayInformationsDialog(GetLanguageData("general.onlyOneInstance.title"), GetLanguageData("general.onlyOneInstance.message"))

		os.Exit(0)
	}

	logging.Information("Instances Checker", "File locked.")
}

func checkUpdates(currentVersion string, route string) (bool, error) {
	response, requestError := variables.HttpClient.Get(fmt.Sprintf("%s/latest_version", route))

	if requestError != nil {
		logging.Information("Update Checker", "Couldn't send request: %s", requestError)

		return false, requestError
	}

	defer response.Body.Close()

	body, readError := io.ReadAll(response.Body)

	if readError != nil {
		logging.Information("Update Checker", "Couldn't read response: %s", readError)

		return false, readError
	}

	data, parseError := gabs.ParseJSON(body)

	if parseError != nil {
		logging.Information("Update Checker", "Couldn't parse response data: %s", parseError)

		return false, parseError
	}

	return currentVersion != data.Path("latestVersion").Data().(string), nil
}

func updateChecker(currentVersion string, route string) {
	logging.Information("Update Checker (Backend)", "Checking for updates at %s", route)

	for {
		if newUpdateAvailable, _ := checkUpdates(currentVersion, route); newUpdateAvailable {
			logging.Information("Update Checker (Backend)", "New update available.")

			displayUpdateDialog()

			break
		}

		time.Sleep(30 * time.Second)
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

	if variables.Development {
		fmt.Print("IN DEVELOPMENT MODE\n\n")
	}

	logging.Information("General", "OS: %s | Architecture: %s", runtime.GOOS, runtime.GOARCH)

	if !variables.Development {
		checkInstances(variables.TemporaryDirectory)
	}

	buildFrontend()

	// go updateChecker(clientVersion, fmt.Sprintf("%s/update", variables.Api))
}