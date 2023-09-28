package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rulu158/gocoderunner/runner"
	"github.com/rulu158/gocoderunner/runner/languages"
)

/*
func (srv *Server) ExecRunner(c *gin.Context) {
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
		c.Writer.Write([]byte("Error: timeout"))
		return
	}

	c.Writer.Write([]byte(sb.String()))
}
*/

func (srv *Server) ExecRunnerPOST(c *gin.Context) {
	var codeItem CodePOST
	if err := c.BindJSON(&codeItem); err != nil {
		response := &Response{ID: "", Error: true, Result: "Invalid JSON."}
		res, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		c.Writer.Write(res)
		return
	}

	var sb strings.Builder
	runner := runner.NewRunner(languages.Go, &runner.RunnerOptions{
		Stdout:  &sb,
		Timeout: 30 * time.Second,
	})

	okCh := make(chan int, 1)
	koCh := make(chan int, 1)
	go func() {
		err := runner.ExecCode([]byte(codeItem.Code))
		if err != nil {
			koCh <- 1
		} else {
			okCh <- 1
		}
	}()

	var (
		timer  *time.Timer
		timerC <-chan time.Time
	)
	if runner.Options.Timeout != 0 {
		timer = time.NewTimer(runner.Options.Timeout)
		timerC = timer.C
	}

	isError := false
	select {
	case <-okCh:
	case <-koCh:
		isError = true
	case <-timerC:
		c.Writer.Write([]byte("Error: timeout"))
		return
	}

	response := &Response{ID: runner.ID, Error: isError, Result: sb.String()}
	res, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	c.Writer.Write(res)
}
