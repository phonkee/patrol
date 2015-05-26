package apitest

import (
	"encoding/json"
	"fmt"
)

/*
Pretty prints json

if data is string - PrettyPrint first unmarshals this to map[string]interface{}
and then pretty prints it. This feature is e.g. for http responses

*/
func PrettyPrint(format string, data interface{}) (err error) {
	format = format + "\n"

	switch data.(type) {
	case string:
		m := map[string]interface{}{}
		if err = json.Unmarshal([]byte(data.(string)), &m); err != nil {
			return err
		}
		data = m
	case []byte:
		m := map[string]interface{}{}
		if err = json.Unmarshal(data.([]byte), &m); err != nil {
			return err
		}
		data = m
	}

	var body []byte
	if body, err = json.MarshalIndent(data, "", "    "); err != nil {
		return err
	}
	fmt.Printf(format, string(body))
	return nil
}
