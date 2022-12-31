package backend

import (
	"github.com/FlufBird/client/src/packages/global/variables"

	"github.com/FlufBird/client/src/packages/global/functions/logging"

	"github.com/FlufBird/client/src/packages/frontend/application"
	frontend "github.com/FlufBird/client/src/packages/frontend/build"

	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/juju/fslock"

	"github.com/Jeffail/gabs/v2"

	// "github.com/cavaliergopher/grab/v3"

	// "github.com/inconshreveable/go-update"

	"github.com/sqweek/dialog"
)

func setGlobalVariables() {
	var server string

	displayDataRetrievalError := func () {
		displayCriticalErrorDialog("Couldn't retrieve data.")
	}

	development := true // dont forget to change this in production 😉

	data := "data"

	variables.Development = development

	variables.Os = runtime.GOOS
	variables.Architecture = runtime.GOARCH

	variables.TemporaryDirectory = os.TempDir()

	variables.Resources = "resources"

	variables.ApplicationData = fmt.Sprintf("%s/application", data)
	variables.UserData = fmt.Sprintf("%s/user", data)

	variables.ApiVersion = "1"

	switch development {
		case true:
			server = "http://localhost:31822"
		case false:
			server = "https://flufbird.is-an.app"
	}

	variables.Api = fmt.Sprintf("%s/api/v%s", server, variables.ApiVersion)

	languages, languagesError := gabs.ParseJSONFile(fmt.Sprintf("%s/languages.json", variables.ApplicationData))

	if languagesError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve languages list: %s", languagesError)

		displayDataRetrievalError()
	}

	variables.Languages = languages

	generalUserData, generalUserDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/general.json", variables.UserData))

	if generalUserDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve general user data: %s", generalUserDataError)

		displayDataRetrievalError()
	}

	variables.GeneralUserData = generalUserData

	language, languageError := gabs.ParseJSONFile(fmt.Sprintf("%s/%s.json",
		fmt.Sprintf("%s/languages", variables.Resources),
		variables.GeneralUserData.Path("language").Data().(string),
	))

	if languagesError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve language data: %s", languageError)

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
		logging.Information("Check Instances", "Another instance is already running. | Lock Error: %s", _error.Error())

		displayInformationsDialog(variables.Language.Path("general.onlyOneInstance.title").Data().(string), variables.Language.Path("general.onlyOneInstance.message").Data().(string))

		os.Exit(0)
	}

	logging.Information("Check Instances", "File locked.")
}

func checkUpdates(currentVersion string, route string) bool {
	response, requestError := variables.HttpClient.Get(fmt.Sprintf("%s/latest_version", route))

	if requestError != nil {
		logging.Information("Updater (Check Updates)", "Couldn't send request.")

		return false
	}

	defer response.Body.Close()

	body, readError := io.ReadAll(response.Body)

	if readError != nil {
		logging.Information("Updater (Check Updates)", "Couldn't read response.")

		return false
	}

	data, parseError := gabs.ParseJSON(body)

	if parseError != nil {
		logging.Information("Updater (Check Updates)", "Couldn't parse response data.")

		return false
	}

	return currentVersion != data.Path("latestVersion").Data().(string)
}

func updateChecker(currentVersion string, route string) {
	logging.Information("Updater", "Checking for updates at %s", route)

	for {
		if checkUpdates(currentVersion, route) {
			logging.Information("Updater", "New update available.")

			application.DisplayUpdateDialog()

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

func Backend() {
	// clientVersion := "1.0.0-a.1"

	setGlobalVariables()

	if variables.Development {
		fmt.Print("IN DEVELOPMENT MODE\n\n")
	}

	logging.Information("General", "OS: %s | Architecture: %s", variables.Os, variables.Architecture)

	checkInstances(variables.TemporaryDirectory)

	frontend.Build()

	// TODO: check for updates on app start first, if the user clicks no, dont start the update checker thread

	// go updateChecker(clientVersion, fmt.Sprintf("%s/update", variables.Api))
}