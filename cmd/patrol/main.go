package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/settings"
)

func main() {

	/* In custom builds of patrol, custom plugins are registered before Setup
	 */

	// Setup patrol application
	if err := patrol.Setup(); err != nil {
		fmt.Printf(settings.RED_COLOR("patrol: setup failed with error: %s\n"), err)
		os.Exit(1)
	}

	if err := patrol.Run(flag.Args()); err != nil {
		fmt.Printf(settings.RED_COLOR("error: %s\n"), err)
		os.Exit(1)
	}
}
