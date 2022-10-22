package variables

import (
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

var DevelopmentMode bool

var RootDirectory string

var RuntimeOS string
var RuntimeArchitecture string

var TemporaryDirectory string

var Resources string

var CurrentVersion string
var ApiVersion string

var OldExecutable, CurrentExecutable string
var Server, Api string

var ResourcesData string
var ResourcesLanguages string

var ApiUpdate string

var ApplicationData, UserData *gabs.Container

var Language string
var LanguageData *gabs.Container

var HttpClient *http.Client