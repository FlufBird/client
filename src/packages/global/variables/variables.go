package variables

import (
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

var DevelopmentMode bool

var Api string

var ApiUpdate string

var ApplicationData, UserData *gabs.Container
var Language *gabs.Container

var HttpClient *http.Client