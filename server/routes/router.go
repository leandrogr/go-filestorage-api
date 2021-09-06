package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leandrogr/go-filestorage-api/controllers"
	"github.com/leandrogr/go-filestorage-api/services"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
	}
}

func ConfigRoutes(router *gin.Engine) *gin.Engine {

	router.Use(CORSMiddleware())
	router.OPTIONS("/*path", CORSMiddleware())

	sess := services.ConnectAWS()
	router.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
	})

	main := router.Group("api/v1")
	{
		files := main.Group("files")
		{
			files.GET("/", controllers.ShowFiles)
			files.GET("/info", controllers.GetFile)
			files.GET("/download", controllers.DownloadFile)
			files.POST("/", controllers.CreateFile)
			files.POST("/copy", controllers.CopyFile)
			files.DELETE("/", controllers.DeleteFile)
		}
	}

	return router
}
