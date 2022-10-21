package main

import (
	"github.com/FlufBird/client/src/packages/global/variables"

	"github.com/FlufBird/client/src/packages/backend"

	"fmt"
)

func main() {
	if variables.DevelopmentMode {
		fmt.Print("LOGGING ENABLED\n\n")
	}

	backend.Backend()
}