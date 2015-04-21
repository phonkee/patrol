package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
)

type APIResult struct {
	Result map[string]interface{} `json:"result"`
	Status int                    `json:"status"`
}

type APIResultList struct {
	Result []map[string]interface{} `json:"result"`
	Status int                      `json:"status"`
}

func NewAPIRequest(method, routename string, body io.Reader) *http.Request {
	url, _ := patrol.Context.Router.Get(routename).URL()
	r, _ := http.NewRequest(method, url.Path, body)
	return r
}

func NewAPIRequestWithToken(method, routename string, what interface{}, token string) *http.Request {
	var body io.Reader

	switch what := what.(type) {
	case io.Reader:
		body = what
	default:
		if b, err := json.Marshal(what); err == nil {
			body = bytes.NewBuffer(b)
		} else {
			body = bytes.NewBufferString(fmt.Sprintf("%v", b))
		}
	}

	r := NewAPIRequest(method, routename, body)
	r.Header.Set(settings.AUTH_TOKEN_HEADER_NAME, token)
	return r
}

func GetAPIResponse(r *http.Request, result interface{}) (*httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	patrol.Context.Router.ServeHTTP(w, r)

	if result != nil {
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			return w, err
		}
	}
	return w, nil
}

func NewBody(body string) io.Reader {
	return strings.NewReader(body)
}

func CreateUser(context *context.Context, funcs ...func(user *models.User)) (*models.User, error) {
	var err error
	manager := models.NewUserManager(context)

	rand.Seed(time.Now().UnixNano())

	/*
		@TODO: big todo, create user should use rest api. som this will be doing two things
		first to create superuser and login under this account
	*/

	funcs = append(funcs, func(user *models.User) {
		user.Email = utils.RandomString(20) + randomdata.Email()
		user.IsActive = true
		user.Username = utils.RandomString(20)
		user.Name = randomdata.FullName(randomdata.Male)
		user.SetPassword("")
	})

	user := manager.NewUser(funcs...)
	err = user.Insert(manager)
	return user, err
}

func CreateUserWithToken(context *context.Context, funcs ...func(user *models.User)) (user *models.User, token string, err error) {
	if user, err = CreateUser(context, funcs...); err != nil {
		return
	}

	manager := models.NewUserManager(context)
	token, err = manager.Login(user)

	return
}

/*
RequestsSession
*/
func NewRequestsSession(context *context.Context) *RequestsSession {
	return &RequestsSession{
		context: context,
	}
}

type RequestsSession struct {
	token   string
	user    *models.User
	context *context.Context
	err     error
}

func (r *RequestsSession) WithNewUser(funcs ...func(user *models.User)) *RequestsSession {
	manager := models.NewUserManager(r.context)

	fs := []func(user *models.User){}
	fs = append(fs, manager.MakeTestUser(""))
	fs = append(fs, funcs...)

	user := manager.NewUser(fs...)
	r.err = user.Insert(manager)
	r.user = user
	r.token, r.err = manager.Login(user)

	return r
}
func (r *RequestsSession) Error() error {
	return r.err
}

func (r *RequestsSession) Token(token string) *RequestsSession {
	r.token = token
	return r
}

func (r *RequestsSession) User() *models.User {
	return r.user
}

func (r *RequestsSession) Request(method, routename string, args ...string) *RequestsSessionRequest {
	url, _ := patrol.Context.Router.Get(routename).URL(args...)

	rsr := &RequestsSessionRequest{
		method:  method,
		path:    url.Path,
		token:   r.token,
		context: r.context,
	}
	rsr.prepareRequest()
	return rsr
}

type RequestsSessionRequest struct {
	context  *context.Context
	request  *http.Request
	response *httptest.ResponseRecorder
	body     io.Reader
	method   string
	path     string
	err      error
	token    string
}

func (r *RequestsSessionRequest) StringBody(body string) *RequestsSessionRequest {
	r.body = strings.NewReader(body)
	return r
}

func (r *RequestsSessionRequest) JSONBody(what interface{}) *RequestsSessionRequest {
	if b, err := json.Marshal(what); err == nil {
		r.body = bytes.NewBuffer(b)
	} else {
		r.err = errors.New("cannot marshal data")
		r.body = nil
	}
	return r
}

func (r *RequestsSessionRequest) Body(body io.Reader) *RequestsSessionRequest {
	r.body = body
	return r
}

func (r *RequestsSessionRequest) Request() *http.Request {
	return r.request
}

func (r *RequestsSessionRequest) Response() *httptest.ResponseRecorder {
	return r.response
}

func (r *RequestsSessionRequest) Do() *RequestsSessionRequest {
	r.prepareRequest()
	patrol.Context.Router.ServeHTTP(r.response, r.request)
	return r
}

func (r *RequestsSessionRequest) prepareRequest() {
	req, _ := http.NewRequest(r.method, r.path, r.body)
	r.request = req
	r.request.Header.Set(settings.AUTH_TOKEN_HEADER_NAME, r.token)
	r.response = httptest.NewRecorder()
	r.err = nil
}

func (r *RequestsSessionRequest) Scan(target interface{}) *RequestsSessionRequest {
	if r.err != nil {
		return r
	}

	if err := json.Unmarshal(r.response.Body.Bytes(), target); err != nil {
		r.err = err
	}

	return r
}

func (r *RequestsSessionRequest) Error() error {
	return r.err
}
