package internal_test

import (
	"errors"
	"reflect"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestErrorType(t *testing.T) {
	er := errors.New(`some error`)
	v := vf.Value(er)
	require.Same(t, typ.Error, v)
	require.Assignable(t, v, er)
	require.Assignable(t, v, v)
	require.NotAssignable(t, v, typ.String)
	require.Equal(t, v, v)

	require.Equal(t, v.HashCode(), v.HashCode())
	require.NotEqual(t, 0, v.HashCode())

	require.Equal(t, `error`, v.String())

	require.True(t, reflect.TypeOf(er).AssignableTo(v.ReflectType()))
}

type testUnwrapError struct {
	err error
}

func (e *testUnwrapError) Error() string {
	return `unwrap me`
}

func (e *testUnwrapError) Unwrap() error {
	return e.err
}

func TestError(t *testing.T) {
	ev := errors.New(`some error`)
	v := vf.Value(ev)

	ve, ok := v.(error)
	require.True(t, ok)

	require.Equal(t, v, ev)
	require.Equal(t, v, errors.New(`some error`))
	require.Equal(t, v, vf.Value(errors.New(`some error`)))
	require.NotEqual(t, v, `some error`)

	require.NotEqual(t, 0, v.HashCode())
	require.Equal(t, v.HashCode(), vf.Value(errors.New(`some error`)).HashCode())
	require.NotEqual(t, v.HashCode(), vf.Value(errors.New(`some other error`)).HashCode())

	require.Equal(t, `some error`, ve.Error())

	require.Equal(t, ve.Error(), v.String())

	type uvt interface {
		Unwrap() error
	}

	u, ok := v.(uvt)
	require.True(t, ok)
	require.Nil(t, u.Unwrap())

	uv := vf.Value(&testUnwrapError{ev})
	u, ok = uv.(uvt)
	require.True(t, ok)
	require.Equal(t, u.Unwrap(), ev)
}

func TestError_ReflectTo(t *testing.T) {
	var err error
	v := vf.Value(errors.New(`some error`))
	vf.ReflectTo(v, reflect.ValueOf(&err).Elem())
	require.NotNil(t, err)
	require.Equal(t, `some error`, err.Error())

	var bp *error
	vf.ReflectTo(v, reflect.ValueOf(&bp).Elem())
	require.NotNil(t, bp)
	require.NotNil(t, *bp)
	require.Equal(t, `some error`, (*bp).Error())

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(error)
	require.True(t, ok)
	require.Same(t, err, ec)
}
