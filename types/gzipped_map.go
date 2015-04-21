package types

import (
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"

	"github.com/phonkee/patrol/settings"
)

/*
GzippedMap
store map[string]interface{} to database as gzipped JSON
*/

var (
	ErrGzippedMapBadValue = errors.New("gzipped map: bad value")
)

type GzippedMap map[string]interface{}

// Scan implements the Scanner interface.
func (g GzippedMap) Scan(value interface{}) (err error) {
	var (
		body []byte
		ok   bool
	)
	if body, ok = value.([]byte); !ok {
		err = ErrGzippedMapBadValue
		return
	}

	buffer := bytes.NewReader(body)
	var gr *gzip.Reader
	if gr, err = gzip.NewReader(buffer); err != nil {
		return
	}
	gr.Close()

	result := new(bytes.Buffer)
	io.Copy(result, gr)

	return json.Unmarshal(result.Bytes(), &g)
}

// Value implements the driver Valuer interface.
func (g GzippedMap) Value() (val driver.Value, err error) {
	var body []byte

	if body, err = json.Marshal(g); err != nil {
		return
	}

	buffer := new(bytes.Buffer)
	var writer *gzip.Writer
	if writer, err = gzip.NewWriterLevel(buffer, settings.GZIPPED_MAP_COMPRESSION_LEVEL); err != nil {
		return
	}
	writer.Write(body)
	writer.Close()
	return buffer.Bytes(), nil
}
