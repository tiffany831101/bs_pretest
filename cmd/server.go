package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/tiffany831101/bs_pretest.git/docs"
	"github.com/tiffany831101/bs_pretest.git/internal/controller"
)

type Server struct {
	engine *gin.Engine
}

func StartServer() *Server {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
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

func (s *Server) RunSwagger() {
	docs.SwaggerInfo.BasePath = "/api/" + viper.GetString("api.version")
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
