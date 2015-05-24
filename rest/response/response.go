/*
Response is structured json response.

Response also supports errors, raw content, etc..

Most methods support chaining so writing response to http is doable on single line
In next examples we use these variables

	w http.ResponseWriter
	r *http.Request

Example of responses

	response.New().Error(errors.New("error")).Write(w, r)
	response.New().Error(structError).Write(w, r)

	response.New().Status(http.StatusOK).Result(product).Write(w, r)
	response.New().Status(http.StatusOK).Result(products).ResultSize(size).Write(w, r)

also there is non required argument status

	// raw supports string, Stringer, []byte, and other data are automatically
	// json marshalled
	response.New(http.StatusOK).Raw("{}")Write(w, r)

	body := map[string]string{
		"version": "1.0beta"
	}
	response.New(http.StatusOK).Raw(body)Write(w, r)
	response.New(http.StatusOK).Result(product).Write(w, r)
	response.New(http.StatusOK).Result(products).ResultSize(size).Write(w, r)
*/
package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	ctx "github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/ordering"
	"github.com/phonkee/patrol/rest/paginator"
	"github.com/phonkee/patrol/settings"
)

// Returns new response
func New(statuses ...int) (result *Response) {
	result = &Response{
		Headers_:     map[string]string{},
		ContentType_: "application/json",
	}
	//
	if len(statuses) > 0 {
		result.Status(statuses[0])
	}
	return
}

// Response
type Response struct {
	ContentType_ string      `json:"-"`
	Status_      int         `json:"status"`
	Message_     string      `json:"message"`
	Result_      interface{} `json:"result,omitempty"`
	ResultSize_  *int        `json:"result_size,omitempty"`

	Paging_   *paginator.Paginator `json:"paging,omitempty"`
	Ordering_ *ordering.Ordering   `json:"ordering,omitempty"`

	Error_   interface{}       `json:"error,omitempty"`
	Headers_ map[string]string `json:"-"`
	Raw_     *string           `json:"-"`
}

func (r *Response) ContentType(ct string) *Response {
	r.ContentType_ = ct
	return r
}

func (r *Response) Error(err interface{}) *Response {
	switch err := err.(type) {
	case error:
		r.Error_ = err.Error()
	case fmt.Stringer:
		r.Error_ = err.String()
	default:
		r.Error_ = err
	}
	return r
}

func (r *Response) Header(key, value string) *Response {
	r.Headers_[key] = value
	return r
}

// Write metadata (set raw content and return cors headers)
func (r *Response) Metadata(md *metadata.Metadata) *Response {
	r.Raw(md)
	r.Header("Access-Control-Allow-Origin", settings.CORS_ACCESS_ORIGIN)
	r.Header("Access-Control-Allow-Methods", strings.Join(md.Methods(), ", "))
	r.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
	return r
}

func (r *Response) Marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (r *Response) Message(message string) *Response {
	r.Message_ = message
	return r
}
func (r *Response) Ordering(ordering *ordering.Ordering) *Response {
	r.Ordering_ = ordering
	return r
}

func (r *Response) Paging(paging *paginator.Paginator) *Response {
	r.Paging_ = paging
	return r
}

func (r *Response) Raw(value interface{}) *Response {
	switch v := value.(type) {
	case fmt.Stringer:
		tmp := v.String()
		r.Raw_ = &tmp
	case string:
		r.Raw_ = &v
	case nil:
		r.Raw_ = nil
	case []byte:
		tmp := string(v)
		r.Raw_ = &tmp
	default:
		if json, err := r.Marshal(value); err == nil {
			tmp := string(json)
			r.Raw_ = &tmp
		} else {
			r.Raw_ = nil
			r.Error(err)
			r.Status(http.StatusInternalServerError)
		}
	}
	return r
}

func (r *Response) Result(result interface{}) *Response {
	r.Result_ = result
	return r
}
func (r *Response) ResultSize(size int) *Response {
	r.ResultSize_ = &size
	return r
}

// Sets status
func (r *Response) Status(status int) *Response {
	r.Status_ = status
	r.Message_ = http.StatusText(status)
	return r
}

// returns string representation
func (r *Response) String() (strbody string) {
	if r.Raw_ == nil {
		var (
			body []byte
			err  error
		)
		// marshal self
		if body, err = r.Marshal(r); err == nil {
			strbody = string(body)
		}
	} else {
		strbody = *r.Raw_
	}
	return
}

// Writes to response writer
func (rr *Response) Write(w http.ResponseWriter, r *http.Request) (err error) {

	var context *ctx.Context
	if context, err = ctx.Get(r); err != nil {
		return
	}

	// write headers
	w.Header().Set("Content-Type", rr.ContentType_)
	for k, v := range rr.Headers_ {
		w.Header().Set(k, v)
	}

	// if not status set we set from context
	if rr.Status_ == 0 {
		rr.Status(context.Status)
	} else {
		context.Status = rr.Status_
	}

	w.WriteHeader(rr.Status_)

	fmt.Fprint(w, rr.String())
	return
}
