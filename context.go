package pantheon

import (
	"errors"
	"io/ioutil"

	avro "github.com/linkedin/goavro"
)

// Context is used to separate nats from the handlers,
// it abstracts the message and gives them access to the application
// without the need for global variables
type Context struct {
	App     *Application
	Data    string
	Decoded interface{}
}

// Copy returns a copy of this context
func (ctx *Context) Copy() *Context {
	return &Context{
		App: ctx.App,
	}
}

// MustGet either finds the parameter on the application or panics
func (ctx *Context) MustGet(key string) interface{} {
	in, err := ctx.Get(key)

	if err != nil {
		panic(err)
	}

	return in
}

// Get fetches a parameter from the application through the context struct
func (ctx *Context) Get(key string) (interface{}, error) {
	in, ok := ctx.App.params[key]

	if !ok {
		return nil, errors.New("Unable for find parameter " + key)
	}

	return in, nil
}

// WithAvro accepts a schema file containing an avro schema
// this method will panic if the file does not exist or if the
// file does not contain valid avro
func (w *HandlerWrapper) WithAvro(schema string) {
	data, err := ioutil.ReadFile(SchemaDir + "/" + schema)

	if err != nil {
		panic(err)
	}

	codec, err := avro.NewCodec(string(data))

	if err != nil {
		panic(err)
	}

	w.Schema = codec
}
