package test

import (
	"testing"

	avro "github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"github.com/valeska-tech/pantheon"
)

type (
	TestData struct {
		Foo string `json:"foo"`
	}
)

func createContext() *pantheon.Context {
	return &pantheon.Context{
		App: &pantheon.Application{
			Handlers: make(map[string]*pantheon.HandlerWrapper),
			Params:   make(map[string]interface{}),
		},
	}
}

func Test_itCanBeCopied(t *testing.T) {
	ctx := createContext()
	copy := ctx.Copy()

	assert.Equal(t, ctx.App, copy.App)
}

func Test_itCanBindVariables(t *testing.T) {
	ctx := createContext()
	ctx.App.With("test", false)

	v := ctx.MustGet("test").(bool)
	_, err := ctx.Get("foo")

	assert.Equal(t, false, v)
	assert.NotNil(t, err)
}

func Test_itCanDeserializeJSON(t *testing.T) {
	ctx := createContext()
	ctx.Data = `{"foo": "bar"}`

	d := &TestData{}
	ctx.MustBind(d)

	assert.Equal(t, "bar", d.Foo)
}

func Test_itCanBindAndDeserialiseFromAvro(t *testing.T) {
	ctx := createContext()

	codec, _ := avro.NewCodec(`{"type": "record", "name": "foo", "fields": [{"name": "foo", "type": "string"}]}`)
	decoded, _, _ := codec.NativeFromTextual([]byte(`{"foo": "bar"}`))
	ctx.Decoded = decoded

	d := &TestData{}
	ctx.MustBind(d)

	assert.Equal(t, "bar", d.Foo)
}
