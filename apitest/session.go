/*
Apitest session
preferred way how to test rest api endpoints

Usage:

	session := apitest.NewSession()

new session with newly created user

	sudosession := apitest.NewSession().WithNewUser(func(user *models.User) {
		user.IsSuperuser = true
	})

session with loaded user
	usersession := apitest.NewSession().WithUser(user)

Perform request
FIrst we create prepared request
	request := session.Request("GET", "auth-user-login").JSONBody(map[string]string{})

then we can do this:
	request.Do().Response().Code
	request.Do().Scan(&response)

*/
package apitest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
)

/*
	Returns new Session
	if context is not given, patrol.Context is used
*/
func NewSession(context ...*context.Context) *Session {
	c := patrol.Context
	if len(context) > 0 {
		c = context[0]
	}
	return &Session{
		context: c,
	}
}

/*
	Session
	provides testing methods for rest api testing

*/
type Session struct {
	token   string
	user    *models.User
	context *context.Context
	err     error
}

/*
	Creates new user and attaches to session
*/
func (r *Session) WithNewUser(funcs ...func(user *models.User)) *Session {
	manager := models.NewUserManager(r.context)

	fs := []func(user *models.User){}
	fs = append(fs, manager.MakeTestUser(""))
	fs = append(fs, funcs...)

	user := manager.NewUser(fs...)
	r.err = user.Insert(r.context)

	return r.WithUser(user)
}

/*
	Attaches user to session
*/
func (s *Session) WithUser(user *models.User) *Session {
	manager := models.NewUserManager(s.context)
	s.user = user
	s.token, s.err = manager.Login(user)
	return s
}

// Returns last error
func (r *Session) Error() error {
	return r.err
}

// Sets token to session
func (r *Session) Token(token string) *Session {
	r.token = token
	return r
}

// Returns user
func (r *Session) User() *models.User {
	return r.user
}

/*
	Returns SessionRequest
*/
func (r *Session) Request(method, routename string, args ...string) *SessionRequest {
	url, _ := r.context.Router.Get(routename).URL(args...)

	rsr := &SessionRequest{
		method:  method,
		path:    url.Path,
		token:   r.token,
		context: r.context,
	}
	rsr.prepareRequest()
	return rsr
}

/*
	Prepared session request
*/
type SessionRequest struct {
	context  *context.Context
	request  *http.Request
	response *httptest.ResponseRecorder
	body     io.Reader
	method   string
	path     string
	err      error
	token    string
}

func (r *SessionRequest) Error() error {
	return r.err
}

func (r *SessionRequest) StringBody(body string) *SessionRequest {
	r.body = strings.NewReader(body)
	return r
}

func (r *SessionRequest) JSONBody(what interface{}) *SessionRequest {
	if b, err := json.Marshal(what); err == nil {
		r.body = bytes.NewBuffer(b)
	} else {
		r.err = errors.New("cannot marshal data")
		r.body = nil
	}
	return r
}

func (r *SessionRequest) Body(body io.Reader) *SessionRequest {
	r.body = body
	return r
}

func (r *SessionRequest) Request() *http.Request {
	return r.request
}

func (r *SessionRequest) Response() *httptest.ResponseRecorder {
	return r.response
}

func (r *SessionRequest) Do() *SessionRequest {
	r.prepareRequest()
	patrol.Context.Router.ServeHTTP(r.response, r.request)
	return r
}

func (r *SessionRequest) prepareRequest() {
	req, _ := http.NewRequest(r.method, r.path, r.body)
	r.request = req
	r.request.Header.Set("Authorization", "Bearer "+r.token)
	r.response = httptest.NewRecorder()
	r.err = nil
}

func (r *SessionRequest) Scan(target interface{}) *SessionRequest {
	if r.err != nil {
		return r
	}

	if err := json.Unmarshal(r.response.Body.Bytes(), target); err != nil {
		r.err = err
	}

	return r
}
