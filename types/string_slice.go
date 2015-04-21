package types

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// construct a regexp to extract values:
var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

type StringSlice []string

func (s *StringSlice) Add(str string) {
	*s = append(*s, str)
}

func (s *StringSlice) AddUnique(str string) {
	for _, item := range *s {
		if item == str {
			return
		}
	}
	s.Add(str)
	return
}

func (s *StringSlice) Compact() {
	found := map[string]int{}

	for _, item := range *s {
		if _, ok := found[item]; ok {
			found[item]++
		} else {
			found[item] = 1
		}
	}

	for k, val := range found {
		if val == 1 {
			continue
		} else {
			s.Remove(k)
		}
	}
}

func (s *StringSlice) Has(str string) bool {
	return s.Index(str) != -1
}

func (s *StringSlice) Index(str string) (result int) {
	result = -1
	for i, item := range *s {
		if item == str {
			result = i
			break
		}
	}
	return
}

func (s *StringSlice) Remove(str string) {
	var index int
	for {
		index = s.Index(str)
		if index == -1 {
			break
		}

		n := StringSlice{}
		for i := 0; i < len(*s); i++ {
			if i == index {
				continue
			}
			n = append(n, (*s)[i])
		}
		*s = n
	}
	return
}

func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}

	asString := string(asBytes)
	parsed := parseArray(asString)
	(*s) = StringSlice(parsed)

	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	parts := []string{}
	for _, val := range s {
		parts = append(parts, quote(val))
	}
	return []byte("{" + strings.Join(parts, ",") + "}"), nil
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseArray(array string) []string {
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}

func quote(s interface{}) string {
	var str string
	switch v := s.(type) {
	case sql.NullString:
		if !v.Valid {
			return "NULL"
		}
		str = v.String
	case string:
		str = v
	default:
		panic("not a string or sql.NullString")
	}

	str = strings.Replace(str, "\\", "\\\\", -1)
	return `"` + strings.Replace(str, "\"", "\\\"", -1) + `"`
}
