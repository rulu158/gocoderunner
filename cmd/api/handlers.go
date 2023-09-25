package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rulu158/gocoderunner/runner"
)

func (srv *Server) ExecRunner(ctx *gin.Context) {
	var sb strings.Builder
	runner := *runner.NewRunner(&runner.RunnerOptions{Stdout: &sb})
	go runner.ExecCode()
	ctx.Writer.Write([]byte(sb.String()))
}
