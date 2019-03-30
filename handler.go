package pantheon

import (
	avro "github.com/linkedin/goavro"
	nats "github.com/nats-io/go-nats"
)

type (
	// HandlerWrapper is used to contain an event handler and any
	// options or settings associated with them, you can bind schemas
	// to the wrapper which are evaluated when a message is recieved
	// before they are sent to the handler.
	HandlerWrapper struct {
		Handler      EventHandler
		Schema       *avro.Codec
		subscription *nats.Subscription
		ch           chan *nats.Msg
		unsubscribe  chan bool
	}

	// EventHandler defines the contract for an event handler
	EventHandler func(ctx *Context)
)

// NewHandler returns a new HandlerWrapper struct
func NewHandler(e EventHandler) *HandlerWrapper {
	return &HandlerWrapper{
		Handler:     e,
		ch:          make(chan *nats.Msg),
		unsubscribe: make(chan bool),
	}
}

// Listen allows the handler to listen to the intake and unsubscribe channels
func (w *HandlerWrapper) Listen(ctx *Context) {
handlerloop:
	for {
		select {
		case msg := <-w.ch:
			// forward it onto the handler
			ctx := ctx.Copy()
			ctx.Data = string(msg.Data)
			ctx.wrapper = w

			w.forward(ctx)
		case <-w.unsubscribe:
			// unsubscribe from the
			w.subscription.Unsubscribe()
			break handlerloop
		}
	}
}

// forward is used to accept a message from Nats, this validates the message
// and forwards onto the handler
func (w *HandlerWrapper) forward(ctx *Context) {
	// We always assume the message is valid, unless there is a schema
	// and the schema invalidates the message
	msgValid := true

	// If there is a schema present, validate the payload
	if w.Schema != nil {
		decoded, _, err := w.Schema.NativeFromTextual([]byte(ctx.Data))

		if err != nil {
			msgValid = false
			ctx.App.Log.Errorf("Data sent is invalid against the schema provided: %s", ctx.Data)
		}

		ctx.Decoded = decoded
	}

	// If the payload is ok, forward this to the handler
	if msgValid {
		// Recover from panics in handlers
		defer func() {
			if r := recover(); r != nil {
				// Log it instead
				ctx.App.Log.Error(r)
			}
		}()

		w.Handler(ctx)
	}
}
