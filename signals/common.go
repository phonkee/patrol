package signals

import "time"

/* Signal handler which notifies on http server start
When http server starts this signal handler it's called. There may be some
scenarios where this can be used e.g. monitoring.
*/
type OnHttpServerStartSignalHandler interface {
	OnHttpServerStart()
}

type OnCleanupSignalHandler interface {
	OnCleanup(time.Time)
}
