package general

import "github.com/FlufBird/client/packages/global/variables"

func GetLanguageData(key string) interface{} {
	return variables.Language.Path(key).Data()
}