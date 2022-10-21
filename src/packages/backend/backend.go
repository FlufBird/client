package backend

import (
	"github.com/FlufBird/client/src/packages/global/variables"

	"github.com/FlufBird/client/src/packages/global/functions/logging"

	"fmt"
	"os"

	// "http/client"

	"github.com/juju/fslock"

	"github.com/Jeffail/gabs/v2"

	"github.com/sqweek/dialog"
)

func setVariables() {
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

	variables.Api = fmt.Sprintf("%s/api/v%s", variables.Server, variables.ApiVersion)

	parsedApplicationData, applicationDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/application.json", variables.ResourcesData))
	parsedUserData, userDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/user.json", variables.ResourcesData))

	if applicationDataError != nil {
		displayCriticalErrorDialog("Couldn't retrieve application data.")
	}

	if userDataError != nil {
		displayCriticalErrorDialog("Couldn't retrieve user data.")
	}

	variables.ApplicationData = parsedApplicationData
	variables.UserData = parsedUserData

	// TODO: http client
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

		displayInformationsDialog("Only 1 Instance of FlufBird can be Running at A Time", "Another instance of FlufBird is already running, only 1 instance of FlufBird can be running at a time!")

		os.Exit(0)
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
}

func Backend() {
	logging.Information("General", "OS: %s | Architecture: %s", variables.RuntimeOS, variables.RuntimeArchitecture)
	logging.Information("General", "Temporary Directory: %s", variables.TemporaryDirectory)

	setVariables()
	checkInstances()
	deleteOldExecutable()

	// TODO: run frontend (starting + application)
	// TODO: start update checker (thread)
}