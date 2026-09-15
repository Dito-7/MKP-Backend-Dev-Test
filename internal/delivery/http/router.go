package http

import (
	"net/http"

	"mkp-cinema-ticketing/internal/config"
	"mkp-cinema-ticketing/internal/delivery/http/handler"
	"mkp-cinema-ticketing/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	scheduleHandler *handler.ScheduleHandler,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"message": "Cinema Ticketing API is running smoothly",
			"env":     cfg.AppEnv,
		})
	})

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
		api.POST("/register", authHandler.Register)
		api.GET("/profile", middleware.AuthMiddleware(cfg), authHandler.Profile)

		schedules := api.Group("/schedules")
		schedules.Use(middleware.AuthMiddleware(cfg))
		{
			schedules.POST("", scheduleHandler.Create)
			schedules.GET("", scheduleHandler.GetAll)
			schedules.GET("/:id", scheduleHandler.GetByID)
			schedules.PUT("/:id", scheduleHandler.Update)
			schedules.DELETE("/:id", scheduleHandler.Delete)
		}
	}

	return r
}
