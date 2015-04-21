package static

import (
	"net/http"

	"github.com/stephens2424/bindataserver"
)

func HttpBindata() http.Handler {
	return bindataserver.Bindata(_bindata)
	// return bindataserver.Bindata(_bindata)
}

// Lists all files available
func ListFileNames() (result []string) {
	result = []string{}
	for f, _ := range _bindata {
		result = append(result, f)
	}

	return
}
