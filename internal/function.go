package internal

import (
	"errors"
	"math"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// functionType is a dynamically created (typically by the parser) dgo.FunctionType
	functionType struct {
		arguments dgo.TupleType
		returns   dgo.TupleType
	}

	// exactFunctionTuple is a dgo.Tuple type that is backed either by the NumIn() and In() methods or
	// the NumOut() and Out() methods of a reflect.Type of kind reflect.Func so that the tuple either
	// represents the arguments or the return values of a that goFunc
	exactFunctionTuple struct {
		count    func() int
		element  func(index int) reflect.Type
		variadic bool
	}

	// goFunc represents a go func
	goFunc reflect.Value

	// goFuncType represents a go func type
	goFuncType struct {
		reflect.Type
	}
)

func (t *exactFunctionTuple) Assignable(other interface{}) bool {
	return Assignable(nil, t, other)
}

func (t *exactFunctionTuple) DeepAssignable(guard dgo.RecursionGuard, other dgo.Value) bool {
	return tupleAssignable(guard, t, other)
}

func (t *exactFunctionTuple) Element(index int) dgo.Value {
	rt := t.element(index)
	if t.variadic {
		n := t.count() - 1
		if n == index {
			rt = rt.Elem()
		}
	}
	return TypeFromReflected(rt)
}

func (t *exactFunctionTuple) ElementType() dgo.Value {
	return tupleElementType(t)
}

func (t *exactFunctionTuple) ElementTypes() dgo.Array {
	return &array{slice: t.typeSlice(), frozen: true}
}

func (t *exactFunctionTuple) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *exactFunctionTuple) deepEqual(seen []dgo.Value, other deepEqual) bool {
	return tupleEquals(nil, t, other)
}

func (t *exactFunctionTuple) HashCode() int {
	return tupleHashCode(t, nil)
}

func (t *exactFunctionTuple) deepHashCode(seen []dgo.Value) int {
	return tupleHashCode(t, seen)
}

func (t *exactFunctionTuple) Len() int {
	return t.count()
}

func (t *exactFunctionTuple) ReflectType() reflect.Type {
	return reflect.SliceOf(t.ElementType().ReflectType())
}

func (t *exactFunctionTuple) Max() int {
	return tupleMax(t)
}

func (t *exactFunctionTuple) Min() int {
	return tupleMin(t)
}

func (t *exactFunctionTuple) String() string {
	return TypeString(t)
}

func (t *exactFunctionTuple) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTuple
}

func (t *exactFunctionTuple) Variadic() bool {
	return t.variadic
}

func (t *exactFunctionTuple) Unbounded() bool {
	return t.variadic && t.count() == 1
}

func (t *exactFunctionTuple) typeSlice() []dgo.Value {
	na := t.count()
	as := make([]dgo.Value, na)
	vdic := t.variadic
	for i := 0; i < na; i++ {
		rt := t.element(i)
		if vdic && i == na-1 {
			rt = rt.Elem()
		}
		as[i] = TypeFromReflected(rt)
	}
	return as
}

// DefaultFunctionType is a function that without any constraints on arguments or return value
var DefaultFunctionType = &functionType{arguments: DefaultTupleType, returns: DefaultTupleType}

// FunctionType returns a new dgo.FunctionType with the given argument and return value
// types.
func FunctionType(args dgo.TupleType, returns dgo.TupleType) dgo.FunctionType {
	if returns.Variadic() && !DefaultTupleType.Equals(returns) {
		panic(errors.New(`tuple describing return values cannot be variadic`))
	}
	if args == DefaultTupleType && returns == DefaultTupleType {
		return DefaultFunctionType
	}
	return &functionType{arguments: args, returns: returns}
}

func (t *functionType) Assignable(other interface{}) bool {
	return Assignable(nil, t, other)
}

func (t *functionType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Value) bool {
	return functionTypeAssignable(guard, t, other)
}

func functionTypeAssignable(guard dgo.RecursionGuard, t dgo.FunctionType, other interface{}) bool {
	if ot, ok := other.(dgo.FunctionType); ok {
		return t.Variadic() == ot.Variadic() &&
			tupleAssignableTuple(guard, t.In(), ot.In()) &&
			tupleAssignableTuple(guard, t.Out(), ot.Out())
	}
	if b, ok := toReflected(other); ok {
		return reflect.Func == b.Kind() && functionTypeAssignable(guard, t, &goFuncType{b.Type()})
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *functionType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *functionType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*functionType); ok {
		return equals(seen, t.arguments, ot.arguments) && equals(seen, t.returns, ot.returns)
	}
	return false
}

func (t *functionType) HashCode() int {
	return t.deepHashCode(nil)
}

func (t *functionType) deepHashCode(seen []dgo.Value) int {
	h := int(dgo.TiFunction)
	h = h*31 + deepHashCode(seen, t.arguments)
	h = h*31 + deepHashCode(seen, t.returns)
	return h
}

func (t *functionType) In() dgo.TupleType {
	return t.arguments
}

