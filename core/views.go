package core

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/types"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const (
	API_JSON_RESULT_KEY = "result"
)

var (
	// @TODO: move this to settings
	flagCorsOrigin = flag.String("cors_allowed_origin", "*", "Cors allowed origin (Access-Control-Allow-Origin header).")

	ErrBreakRequest   = errors.New("")
	ErrParamNotFound  = errors.New("param not found")
	ErrMethodNotFound = errors.New("method not found")
	ErrMuxVarNotFound = errors.New("mux var not found")

	// mapping of methods
	methods = map[string]func(Viewer) (http.HandlerFunc, error){
		"GET": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(GetViewer); ok {
				return t.GET, nil
			}
			return nil, ErrMethodNotFound
		},
		"POST": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(PostViewer); ok {
				return t.POST, nil
			}
			return nil, ErrMethodNotFound
		},
		"PUT": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(PutViewer); ok {
				return t.PUT, nil
			}
			return nil, ErrMethodNotFound
		},
		"PATCH": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(PatchViewer); ok {
				return t.PATCH, nil
			}
			return nil, ErrMethodNotFound
		},
		"DELETE": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(DeleteViewer); ok {
				return t.DELETE, nil
			}
			return nil, ErrMethodNotFound
		},
		"OPTIONS": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(OptionsViewer); ok {
				return t.OPTIONS, nil
			}
			return nil, ErrMethodNotFound
		},
		"HEAD": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(HeadViewer); ok {
				return t.HEAD, nil
			}
			return nil, ErrMethodNotFound
		},
	}
)

/*
Dispatch if defined is called before calling any http method
The only difference is that Dispatch can return error, in which case
handler is not called.
*/
type Viewer interface {
	// Method Before runs before actual http method
	// if it returns error http method is not called
	// can be used to terminate (not authorized).
	// In this case no response is written
	Before(w http.ResponseWriter, r *http.Request) error

	// called after http method call
	After(w http.ResponseWriter, r *http.Request) error

	// returns list of middlewares in which it's wrapped
	Middlewares() []alice.Constructor

	// method not allowed implementation
	MethodNotAllowed(w http.ResponseWriter, r *http.Request)

	// returns context by request
	Context(r *http.Request) *context.Context
}

/*
View that supports http DELETE method
*/
type DeleteViewer interface {
	DELETE(w http.ResponseWriter, r *http.Request)
}

/*
View that supports http GET method
*/
type GetViewer interface {
	GET(w http.ResponseWriter, r *http.Request)
}

/*
View that supports http HEAD method
*/
type HeadViewer interface {
	HEAD(w http.ResponseWriter, r *http.Request)
}

/*
View that supports http OPTIONS method
*/
type OptionsViewer interface {
	OPTIONS(w http.ResponseWriter, r *http.Request)
}

/*
View that supports http PATCH method
*/
type PatchViewer interface {
	PATCH(w http.ResponseWriter, r *http.Request)
}

/*
View that supports http POST method
*/
type PostViewer interface {
	POST(w http.ResponseWriter, r *http.Request)
}

/*
View that supports http PUT method
*/
type PutViewer interface {
	PUT(w http.ResponseWriter, r *http.Request)
}

type ViewerFactoryFunc func() Viewer

/*
   URLView
*/
type URLView struct {
	factory     ViewerFactoryFunc
	view        Viewer
	url         string
	name        string
	middlewares []alice.Constructor
}

func NewURLView(url string, vff ViewerFactoryFunc) *URLView {
	return &URLView{factory: vff, url: url}
}

func (u *URLView) Name(name string) *URLView {
	u.name = name
	return u
}
func (u *URLView) URL() string { return u.url }
func (u *URLView) Middlewares(middlewares ...alice.Constructor) *URLView {
	u.middlewares = middlewares
	return u
}

func (u *URLView) GetName() string {
	return u.name
}

