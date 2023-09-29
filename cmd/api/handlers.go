package main

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var sbStdout, sbStderr strings.Builder
	runner := runner.NewRunner(languages.Go, &runner.RunnerOptions{
		Stdout:  &sbStdout,
		Stderr:  &sbStderr,
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
		response := &Response{ID: "", Error: true, Result: "Timeout"}
		c.JSON(http.StatusRequestTimeout, response)
		return
	}

	var id, result string
	var status int
	if isError {
		id = ""
		status = http.StatusInternalServerError
		result = sbStderr.String()
	} else {
		id = runner.ID
		status = http.StatusOK
		result = sbStdout.String()
	}
	response := &Response{ID: id, Error: isError, Result: result}
	c.JSON(status, response)
}
