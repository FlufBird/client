package variables

import (
	"fmt"
	"os"
	"runtime"

	"net/http"
)

const RuntimeOS string = runtime.GOOS
const RuntimeArchitecture string = runtime.GOARCH

const resources string = "resources"

const CurrentVersion string = "1.0.0"
const ApiVersion string = "1"

const DevelopmentMode bool = true

var OldExecutable, CurrentExecutable string
var Server, Api string

var ResourcesData string = fmt.Sprintf("%s/data", resources)
var ResourcesLanguages string = fmt.Sprintf("%s/languages", resources)

var ApiUpdate string = fmt.Sprintf("%s/update", Api)

var ApplicationData, UserData interface{}

var TemporaryDirectory string = os.TempDir()

var HttpClient http.Client