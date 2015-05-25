/*
views.NewURL()
*/
package views

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/response"
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
	GetContext(r *http.Request) *context.Context
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

/*
Implementation of generic view
*/
type GenericView struct {
}

/*
Returns context view
*/
func (g *GenericView) GetContext(r *http.Request) (c *context.Context) {
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
APIView
*/
type APIView struct {
	GenericView
}

/*
Method not allowed implementation
*/
func (a *APIView) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusMethodNotAllowed).Write(w, r)
}
