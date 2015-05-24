package apitest

import (
	"os"

	"github.com/golang/glog"
	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/settings"
)

/*
Sets up test environment
*/
func Setup() {
	settings.SETTINGS_SECRET_KEY = os.Getenv("SECRET_KEY")

	if err := patrol.Setup(); err != nil && err != patrol.ErrPatrolAlreadySetup {
		glog.Errorf("patrol: setup error: %v", err)
	}

	patrol.Run([]string{"migrate"})
}
