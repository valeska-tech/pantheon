# Pantheon

A wrapper framework around Nats for Go. Allows the use of abstracted handlers and context and adds avro schema validation.

## Usage

``` go

    package main

    import "github.com/valeska-tech/pantheon"

    func main() {
        app := pantheon.NewApp()
        defer app.Run()

        app.Handler("topic.name", handler)
        app.Handler("test.topic", handlerWithSchema).WithAvro("schema.json")
    }

    func handler(ctx *pantheon.Context) {
        data := &DataStuct{}
        ctx.MustBind(data)

        // Do stuff
    }

    func handlerWithSchema(ctx *pantheon.Context) {
        ctx.App.Log.Info("handler")
        ctx.App.Produce("topic.name", nil)
    }
```

## Why?

After moving to Nats from Kafka I found myself building this structure in all my services as a wrapper around Nats.
