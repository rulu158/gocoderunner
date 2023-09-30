package main

import "github.com/gin-gonic/gin"

type Server struct {
	status int32
}

func NewServer() *Server {
	srv := &Server{}
	return srv
}

func (srv *Server) Run(port string) {
	engine := gin.New()
	_ = engine

	engine.POST("/api/exec-runner", srv.ExecRunnerPOST)

	engine.Run(port)
}
