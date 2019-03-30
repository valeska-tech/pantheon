package pantheon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	avro "github.com/linkedin/goavro"
	"github.com/mitchellh/mapstructure"
)

// Context is used to separate nats from the handlers,
// it abstracts the message and gives them access to the application
// without the need for global variables
type Context struct {
	App     *Application
	Data    string
	Decoded interface{}
	wrapper *HandlerWrapper
}

// Copy returns a copy of this context
func (ctx *Context) Copy() *Context {
	return &Context{
		App: ctx.App,
	}
}

// Unsubscribe sends a message to the unsubscribe channel of the handler
// attached to the current context
func (ctx *Context) Unsubscribe() {
	ctx.wrapper.unsubscribe <- true
}

// MustBind ensures that the data from an event binds to the
// given interface
func (ctx *Context) MustBind(in interface{}) {
	// Does it have a decoded version?
	if ctx.Decoded != nil {
		// Use mapstructure to decode this biz.
		if err := mapstructure.Decode(ctx.Decoded, in); err != nil {
			panic(err)
		}

		return
	}

	// Otherwise, lets deocde this as json
	if err := json.NewDecoder(strings.NewReader(ctx.Data)).Decode(in); err != nil {
		panic(err)
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
