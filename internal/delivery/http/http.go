package http

import (
	"arch/internal/application"
	"arch/internal/delivery/http/ai"
	"arch/internal/delivery/http/place"
	"arch/internal/delivery/middleware"
	"arch/internal/domain/entity"
	"arch/internal/server"
	"strings"
	"time"

	"github.com/Aurivena/spond/v2/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Http struct {
	Ai         *ai.Handler
	Place      *place.Handler
	Middleware *middleware.Middleware
}

func NewHttp(application *application.Application, spond *core.Spond, middleware *middleware.Middleware) *Http {
	return &Http{
		Ai:         ai.New(application, spond),
		Place:      place.New(application, spond),
		Middleware: middleware,
	}
}

func (h *Http) InitHTTPHttps(config *entity.ServerConfig) *gin.Engine {
	ginSetMode(config.ServerMode)
	gHttp := gin.Default()
	allowOrigins := strings.Split(config.Domain, ",")

	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"X-Session-ID", "X-Password", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := gHttp.Group("/api")
	{

		aiRouter := api.Group("/ai")
		{
			aiRouter.POST("/send", h.Ai.Send)
		}

		placeRouter := api.Group("/place")
		{
			placeRouter.GET("/list", h.Place.List)
			placeRouter.GET("/list/:kind", h.Place.ListByKind)
			placeRouter.GET("/:id", h.Place.ByID)
		}

	}

	return gHttp
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