func (u *URLView) Methods() (result []string) {
	view := u.factory()
	result = []string{}
	for method, fn := range methods {
		if _, err := fn(view); err == nil {
			result = append(result, method)
		}
	}
	return
}
func (u *URLView) NotImplementedMethods() (result []string) {
	view := u.factory()
	result = []string{}
	for method, fn := range methods {
		if _, err := fn(view); err != nil {
			result = append(result, method)
		}
	}
	return
}

func (u *URLView) Register(router *mux.Router, chain alice.Chain) (err error) {

	availMethods := u.Methods()
	niMethods := u.NotImplementedMethods()
	view := u.factory()

	if len(availMethods) == 0 {
		return fmt.Errorf("register %T failed: view must satisfy at least one <Method>Viewer interface.", view)
	}

	// build chain of middlewares for given view
	handlerChain := chain.Append(view.Middlewares()...).Append(u.middlewares...)

	// register to mux
	router.Handle(u.url, handlerChain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		// get fresh view instance from factory
		curView := u.factory()

		// get method func
		method := methods[r.Method]
		// @TODO: this should never fail
		handler, _ := method(curView)

		// run view.Before
		if err := curView.Before(w, r); err != nil {
			if err != ErrBreakRequest {
				// glog.Error(err)
			}
			return
		}

		// process post
		handler(w, r)

		// process Before method
		if err := curView.After(w, r); err != nil {
			return
		}
	})).Methods(availMethods...).Name(u.name)

	// register method not allowed
	router.Handle(u.url, chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		curView := u.factory()

		// run method not allowed
		curView.MethodNotAllowed(w, r)

	})).Methods(niMethods...).Name(u.name)

	glog.V(1).Infof("patrol: register %T, url:\"%s\" methods:%v, NA: %v\n", view, u.url, availMethods, niMethods)

	return
}

/*
   Base views
*/

/*
Generic view
provides some helpers above response (status)
*/
type GenericView struct{}

func (g *GenericView) Context(r *http.Request) (c *context.Context) {
	var err error
	if c, err = context.Get(r); err != nil {
		panic("context is not there")
	}

	return
}

/*
Method Before is called before calling view handler. One difference is
that Method Before can return an error. If error is returned view handler
is not called
*/
func (g GenericView) Before(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (g GenericView) After(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (g GenericView) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusMethodNotAllowed).Raw("method not allowed").Write(w, r)
}

/*
Returns chain of middlewares for all handlers that view supports
*/
func (g GenericView) Middlewares() []alice.Constructor {
	return []alice.Constructor{}
}

/*
returns mux variable as int64
*/
func (g *GenericView) GetMuxVarInt64(r *http.Request, name string) (value int64, err error) {
	vars := mux.Vars(r)
	stringvar, ok := vars[name]
	if !ok {
		return value, ErrMuxVarNotFound
	}

	if value, err = strconv.ParseInt(stringvar, 10, 0); err != nil {
		return
	}

	return
}

func (g *GenericView) GetMuxVarForeignKey(r *http.Request, name string) (value *types.ForeignKey, err error) {
	var intvar int64

	if intvar, err = g.GetMuxVarInt64(r, name); err != nil {
		return
	}
	fk := types.ForeignKey(intvar)
	return &fk, nil
}

/*JSONView
 */
type JSONView struct {
	GenericView
}

/* Helper to unmarshal json to target
 */
func (j *JSONView) Unmarshal(body io.Reader, target interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(target)
}

func (j *JSONView) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusMethodNotAllowed).Write(w, r)
}

func (j *JSONView) GetParam(r *http.Request, param string) (value string, ok bool) {
	vars := mux.Vars(r)
	value, ok = vars[param]
	return
}

func (j *JSONView) GetParamInt(r *http.Request, param string) (id int, err error) {
	var (
		value string
		ok    bool
	)

	if value, ok = j.GetParam(r, param); !ok {
		err = ErrParamNotFound
		return
	}

	if id, err = strconv.Atoi(value); err != nil {
		return
	}

	return
}
