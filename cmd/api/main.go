package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rulu158/gocoderunner/runner"
)

const port = ":9920"

func main() {
	gin.SetMode(gin.ReleaseMode)

	handleSignals(func() { runner.FreeAllResources() })

	runner.FreeAllResources()

	srv := NewServer()
	srv.Run(port)
}

func handleSignals(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				fn()
				os.Exit(1)
			}
		}
	}()
}
