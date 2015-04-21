package models

import (
	"net/http"
	"strings"
	"time"

	"github.com/phonkee/patrol/context"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/settings"
)

const (
	HTTP_STATUS_CONTEXT_KEY = "__patrol_http_status"
)

func NewRequestManager(context *context.Context) *RequestManager {
	return &RequestManager{
		context: context,
	}
}

/* Request manager
encapsulates various helpers for request
*/
type RequestManager struct {
	Manager
	context *context.Context
}

func (rm *RequestManager) Status(status int) {
	rm.context.Status = status
}

func (rm *RequestManager) GetStatus() int {
	return rm.context.Status
}

// Logs request
func (rm *RequestManager) LogRequest(r *http.Request, start time.Time) {
	status := rm.GetStatus()
	referer := r.Referer()
	if referer == "" {
		referer = "-"
	}

	userAgent := r.Header.Get("User-Agent")
	userAgent = ""

	glog.Infof("\"%s %s %s\" %d %v \"%s\" \"%s\"",
		r.Method, r.RequestURI, r.Proto, status, time.Since(start), referer,
		userAgent,
	)
}

// Logs request
func (rm *RequestManager) SentryAuthHeaders(r *http.Request) (values map[string]string, err error) {
	values = map[string]string{}
	auth := r.Header.Get(settings.SENTRY_AUTH_HEADER_NAME)
	index := strings.Index(auth, " ")
	if index == -1 {
		err = ErrCannotParseAuthHeaders
	} else {
		remainder := auth[index+1:]
		for _, part := range strings.Split(remainder, ",") {
			part = strings.TrimSpace(part)
			kv := strings.SplitN(part, "=", 2)
			values[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return
}
