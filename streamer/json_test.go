package streamer_test

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/vf"
)

type badWriter int

func (b badWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New(`bang`)
}

func ExampleJSON() {
	s := streamer.New(nil, nil)
	s.Stream(vf.Map(`a`, 1, `b`, []int{1, 2}), streamer.JSON(os.Stdout))
	// Output: {"a":1,"b":[1,2]}
}

func TestJSON_AddRef(t *testing.T) {
	v := vf.Strings(`a`, `b`)
	a := vf.Values(v, v)
	s := streamer.New(nil, nil)
	b := bytes.Buffer{}
	s.Stream(a, streamer.JSON(&b))
	require.Equal(t, `[["a","b"],{"__ref":1}]`, b.String())
}

func TestJSON_primitives(t *testing.T) {
	v := vf.Values(true, nil, 1, 2.1, `string`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(v, streamer.JSON(&b))
	require.Equal(t, `[true,null,1,2.1,"string"]`, b.String())
}

func TestJSON_badWrite(t *testing.T) {
	require.Panic(t, func() { streamer.New(nil, nil).Stream(vf.Integer(3), streamer.JSON(badWriter(0))) }, `bang`)
}

func TestJSON_CanDoBinary(t *testing.T) {
	v := vf.Values(vf.BinaryFromString(`AQID`))
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(v, streamer.JSON(&b))
	require.Equal(t, `[{"__type":"binary","__value":"AQID"}]`, b.String())
}

func TestJSON_CanDoTime(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(vf.Time(ts), streamer.JSON(&b))
	require.Equal(t, `{"__type":"time","__value":"2019-10-06T07:15:00-07:00"}`, b.String())
}

func TestJSON_ComplexKeys(t *testing.T) {
	v := vf.Map(vf.BinaryFromString(`AQID`), `value of binary`, `hey`, `value of hey`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(v, streamer.JSON(&b))
	require.Equal(t, `{"__type":"map","__value":[{"__type":"binary","__value":"AQID"},"value of binary","hey","value of hey"]}`, b.String())
}