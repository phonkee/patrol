package apitest

import (
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/settings"
)

var (
	once sync.Once
)

/*
Sets up test environment
*/
func Setup() {
	settings.SETTINGS_SECRET_KEY = os.Getenv("SECRET_KEY")

	once.Do(func() {
		if err := patrol.Setup(); err != nil && err != patrol.ErrPatrolAlreadySetup {
			glog.Errorf("patrol: setup error: %v", err)
		}
		patrol.Run([]string{"migrate"})
	})
}
