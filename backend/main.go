package main

import (
	"github.com/gin-gonic/gin"
	"github.com/user/kahoot-clone/controllers"
	"github.com/user/kahoot-clone/db"
	"github.com/user/kahoot-clone/middleware"
	"github.com/user/kahoot-clone/socket"
)

func main() {
	db.ConnectDB()
	controllers.SeedAdmin()

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public Routes
	r.POST("/api/login", controllers.Login)
	r.Static("/uploads", "./uploads")

	// Protected Admin Routes
	admin := r.Group("/api")
	admin.Use(middleware.JWTAuth())
	{
		admin.POST("/games", controllers.CreateGame)
		admin.GET("/games", controllers.GetGames)
		admin.GET("/games/:id", controllers.GetGame)
		admin.PUT("/games/:id", controllers.UpdateGame)
		admin.DELETE("/games/:id", controllers.DeleteGame)
		admin.POST("/games/:id/duplicate", controllers.DuplicateGame)
		admin.POST("/upload", controllers.UploadFile)

		// History Routes
		admin.GET("/history", controllers.GetHistory)
		admin.GET("/history/:id", controllers.GetHistoryDetail)
		admin.POST("/history", controllers.SaveHistory)
	}

	// WebSocket Route
	r.GET("/ws", func(c *gin.Context) {
		socket.HandleConnection(c.Writer, c.Request)
	})

	r.Run(":8081")
}
