package validator

import "errors"

/*
Constructor to return new result
*/
func NewResult() *Result {
	return &Result{
		Fields:  map[string][]string{},
		Unbound: []string{},
	}
}

/*
Validator result which is json marshallable and supports 2 kinds of errors
fields - bound to fields
unbound - other errors (such as usernot found)
*/
type Result struct {
	Fields  map[string][]string `json:"fields,omitempty"`
	Unbound []string            `json:"unbound,omitempty"`
}

/*
Adds field error, if duplicate error found, will not be added
*/
func (r *Result) AddFieldError(field string, err error) *Result {
	errstr := err.Error()

	if _, ok := r.Fields[field]; !ok {
		r.Fields[field] = []string{errstr}
	}
	for _, v := range r.Fields[field] {
		if v == errstr {
			return r
		}
	}
	r.Fields[field] = append(r.Fields[field], errstr)

	return r
}

func (r *Result) AddUnboundError(err error) *Result {
	r.Unbound = append(r.Unbound, err.Error())
	return r
}

/*
Returns errors for given field, if not found blank slice is returned
*/
func (r *Result) GetFieldErrors(field string) []string {
	if _, ok := r.Fields[field]; !ok {
		return []string{}
	}
	return r.Fields[field]
}

/*
Returns whether given field has any errors
*/
func (r *Result) HasFieldErrors(field string) bool {
	return len(r.GetFieldErrors(field)) > 0
}

/*
Appends errors from other result
*/
func (r *Result) Append(other *Result) *Result {
	for field, errs := range other.Fields {
		for _, err := range errs {
			r.AddFieldError(field, errors.New(err))
		}
	}
	return r
}

/*
Excludes fields from result
*/
func (r *Result) Exclude(fields ...string) (nr *Result) {
	nr = NewResult()
	nr.Append(r)
	for _, field := range fields {
		delete(nr.Fields, field)
	}
	return nr
}

/*
Excludes fields from result
*/
func (r *Result) Allow(fields ...string) (nr *Result) {
	nr = NewResult()
	for field, errs := range r.Fields {
		for _, fieldi := range fields {
			if field == fieldi {
				for _, err := range errs {
					nr.AddFieldError(field, errors.New(err))
				}
			}
		}
	}
	return
}

/*
Returns if is valid result
*/
func (r *Result) IsValid() bool {
	return len(r.Unbound)+len(r.Fields) == 0
}

// Try to read postgres error
// func (r *Result) AddPostgresError(err error, context *context.Context) bool {
// 	if err, ok := err.(*pq.Error); ok {

// 		// Postgres errors
// 		// http://www.postgresql.org/docs/9.4/static/errcodes-appendix.html

// 		// @TODO: add all these
// 		// 23000	integrity_constraint_violation
// 		// 23001	restrict_violation
// 		// 23P01	exclusion_violation

// 		switch err.Code.Name() {
// 		case "unique_violation", "check_violation", "foreign_key_violation", "not_null_violation":
// 			v.AddFieldError(context.DBInfo.Constraints[err.Constraint], err.Code.Name())
// 		default:
// 		}
// 	}

// 	return true
// }
