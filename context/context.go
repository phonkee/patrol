package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/phonkee/ergoq"
	"github.com/phonkee/gocacher"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
)

const (
	CONTEXT_KEY = "patrol:context"
	SECRET_KEY  = "patrol:secretkey"
)

var (
	ErrContextNotFound       = errors.New("context_not_found")
	ErrTransactionNotStarted = errors.New("transaction not started")

	mutex sync.RWMutex
	data  = make(map[*http.Request]*Context)
)

// Returns new context
func New(DatabaseDSN, MessageQueueDSN, CacheDSN string) (result *Context, err error) {
	result = &Context{}
	result.Router = mux.NewRouter()
	result.Router.StrictSlash(true)
	result.Vars = map[interface{}]interface{}{}
	result.Status = http.StatusTeapot

	// This was not ok
	// result.Request, _ = http.NewRequest("GET", "/", nil)

	funcs := []func() error{
		func() error { return result.dialDB(DatabaseDSN) },
		func() error { return result.dialMQ(MessageQueueDSN) },
		func() error { return result.dialCache(CacheDSN) },
	}

	for _, f := range funcs {
		if err = f(); err != nil {
			return nil, err
		}
	}

	// resultinfo
	result.DBInfo, err = utils.NewDBInfo(result.DB)

	return
}

/*
Returns test context from environment variables
*/
func NewTest() (*Context, error) {
	c, e := New(os.Getenv("DB_DSN"), os.Getenv("ERGOQ_DSN"), os.Getenv("CACHER_DSN"))
	if e != nil {
		return c, e
	}
	return c, nil
}

// returns context by request
func Get(r *http.Request) (cres *Context, err error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var (
		result interface{}
		ok     bool
	)

	if result, ok = data[r]; !ok {
		err = ErrContextNotFound
		return
	}

	cres = result.(*Context)
	return
}

func Clear(r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(data, r)
}

type Context struct {
	// cache connection
	Cache gocacher.Cacher

	// database connection
	DB *sqlx.DB

	// transaction
	Tx *sqlx.Tx

	// dbinfo
	DBInfo *utils.DBInfo

	// message queue connection
	Queue ergoq.MessageQueuer

	// quit channel
	Quit chan struct{}

	// store request
	Request *http.Request

	// patrol router
	Router *mux.Router

	// request status
	Status int

	// variables
	Vars map[interface{}]interface{}
}

func (c *Context) Copy() (ccopy *Context) {
	ccopy = &Context{
		Cache:   c.Cache,
		DB:      c.DB,
		DBInfo:  c.DBInfo,
		Queue:   c.Queue,
		Quit:    c.Quit,
		Request: c.Request,
		Router:  c.Router,
		Status:  c.Status,
		Vars:    map[interface{}]interface{}{},
	}

	// Copy Vars
	for k, v := range c.Vars {
		ccopy.Vars[k] = v
	}

	return
}

func (c *Context) WithRequest(r *http.Request) (ccopy *Context) {
	mutex.Lock()
	defer mutex.Unlock()
	ccopy = c.Copy()
	ccopy.Request = r
	data[r] = ccopy
	return ccopy
}

func (c *Context) dialDB(dsn string) (err error) {
	var connstr string
	if connstr, err = pq.ParseURL(dsn); err != nil {
		return fmt.Errorf("patrol: database parse dsn error %s.", err)
	}
	if c.DB, err = sqlx.Connect("postgres", connstr); err != nil {
		return fmt.Errorf("patrol: database connect error %s.", err)
	}
	return nil
}

func (c *Context) dialMQ(dsn string) (err error) {
	if c.Queue, err = ergoq.Open(settings.SETTINGS_MESSAGE_QUEUE_DSN); err != nil {
		return fmt.Errorf("patrol: queue open error %s.", err)
	}
	return nil
}

func (c *Context) dialCache(dsn string) (err error) {
	if c.Cache, err = gocacher.Open(dsn); err != nil {
		return fmt.Errorf("patrol: cache open error %s.", err)
	}
	return nil
}

func (c *Context) Set(key, value interface{}) {
	c.Vars[key] = value
}

func (c *Context) Get(key interface{}) (value interface{}) {
	value, _ = c.Vars[key]
	return
}

func (c *Context) GetOk(key interface{}) (interface{}, bool) {
	value, ok := c.Vars[key]
	return value, ok
}

/*
Transaction related methods
*/
func (c *Context) Begin() (err error) {
	if c.Tx, err = c.DB.Beginx(); err != nil {
		return
	}
	return
}

func (c *Context) Commit() (err error) {
	if c.Tx == nil {
		return ErrTransactionNotStarted
	}
	return c.Tx.Commit()
}

func (c *Context) Rollback() (err error) {
	if c.Tx == nil {
		return ErrTransactionNotStarted
	}
	return c.Tx.Rollback()
}

// bind request data to structure
func (c *Context) Bind(target interface{}) (err error) {
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(target)
}
