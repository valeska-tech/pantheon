package test

import (
	"os"
	"testing"

	gnatsd "github.com/nats-io/gnatsd/test"
	"github.com/stretchr/testify/assert"
	"github.com/valeska-tech/pantheon"
)

func MsgHandler(ctx *pantheon.Context) {
	done := ctx.MustGet("done").(chan bool)
	done <- true
}

func UnsubHandler(ctx *pantheon.Context) {
	go ctx.Unsubscribe()
	done := ctx.MustGet("done").(chan bool)
	done <- true
}

func Test_itCanHandleAndProduceAMessage(t *testing.T) {
	s := gnatsd.RunDefaultServer()
	defer s.Shutdown()
	os.Setenv("NATS_URLS", "nats://localhost:4222")
	done := make(chan bool)

	app := pantheon.NewApp()
	app.Handler("test.event", MsgHandler)
	app.With("done", done)
	app.RegisterHandlers()

	app.Produce("test.event", nil)

	<-done
}

func Test_itCanUnsubscribeFromTopic(t *testing.T) {
	s := gnatsd.RunDefaultServer()
	defer s.Shutdown()
	os.Setenv("NATS_URLS", "nats://localhost:4222")
	done := make(chan bool)

	app := pantheon.NewApp()
	app.Handler("test.unsub", UnsubHandler)
	app.With("done", done)
	app.RegisterHandlers()

	app.Produce("test.unsub", nil)

	<-done

	_, ok := app.Handlers["test_unsub"]

	assert.Equal(t, false, ok)
}
