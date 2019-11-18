package internal

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// regexpType represents an regexp type without constraints
	regexpType int

	regexpVal regexp.Regexp
)

// DefaultRegexpType is the unconstrained Regexp type
const DefaultRegexpType = regexpType(0)

var reflectRegexpType = reflect.TypeOf(&regexp.Regexp{})

func (t regexpType) Assignable(ot interface{}) bool {
	switch ot.(type) {
	case regexpType, *regexpVal, *regexp.Regexp:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t regexpType) Equals(v interface{}) bool {
	return t == v
}

func (t regexpType) HashCode() int {
	return int(dgo.TiRegexp)
}

func (t regexpType) ReflectType() reflect.Type {
	return reflectRegexpType
}

func (t regexpType) String() string {
	return TypeString(t)
}

func (t regexpType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexp
}

// Regexp returns the given regexp as a dgo.Regexp
func Regexp(rx *regexp.Regexp) dgo.Regexp {
	return (*regexpVal)(rx)
}

func (v *regexpVal) Assignable(other interface{}) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *regexpVal) GoRegexp() *regexp.Regexp {
	return (*regexp.Regexp)(v)
}

func (v *regexpVal) Equals(other interface{}) bool {
	if ot, ok := other.(*regexpVal); ok {
		return (*regexp.Regexp)(v).String() == (*regexp.Regexp)(ot).String()
	}
	if ot, ok := other.(*regexp.Regexp); ok {
		return (*regexp.Regexp)(v).String() == (ot).String()
	}
	return false
}

func (v *regexpVal) Generic() dgo.Value {
	return DefaultRegexpType
}

func (v *regexpVal) HashCode() int {
	return util.StringHash((*regexp.Regexp)(v).String())
}

func (v *regexpVal) IsInstance(rx *regexp.Regexp) bool {
	return (*regexp.Regexp)(v).String() == rx.String()
}

func (v *regexpVal) ReflectTo(value reflect.Value) {
	rv := reflect.ValueOf((*regexp.Regexp)(v))
	k := value.Kind()
	if !(k == reflect.Ptr || k == reflect.Interface) {
		rv = rv.Elem()
	}
	value.Set(rv)
}

func (v *regexpVal) String() string {
	return (*regexp.Regexp)(v).String()
}

func (v *regexpVal) ReflectType() reflect.Type {
	return reflectRegexpType
}

func (v *regexpVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexpExact
}

// RegexpSlashQuote converts the given string into a slash delimited string with internal slashes escaped
// and writes it on the given builder.
func RegexpSlashQuote(sb *strings.Builder, str string) {
	util.WriteByte(sb, '/')
	for _, c := range str {
		switch c {
		case '\t':
			util.WriteString(sb, `\t`)
		case '\n':
			util.WriteString(sb, `\n`)
		case '\r':
			util.WriteString(sb, `\r`)
		case '/':
			util.WriteString(sb, `\/`)
		case '\\':
			util.WriteString(sb, `\\`)
		default:
			if c < 0x20 {
				util.Fprintf(sb, `\u{%X}`, c)
			} else {
				util.WriteRune(sb, c)
			}
		}
	}
	util.WriteByte(sb, '/')
}
