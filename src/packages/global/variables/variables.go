package variables

import (
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

var DevelopmentMode bool

var Os, Architecture string

var TemporaryDirectory string

var Resources string
var ApplicationData, UserData string

var GeneralUserData *gabs.Container

var Languages, Language *gabs.Container

var ApiVersion, Api string

var HttpClient *http.Client