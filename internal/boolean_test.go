package internal_test

import (
	"reflect"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestBooleanType(t *testing.T) {
	v := vf.True
	require.Assignable(t, v, v)
	require.Assignable(t, typ.Boolean, v)
	require.Assignable(t, typ.True, v)
	require.Assignable(t, v, true)
	require.Assignable(t, typ.Boolean, true)
	require.Assignable(t, typ.True, true)
	require.NotAssignable(t, typ.False, v)
	require.NotAssignable(t, typ.False, true)
	require.NotAssignable(t, v, typ.Boolean)
	require.Equal(t, v, v)
	require.True(t, v.IsInstance(true))
	require.False(t, v.IsInstance(false))
	require.Equal(t, `true`, v.String())

	v = vf.False
	require.Assignable(t, typ.Boolean, v)
	require.Assignable(t, typ.False, v)
	require.NotAssignable(t, typ.True, v)
	require.Assignable(t, typ.Boolean, v)
	require.NotAssignable(t, v, typ.Boolean)
	require.True(t, v.IsInstance(false))
	require.False(t, v.IsInstance(true))
	require.Equal(t, `false`, v.String())

	require.Equal(t, typ.Boolean.HashCode(), typ.Boolean.HashCode())
	require.NotEqual(t, 0, typ.Boolean.HashCode())
	require.NotEqual(t, 0, typ.True.HashCode())
	require.NotEqual(t, 0, typ.False.HashCode())
	require.NotEqual(t, typ.Boolean.HashCode(), typ.True.HashCode())
	require.NotEqual(t, typ.Boolean.HashCode(), typ.False.HashCode())
	require.NotEqual(t, typ.True.HashCode(), typ.False.HashCode())
	require.Equal(t, `bool`, typ.Boolean.String())

	require.Equal(t, reflect.TypeOf(false), typ.True.ReflectType())
}

func TestNew_bool(t *testing.T) {
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.String(`y`)))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.String(`Yes`)))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.String(`TRUE`)))
	require.Equal(t, vf.False, vf.New(typ.Boolean, vf.String(`N`)))
	require.Equal(t, vf.False, vf.New(typ.Boolean, vf.String(`no`)))
	require.Equal(t, vf.False, vf.New(typ.Boolean, vf.String(`False`)))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.Float(1)))
	require.Equal(t, vf.False, vf.New(typ.Boolean, vf.Float(0)))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.Float(1)))
	require.Equal(t, vf.False, vf.New(typ.Boolean, vf.Integer(0)))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.Integer(1)))
	require.Equal(t, vf.False, vf.New(typ.Boolean, vf.False))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.True))
	require.Equal(t, vf.True, vf.New(typ.Boolean, vf.Arguments(vf.True)))
	require.Panic(t, func() { vf.New(typ.Boolean, vf.String(`unhappy`)) }, `unable to create a bool from unhappy`)
	require.Panic(t, func() { vf.New(typ.Boolean, vf.Arguments(vf.True, vf.True)) }, `illegal number of arguments`)
}

func TestBoolean(t *testing.T) {
	require.Equal(t, vf.True, vf.Boolean(true))
	require.Equal(t, vf.False, vf.Boolean(false))
}

func TestBoolean_Equals(t *testing.T) {
	require.True(t, vf.True.Equals(vf.True))
	require.True(t, vf.True.Equals(true))
	require.False(t, vf.True.Equals(vf.False))
	require.False(t, vf.True.Equals(false))
	require.True(t, vf.False.Equals(vf.False))
	require.True(t, vf.False.Equals(false))
	require.False(t, vf.False.Equals(vf.True))
	require.False(t, vf.False.Equals(true))
	require.True(t, vf.True.GoBool())
	require.False(t, vf.False.GoBool())
}

func TestBoolean_HashCode(t *testing.T) {
	require.NotEqual(t, 0, vf.True.HashCode())
	require.NotEqual(t, 0, vf.False.HashCode())
	require.Equal(t, vf.True.HashCode(), vf.True.HashCode())
	require.NotEqual(t, vf.True.HashCode(), vf.False.HashCode())
}

func TestBoolean_CompareTo(t *testing.T) {
	c, ok := vf.True.CompareTo(vf.True)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.True.CompareTo(vf.False)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.False.CompareTo(vf.True)
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.True.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.False.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = vf.True.CompareTo(vf.Integer(1))
	require.False(t, ok)
}

func TestBoolean_ReflectTo(t *testing.T) {
	var b bool
	vf.True.ReflectTo(reflect.ValueOf(&b).Elem())
	require.True(t, b)

	var bp *bool
	vf.True.ReflectTo(reflect.ValueOf(&bp).Elem())
	require.True(t, *bp)

	var mi interface{}
	mip := &mi
	vf.True.ReflectTo(reflect.ValueOf(mip).Elem())
	bc, ok := mi.(bool)
	require.True(t, ok)
	require.True(t, bc)
}

func TestBoolean_String(t *testing.T) {
	require.Equal(t, `true`, vf.True.String())
	require.Equal(t, `false`, vf.False.String())
}
