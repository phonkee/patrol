package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

var (
	chars                = []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ23456789")
	ErrCannotParseAuthID = errors.New("Cannot parse auth ID")
	ErrInvalidInput      = errors.New("invalid input")
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func DBArrayPlaceholder(count int) string {
	if count == 0 {
		return ""
	}
	return strings.Repeat(",?", count)[1:]
}

func StringPadLeft(s, char string, l int) string {
	ls := len(s)
	if ls >= l {
		return s
	}
	return strings.Repeat(char, l-ls) + s
}

func StringPadRight(s, char string, l int) string {
	ls := len(s)
	if ls >= l {
		return s
	}
	return s + strings.Repeat(char, l-ls)
}

/*
	Returns index where is string located
	third parameter is not required (default is 0 - from start)
*/
func StringIndex(list []string, s string, start ...int) (result int) {
	st := 0
	ll := len(list)
	if len(start) > 0 {
		if start[0] < 0 {
			st = start[0] % ll
			if st < 0 {
				st = st + ll
			}
		} else if start[0] < ll {
			st = start[0]
		}
	}
	result = -1
	for i := st; i < ll; i++ {
		if list[i] == s {
			result = i
			break
		}
	}
	return
}

/*
	splits migration identifier into id and pluginId
	so e.g.
		auth:initial-migration will be splitted into "auth", "initial-migration"
		in case of pluginId is not given in string core id will be used (patrol internals)
		you can definitely depend on core migrations
*/
func SplitIdentifier(identifier, defaultPlugin string) (string, string, error) {
	parts := strings.SplitN(strings.TrimSpace(identifier), ":", 2)

	pid := defaultPlugin
	mid := ""

	if len(parts) == 2 {
		if parts[0] != "" {
			pid = parts[0]
		}

		if parts[1] == "" {
			return "", "", ErrInvalidInput
		}

		mid = parts[1]
	} else {
		if parts[0] == "" {
			return "", "", ErrInvalidInput
		}
		mid = parts[0]
	}

	return pid, mid, nil
}

/*
	Generates random string
*/
func RandomString(n int, characters ...string) string {
	c := chars
	if len(characters) > 0 {
		c = []rune(characters[0])
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = c[rand.Intn(len(c))]
	}
	return string(b)
}

func ParseAuth(auth string) (id string, values map[string]string, err error) {
	values = make(map[string]string)
	index := strings.Index(auth, " ")
	if index == -1 {
		err = ErrCannotParseAuthID
	} else {
		id = auth[:index]
		remainder := auth[index+1:]
		for _, part := range strings.Split(remainder, ",") {
			part = strings.TrimSpace(part)
			kv := strings.SplitN(part, "=", 2)
			values[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	return
}
