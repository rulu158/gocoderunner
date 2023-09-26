package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rulu158/gocoderunner/runner"
	"github.com/rulu158/gocoderunner/runner/languages"
)

func (srv *Server) ExecRunner(ctx *gin.Context) {
	var sb strings.Builder
	runner := *runner.NewRunner(languages.Go, &runner.RunnerOptions{
		Stdout:  &sb,
		Timeout: 30 * time.Second,
	})

	stdCh := make(chan int, 1)
	go func() {
		runner.ExecCode()
		stdCh <- 1
	}()

	var (
		timer  *time.Timer
		timerC <-chan time.Time
	)
	if runner.Options.Timeout != 0 {
		timer = time.NewTimer(runner.Options.Timeout)
		timerC = timer.C
	}

	select {
	case <-stdCh:
	case <-timerC:
		ctx.Writer.Write([]byte("Error: timeout"))
		return
	}

	ctx.Writer.Write([]byte(sb.String()))
}
