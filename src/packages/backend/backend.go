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

	"github.com/sqweek/dialog"
)

func setGlobalVariables(
	resources string,
	data string,

	server string,

	apiVersion string,
) {
	variables.DevelopmentMode = true

	variables.Api = fmt.Sprintf("%s/api/v%s", server, apiVersion)

	variables.ApiUpdate = fmt.Sprintf("%s/update", variables.Api)

	// parsedApplicationData, applicationDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/application.json", data))
	parsedUserData, userDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/user.json", data))

	// if applicationDataError != nil {
	// 	logging.Fatal("Variables Setting", "Couldn't retrieve application data: %s", applicationDataError)

	// 	displayCriticalErrorDialog("Couldn't retrieve application data.")

	// 	os.Exit(1)
	// }

	if userDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve user data: %s", userDataError)

		displayCriticalErrorDialog("Couldn't retrieve user data.")

		os.Exit(1)
	}

	language, languageDataError := gabs.ParseJSONFile(fmt.Sprintf("%s/%s.json", fmt.Sprintf("%s/languages", resources), parsedUserData.Path("language").Data().(string)))

	if languageDataError != nil {
		logging.Fatal("Variables Setting", "Couldn't retrieve language data: %s", languageDataError)

		displayCriticalErrorDialog("Couldn't retrieve language data.")

		os.Exit(1)
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

func checkUpdated(
	updateArchive string,

	oldExecutable string,
	currentExecutable string,

	resources string,
	newResources string,

	data string,
	newData string,
) {
	if os.Remove(oldExecutable) == nil {
		logging.Information("Check Updated", "Old executable deleted.")
	}

	// if os.RemoveAll(oldResources) == nil {
	// 	logging.Information("Check Updated", "Old resources deleted.")
	// }

	// if os.RemoveAll(oldData) == nil {
	// 	logging.Information("Check Updated", "Old data deleted.")
	// }
}

func checkInstances(temporaryDirectory string) {
	lock := fslock.New(fmt.Sprintf("%s/%s", temporaryDirectory, "flufbird_single_instance_lock"))
	_error := lock.TryLock()

	if _error != nil {
		logging.Information("Check Instances", "Another instance is already running. | Lock Error: %s", _error.Error())

		displayInformationsDialog(variables.Language.Path("general.onlyOneInstance.title").Data().(string), variables.Language.Path("general.onlyOneInstance.message").Data().(string))

		os.Exit(0)
	}

	logging.Information("Check Instances", "File locked.")
}

func checkUpdates(clientVersion string) bool {
	response, requestError := variables.HttpClient.Get(fmt.Sprintf("%s/latest_version", variables.ApiUpdate))

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

	return clientVersion != data.Path("latestVersion").Data().(string)
}

func updateChecker(clientVersion string) {
	logging.Information("Updater", "Checking for updates at %s", variables.ApiUpdate)

	for {
		if checkUpdates(clientVersion) {
			logging.Information("Updater", "New update available, asking user.")

			if application.AskUpdate() {
				logging.Information("Updater", "Got confirmation.")

				update()
			} else {
				logging.Information("Updater", "User denied, not checking for updates anymore.")
			}

			break
		}

		time.Sleep(30 * time.Second)
	}
}

func update() {
	// TODO: hide application (remember to hide all events) and display progress window, if this errors, display error, close the progress window and reshow the application

	logging.Information("Updater", "Downloading update archive...")

	logging.Information("Updater", "Unarchiving update archive...")

	logging.Information("Updater", "Exiting.")

	os.Exit(0)
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
	var oldExecutable, currentExecutable string
	var server string

	runtimeOS := runtime.GOOS
	runtimeArchitecture := runtime.GOARCH

	temporaryDirectory := os.TempDir()

	resources := "resources"
	newResources := "new_resources"

	data := "data"
	newData := "new_data"

	clientVersion := "1.0.0"
	apiVersion := "1"

	switch runtimeOS {
		case "windows":
			oldExecutable = "flufbird.exe.old"
			currentExecutable = "flufbird.exe"
		default:
			oldExecutable = "flufbird.old"
			currentExecutable = "flufbird"
	}

	updateArchive := "update.zip"

	switch variables.DevelopmentMode {
		case true:
			server = "http://localhost:5000"
		case false:
			server = "https://flufbird.is-an.app"
	}
	// end setting variables

	checkUpdated(
		updateArchive,

		oldExecutable,
		currentExecutable,

		resources,
		newResources,

		data,
		newData,
	)

	setGlobalVariables(
		resources,
		data,

		server,

		apiVersion,
	)

	if variables.DevelopmentMode {
		fmt.Print("DEVELOPMENT MODE ENABLED\n\n")
	}

	logging.Information("General", "OS: %s | Architecture: %s", runtimeOS, runtimeArchitecture)

	checkInstances(temporaryDirectory)

	go updateChecker(clientVersion)

	frontend.Build()

	for { // prevents the program from exiting for development, we dont yet have the application
		time.Sleep(time.Hour)
	}
}