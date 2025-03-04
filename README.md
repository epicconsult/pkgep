# EPIC Go Standard Library
Collection of Go libraries for Epic-Consulting, Designed with "*Simplicity* and *Flexibility*" in mind.

## Logger

### Example
```go
package main

import (
  "github.com/epicconsult/pkgep/logger"
)

func main() {
	logger.SetLogger(logger.NewLogrus())

	logger.Logger.Info("Start application...")

  
	app := fiber.New()

	logger.Logger.Info("Created api server")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	logger.Logger.Info("Api server running on port 3000")

	app.Listen(":3000")
}
```
### Logrus
Logrus is implemented as a default logger client. you can config logrus to suite your need
* Rotation: log rotation is default to used date based log file ```date```, you can opt to use ```timestamp``` which is more optimization.

```go
logger.SetLogger(logger.NewLogrus(
  logger.WithAppName("epic-app"),
  logger.WithRotationType(logger.Timestamp),
))
```

### Implement your own logger client
You can even provide your own logger client if logrus does not suite your need by implement Standard Logger Interface and pass it to ```SetLogger``` function.

```go
// Standard Logger Interface
type Logger interface {
    Info(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Trace(msg string, args ...interface{})
}
```