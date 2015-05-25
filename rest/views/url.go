/*
URL is stolen concept from django where you can point to factory that produces
views (struct based)
*/
package views

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// factory that generates Viewer instances
type ViewerFactoryFunc func() Viewer

// Returns new url
func NewURL(url string, factory ViewerFactoryFunc) *URL {
	return &URL{
		url:     url,
		factory: factory,
	}
}

// URL for view
type URL struct {
	factory     ViewerFactoryFunc
	url         string
	name        string
	middlewares []alice.Constructor
}

// sets name
func (u *URL) Name(name string) *URL {
	u.name = name
	return u
}

// Returns name
func (u *URL) GetName() string {
	return u.name
}

// sets middlewares for URL
func (u *URL) Middlewares(middlewares ...alice.Constructor) *URL {
	u.middlewares = middlewares
	return u
}

func (u *URL) URL() string {
	return u.url
}

/*
Returns list of available methods on views
*/
func (u *URL) AvailableMethods() (result []string) {
	view := u.factory()
	result = []string{}
	for method, fn := range methods {
		if _, err := fn(view); err == nil {
			result = append(result, method)
		}
	}
	return
}

/*
Returns list of unavailable methods
*/
func (u *URL) UnavailableMethods() (result []string) {
	result = []string{}
	avail := u.AvailableMethods()
label:
	for method, _ := range methods {
		for _, am := range avail {
			if am == method {
				continue label
			}
		}
		result = append(result, method)
	}

	return
}

/*
Registers url to mux router
*/
func (u *URL) Register(router *mux.Router, chain alice.Chain) (err error) {

	availm := u.AvailableMethods()
	unavailm := u.UnavailableMethods()
	view := u.factory()

	// At leas one http method viewer interface must be satisfied
	if len(availm) == 0 {
		return fmt.Errorf("register %T failed: view must satisfy at least one <Method>Viewer interface.", view)
	}

	// build chain of middlewares for given view
	handlerchain := chain.Append(view.Middlewares()...).Append(u.middlewares...)

	// register to mux
	router.Handle(u.url, handlerchain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		// get fresh view instance from factory
		freshview := u.factory()

		// get method func
		method := methods[r.Method]

		// @TODO: this should never fail
		handler, _ := method(freshview)

		// run view.Before
		if err := freshview.Before(w, r); err != nil {

			// @TODO: what with errors in here, log them?
			if err != nil {
				glog.V(1).Infof("view.Before returned %v", err)
			}
			return
		}

		// process post
		handler(w, r)

		// process Before method
		if err := freshview.After(w, r); err != nil {
			return
		}
	})).Methods(availm...).Name(u.name)

	// register method not allowed
	route := router.Handle(u.url, chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		freshview := u.factory()

		// run method not allowed
		freshview.MethodNotAllowed(w, r)
	})).Methods(unavailm...)
	if u.name != "" {
		route.Name(u.name)
	}

	glog.V(1).Infof("patrol: register %T, url:\"%s\" methods:%v, NA: %v\n", view, u.url, availm, unavailm)

	return
}
