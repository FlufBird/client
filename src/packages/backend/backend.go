package backend

import (
	"github.com/FlufBird/client/src/packages/global/variables"

	"github.com/FlufBird/client/src/packages/global/functions/logging"

	"fmt"
	"time"
	"runtime"
	"os"
	"net/http"

	"github.com/juju/fslock"

	"github.com/Jeffail/gabs/v2"

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
		variables.RootDirectory = "../" // we're currently inside "src" in development so we have to move 1 directory above
		variables.Server = "http://localhost:5000"
	} else {
		variables.RootDirectory = ""
		variables.Server = "https://flufbird.deta.dev"
	}

	variables.Resources = fmt.Sprintf("%sresources", variables.RootDirectory)

	variables.ResourcesData = fmt.Sprintf("%s/data", variables.Resources)
	variables.ResourcesLanguages = fmt.Sprintf("%s/languages", variables.Resources)

	variables.Api = fmt.Sprintf("%s/api/v%s", variables.Server, variables.ApiVersion)

	variables.ApiUpdate = fmt.Sprintf("%s/update", variables.Api)

	parsedApplicationData, applicationDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/application.json", variables.ResourcesData))
	parsedUserData, userDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/user.json", variables.ResourcesData))

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

	languageData, languageDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/%s.json", variables.ResourcesLanguages, variables.Language))

	if languageDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve language data: %s", languageDataError)

		displayCriticalErrorDialog("Couldn't retrieve language data.")

		os.Exit(1)
	}

	variables.LanguageData = languageData

	variables.HttpClient = &http.Client{ // TODO: options
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

func checkUpdates() {}

func updateChecker() {}

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

	// TODO: run frontend (starting + application)

	go updateChecker()
}