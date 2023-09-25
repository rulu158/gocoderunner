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

	/*engine.GET("/new-runner", srv.GetNewRunner)
	engine.POST("/stop-runner", srv.StopRunner)*/
	engine.GET("/exec-runner", srv.ExecRunner)

	engine.Run(port)
}
