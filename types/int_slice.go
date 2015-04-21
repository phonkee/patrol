package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []int

// adds num value to IntSlice
func (i *IntSlice) Add(num int) *IntSlice {
	*i = append(*i, num)
	return i
}

// adds num value to IntSlice if not exists
func (i *IntSlice) AddUnique(num int) (result *IntSlice) {
	result = i
	if i.Index(num) == -1 {
		i.Add(num)
	}
	return
}

// Compacts array (removes duplicates)
func (i *IntSlice) Compact() {
	vals := map[int]int{}
	for _, val := range *i {
		if _, ok := vals[val]; ok {
			vals[val]++
		} else {
			vals[val] = 1
		}
	}
	is := make(IntSlice, 0, len(vals))
	for _, v := range *i {
		if vals[v] == 0 {
			continue
		}
		is.Add(v)
		vals[v] = 0
	}
	*i = is
}

// Returns whether intslice has (contains) num
func (i *IntSlice) Has(num int) bool {
	return i.Index(num) != -1
}

// returns index (position) in IntSlice, if not found -1 is returned
func (i *IntSlice) Index(num int) (result int) {
	result = -1
	for index, val := range *i {
		if val == num {
			result = index
			break
		}
	}
	return
}

// Removes num from intSlice
func (i *IntSlice) Remove(num int) {
	var index int
	if index = i.Index(num); index == -1 {
		return
	}
	is := make(IntSlice, 0, len(*i))

	for _, val := range *i {
		if val == num {
			continue
		}
		is = append(is, val)
	}

	(*i) = is
}

// Scan from database data
func (s *IntSlice) Scan(src interface{}) (err error) {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	asString := strings.Trim(string(asBytes), "{}")
	if asString == "" {
		(*s) = IntSlice{}
	}
	splitted := strings.Split(asString, ",")

	ints := make([]int, len(splitted))
	var ival int
	for i, val := range splitted {
		if ival, err = strconv.Atoi(val); err != nil {
			return err
		}
		ints[i] = ival
	}

	(*s) = IntSlice(ints)

	return nil
}

// from intslice to []byte
func (s IntSlice) Value() (driver.Value, error) {
	strs := make([]string, len(s))
	for i, value := range s {
		strs[i] = fmt.Sprintf("%d", value)
	}

	return []byte("{" + strings.Join(strs, ",") + "}"), nil
}
