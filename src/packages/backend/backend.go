package backend

import (
	"github.com/FlufBird/client/src/packages/global/variables"

	"github.com/FlufBird/client/src/packages/global/functions/logging"

	frontend "github.com/FlufBird/client/src/packages/frontend/build"
	"github.com/FlufBird/client/src/packages/frontend/application"

	"errors"
	"fmt"
	"io"
	"runtime"
	"os"
	"time"
	"net/http"

	"github.com/juju/fslock"

	"github.com/Jeffail/gabs/v2"

	// "github.com/cavaliergopher/grab/v3"

	"github.com/sqweek/dialog"
)

func setVariables() {
	variables.DevelopmentMode = true

	variables.CurrentVersion = "1.0.0"
	variables.ApiVersion = "1"

	variables.RuntimeOS = runtime.GOOS
	variables.RuntimeArchitecture = runtime.GOARCH

	variables.TemporaryDirectory = os.TempDir()

	if variables.RuntimeOS != "windows" {
		variables.OldExecutable = "flufbird.old"
		variables.CurrentExecutable = "flufbird"
	} else {
		variables.OldExecutable = "flufbird.exe.old"
		variables.CurrentExecutable = "flufbird.exe"
	}

	if variables.DevelopmentMode {
		variables.Server = "http://localhost:5000"
	} else {
		variables.Server = "https://flufbird.deta.dev"
	}

	variables.Resources = "resources"
	variables.Data = "data"

	variables.Languages = fmt.Sprintf("%s/languages", variables.Resources)

	variables.Api = fmt.Sprintf("%s/api/v%s", variables.Server, variables.ApiVersion)

	variables.ApiUpdate = fmt.Sprintf("%s/update", variables.Api)

	parsedApplicationData, applicationDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/application.json", variables.Data))
	parsedUserData, userDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/user.json", variables.Data))

	if applicationDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve application data: %s", applicationDataError)

		displayCriticalErrorDialog("Couldn't retrieve application data.")

		os.Exit(1)
	}

	if userDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve user data: %s", userDataError)

		displayCriticalErrorDialog("Couldn't retrieve user data.")

		os.Exit(1)
	}

	variables.ApplicationData = parsedApplicationData
	variables.UserData = parsedUserData

	variables.Language = variables.UserData.Path("language").Data().(string)

	languageData, languageDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/%s.json", variables.Languages, variables.Language))

	if languageDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve language data: %s", languageDataError)

		displayCriticalErrorDialog("Couldn't retrieve language data.")

		os.Exit(1)
	}

	variables.LanguageData = languageData

	variables.HttpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { // don't allow redirects since our api doesn't redirect
			return errors.New("")
		},
	}
}

func deleteOldExecutable() {
	if os.Remove(variables.OldExecutable) == nil {
		logging.Information("Delete Old Executable", "Old executable deleted.")
	}
}

func checkInstances() {
	lock := fslock.New(fmt.Sprintf("%s/%s", variables.TemporaryDirectory, "flufbird_single_instance_check"))
	_error := lock.TryLock()

	if _error != nil {
		logging.Information("Check Instances", "Another instance is already running. | Lock Error: %s", _error.Error())

		displayInformationsDialog(variables.LanguageData.Path("general.onlyOneInstance.title").Data().(string), variables.LanguageData.Path("general.onlyOneInstance.message").Data().(string))

		os.Exit(0)
	}

	logging.Information("Check Instances", "File locked.")
}

func checkUpdates() bool {
	response, request_error := variables.HttpClient.Get(fmt.Sprintf("%s/latest_version", variables.ApiUpdate))

	if request_error != nil {
		logging.Information("Check Updates", "Couldn't send request.")

		return false
	}

	body, read_error := io.ReadAll(response.Body)

	if read_error != nil {
		logging.Information("Check Updates", "Couldn't read response.")

		return false
	}

	data, parse_error := gabs.ParseJSON(body)

	if parse_error != nil {
		logging.Information("Check Updates", "Couldn't parse response data.")

		return false
	}

	defer response.Body.Close()

	return variables.CurrentVersion != data.Path("latestVersion").Data().(string)
}

func updateChecker() {
	logging.Information("Updater", "Checking for updates at %s", variables.ApiUpdate)

	for {
		if checkUpdates() {
			logging.Information("Updater", "New update available, asking user.")

			if application.AskUpdate() {
				update()
			} else {
				logging.Information("Updater", "User denied.")
			}

			break
		}

		time.Sleep(30 * time.Second)
	}
}

func update() {
	logging.Information("Updater", "Got confirmation, updating...")

	// TODO: hide application (remember to hide all events) and display progress window, if this errors, display error, close the progress window and reshow the application
}

func displayDialog(title string, message string) *dialog.MsgBuilder {
	return dialog.Message(message).Title(title)
}

func displayInformationsDialog(title string, message string) {
	displayDialog(fmt.Sprintf("%s - FlufBird", title), message).Info()
}

func displayCriticalErrorDialog(message string) {
	displayDialog("Critical Error - FlufBird", message).Error()
}

func Backend() {
	setVariables()

	if variables.DevelopmentMode {
		fmt.Print("DEVELOPMENT MODE ENABLED\n\n")
	}

	logging.Information("Variables", "OS: %s | Architecture: %s", variables.RuntimeOS, variables.RuntimeArchitecture)
	logging.Information("Variables", "Temporary Directory: %s", variables.TemporaryDirectory)

	checkInstances()

	deleteOldExecutable()

	go updateChecker()

	frontend.Build()

	for { // prevents the program from exiting for development
		time.Sleep(time.Hour)
	}
}