func (t *functionType) Out() dgo.TupleType {
	return t.returns
}

func (t *functionType) ReflectType() reflect.Type {
	// There is currently no way to build a goFunc type dynamically
	panic(errors.New(`unable to build reflect.Type of go func dynamically`))
}

func (t *functionType) String() string {
	return TypeString(t)
}

func (t *functionType) Variadic() bool {
	return t.arguments.Variadic()
}

func (t *functionType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFunction
}

func (f *goFuncType) Assignable(other interface{}) bool {
	return Assignable(nil, f, other)
}

func (f *goFuncType) DeepAssignable(guard dgo.RecursionGuard, other interface{}) bool {
	return functionTypeAssignable(guard, f, other)
}

func (f *goFuncType) Equals(other interface{}) bool {
	if ov, ok := other.(*goFuncType); ok {
		return f.Type == ov.Type
	}
	return false
}

func (f *goFuncType) HashCode() int {
	return int(reflect.ValueOf(f.Type).Pointer())
}

func (f *goFuncType) In() dgo.TupleType {
	rt := f.Type
	return &exactFunctionTuple{count: rt.NumIn, element: rt.In, variadic: rt.IsVariadic()}
}

func (f *goFuncType) ReflectType() reflect.Type {
	return f.Type
}

func (f *goFuncType) Out() dgo.TupleType {
	rt := f.Type
	return &exactFunctionTuple{count: rt.NumOut, element: rt.Out, variadic: false}
}

func (f *goFuncType) String() string {
	return TypeString(f)
}

func (f *goFuncType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFunction
}

func (f *goFuncType) Variadic() bool {
	return f.Type.IsVariadic()
}

func (f *goFunc) Assignable(other interface{}) bool {
	return Assignable(nil, f, other)
}

func (f *goFunc) DeepAssignable(guard dgo.RecursionGuard, other interface{}) bool {
	return functionTypeAssignable(guard, f, other)
}

func (f *goFunc) Call(args dgo.Array) []dgo.Value {
	convertReturn := func(rr []reflect.Value) []dgo.Value {
		vr := make([]dgo.Value, len(rr))
		for i := range rr {
			re := rr[i]
			v := ValueFromReflected(re)
			if v == Nil {
				_, ok := re.Interface().(dgo.Value)
				if !ok {
					v = nil
				}
			}
			vr[i] = v
		}
		return vr
	}

	mx := args.Len()
	m := (*reflect.Value)(f)
	t := m.Type()
	if t.IsVariadic() {
		nv := t.NumIn() - 1 // number of non variadic
		if mx < nv {
			panic(illegalArgumentCount(t.Name(), nv, math.MaxInt64, mx))
		}
		rr := make([]reflect.Value, nv+1)
		for i := 0; i < nv; i++ {
			rr[i] = reflect.New(t.In(i)).Elem()
			ReflectTo(args.Get(i), rr[i])
		}

		// Create the variadic slice
		vt := t.In(nv)
		vz := mx - nv
		vs := reflect.MakeSlice(vt, vz, vz)
		rr[nv] = vs

		for i := 0; i < vz; i++ {
			ReflectTo(args.Get(i+nv), vs.Index(i))
		}
		return convertReturn(m.CallSlice(rr))
	}

	if mx != t.NumIn() {
		panic(illegalArgumentCount(t.Name(), t.NumIn(), t.NumIn(), mx))
	}

	rr := make([]reflect.Value, mx)
	for i := 0; i < mx; i++ {
		rr[i] = reflect.New(t.In(i)).Elem()
		ReflectTo(args.Get(i), rr[i])
	}
	return convertReturn(m.Call(rr))
}

func (f *goFunc) Equals(other interface{}) bool {
	if ov, ok := other.(*goFunc); ok {
		return (*reflect.Value)(f).Pointer() == (*reflect.Value)(ov).Pointer()
	}
	if b, ok := toReflected(other); ok {
		return reflect.Func == b.Kind() && (*reflect.Value)(f).Pointer() == b.Pointer()
	}
	return false
}

func (f *goFunc) GoFunc() interface{} {
	return (*reflect.Value)(f).Interface()
}

func (f *goFunc) HashCode() int {
	return int((*reflect.Value)(f).Pointer())
}

func (f *goFunc) In() dgo.TupleType {
	rt := (*reflect.Value)(f).Type()
	return &exactFunctionTuple{count: rt.NumIn, element: rt.In, variadic: rt.IsVariadic()}
}

func (f *goFunc) ReflectType() reflect.Type {
	return (*reflect.Value)(f).Type()
}

func (f *goFunc) Out() dgo.TupleType {
	rt := (*reflect.Value)(f).Type()
	return &exactFunctionTuple{count: rt.NumOut, element: rt.Out, variadic: false}
}

func (f *goFunc) String() string {
	return (*reflect.Value)(f).String()
}

func (f *goFunc) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFunction
}

func (f *goFunc) Variadic() bool {
	return f.ReflectType().IsVariadic()
}
