/*
	Params is implementation of project/patrol parameters for plugins
	@TODO: architect these
*/
package utils

type PatrolParam interface {
	FromNative(interface{}) []byte
	ToNative([]byte) interface{}
	Validate()
	Description() string
}

type ParamInt struct{}
type ParamBool struct{}
type ParamString struct{}
type ParamFloat struct{}
