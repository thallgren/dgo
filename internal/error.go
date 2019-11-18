package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

type (
	errw struct {
		error
	}

	errType int
)

// DefaultErrorType is the unconstrained Error type
const DefaultErrorType = errType(0)

var reflectErrorType = reflect.TypeOf((*error)(nil)).Elem()

func (t errType) Assignable(other interface{}) bool {
	_, ok := other.(error)
	return ok || t == other || CheckAssignableTo(nil, other, t)
}

func (t errType) Equals(other interface{}) bool {
	return t == DefaultErrorType
}

func (t errType) HashCode() int {
	return int(t.TypeIdentifier())
}

func (t errType) ReflectType() reflect.Type {
	return reflectErrorType
}

func (t errType) String() string {
	return TypeString(t)
}

func (t errType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiError
}

func (e *errw) Assignable(other interface{}) bool {
	return e.Equals(other) || CheckAssignableTo(nil, other, e)
}

func (e *errw) Equals(other interface{}) bool {
	if oe, ok := other.(*errw); ok {
		return e.error.Error() == oe.error.Error()
	}
	if oe, ok := other.(error); ok {
		return e.error.Error() == oe.Error()
	}
	return false
}

func (e *errw) HashCode() int {
	return util.StringHash(e.error.Error())
}

func (e *errw) Error() string {
	return e.error.Error()
}

func (e *errw) ReflectTo(value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		value.Set(reflect.ValueOf(&e.error))
	} else {
		value.Set(reflect.ValueOf(e.error))
	}
}

func (e *errw) ReflectType() reflect.Type {
	return reflectErrorType
}

func (e *errw) String() string {
	return e.error.Error()
}

func (e *errw) Unwrap() error {
	if u, ok := e.error.(interface {
		Unwrap() error
	}); ok {
		return u.Unwrap()
	}
	return nil
}

func (e *errw) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiErrorExact
}
