package pantheon

import (
	"encoding/json"
	"os"
	"sync"

	nats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

var (
	// SchemaDir holds the directory name where schemas are kept
	SchemaDir string
)

type (
	// Application struct takes care of libraries needed
	// to run the app, at the moment this has a single connection
	// to nats, but in the future this will need to be a pooled connection.
	Application struct {
		Log      *logrus.Logger
		Handlers map[string]*HandlerWrapper
		params   map[string]interface{}
		nats     *nats.Conn
	}
)

// NewApp generates a new application using env vars
func NewApp() *Application {
	// Attempt to connect to the nats cluster
	nc, err := nats.Connect(os.Getenv("NATS_URLS"))

	if err != nil {
		panic(err)
	}

	app := &Application{
		Log:      createLog(),
		Handlers: make(map[string]*HandlerWrapper),
		nats:     nc,
		params:   make(map[string]interface{}),
	}

	// Set the schema and template directories from envs
	SchemaDir = os.Getenv("SCHEMAS")

	return app
}

// Sets up a logrus logger instance and returns it
func createLog() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout

	return log
}

// Produce puts the given message onto the specified topicz
func (a *Application) Produce(subject string, obj interface{}) {
	js, err := json.Marshal(obj)

	if err != nil {
		a.Log.Error(err)
		return
	}

	a.nats.Publish(subject, js)
}

// Handler simply adds a new handler to the handlers map
func (a *Application) Handler(subject string, handler EventHandler) *HandlerWrapper {
	// Create and return the handler wrapper
	w := NewHandler(handler)
	a.Handlers[subject] = w
	return w
}

// With adds a dependency into the params map
func (a *Application) With(key string, param interface{}) {
	a.params[key] = param
}

// Run starts the application daemon, this also starts any consumer
func (a *Application) Run() {
	// Start the handlers
	for subject, handler := range a.Handlers {
		a.Log.Infof("Processing handler for %s subject", subject)
		// Register the handlers channel
		sub, err := a.nats.ChanSubscribe(subject, handler.ch)

		if err != nil {
			a.Log.Error(err)
			continue
		}

		// Set the subscription
		handler.subscription = sub

		// Listen for messages
		go handler.Listen(&Context{App: a})
	}

	// Log a message to inform that this is now running
	a.Log.Info("Running Application...")

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
