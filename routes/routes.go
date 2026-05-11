package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", middleware.AuthRequired(), handlers.Logout)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/me", handlers.Me)

			boards := protected.Group("/boards")
			{
				boards.GET("", handlers.ListBoards)
				boards.POST("", handlers.CreateBoard)
				boards.GET("/:id", handlers.GetBoard)
				boards.PATCH("/:id", handlers.UpdateBoard)
				boards.DELETE("/:id", handlers.DeleteBoard)
				boards.POST("/:id/cards", handlers.CreateCard)
				boards.GET("/:id/labels", handlers.ListLabels)
				boards.POST("/:id/labels", handlers.CreateLabel)
			}

			cards := protected.Group("/cards")
			{
				cards.PATCH("/:id", handlers.UpdateCard)
				cards.PATCH("/:id/move", handlers.MoveCard)
				cards.DELETE("/:id", handlers.DeleteCard)
				cards.POST("/:id/labels/:labelId", handlers.AttachLabel)
				cards.DELETE("/:id/labels/:labelId", handlers.DetachLabel)
				cards.POST("/:id/subcards", handlers.CreateSubCard)
			}

			subcards := protected.Group("/subcards")
			{
				subcards.PATCH("/:id", handlers.UpdateSubCard)
				subcards.DELETE("/:id", handlers.DeleteSubCard)
			}

			labels := protected.Group("/labels")
			{
				labels.PATCH("/:id", handlers.UpdateLabel)
				labels.DELETE("/:id", handlers.DeleteLabel)
			}
		}
	}

	return r
}
