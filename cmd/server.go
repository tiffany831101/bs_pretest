package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tiffany831101/bs_pretest.git/internal/controller"
)

type Server struct {
	engine *gin.Engine
}

func StartServer() *Server {
	router := gin.Default()
	return &Server{
		engine: router,
	}
}

func (s *Server) Run() {
	port := viper.GetString("server.port")
	s.engine.Run(":" + port)
}

func (s *Server) SetUpRoutes() {

	controller.NewTasksController()
	controller.SetUpTasksRoutes(s.engine)
}
