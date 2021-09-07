package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leandrogr/go-filestorage-api/server/routes"
)

type Server struct {
	port   string
	server *gin.Engine
}

func NewServer() Server {
	return Server{
		port:   os.Getenv("PORT"),
		server: gin.Default(),
	}
}

func (s *Server) Run() {
	router := routes.ConfigRoutes(s.server)

	log.Print("Server is running at port: ", s.port)
	log.Fatal(router.Run(":" + s.port))
}
