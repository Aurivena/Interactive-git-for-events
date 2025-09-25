package http

import (
	"arch/internal/application"
	"arch/internal/delivery/http/ai"
	"arch/internal/delivery/http/client"
	"arch/internal/delivery/http/history"
	"arch/internal/delivery/http/place"
	"arch/internal/delivery/http/tour"
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
	History    *history.Handler
	Client     *client.Handler
	Tour       *tour.Handler
	Middleware *middleware.Middleware
}

func NewHttp(application *application.Application, spond *core.Spond, middleware *middleware.Middleware) *Http {
	return &Http{
		Ai:         ai.New(application, spond),
		Place:      place.New(application, spond),
		History:    history.New(application, spond),
		Client:     client.New(application, spond),
		Tour:       tour.New(application, spond),
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

		aiRouter := api.Group("/ai", h.Middleware.Session)
		{
			aiRouter.POST("/send", h.Ai.Send)
			aiRouter.POST("/generate/tour", h.Ai.GenerateTour)
			aiRouter.GET("/history", h.History.ListHistory)
		}

		tours := api.Group("/tours", h.Middleware.Session)
		{
			tours.GET("", h.Tour.All)
			tours.GET("/:id", h.Tour.ByID)
		}

		clientApp := api.Group("/client", h.Middleware.Session)
		{
			clientApp.POST("/upsert", h.Client.Upsert)
		}

		places := api.Group("/places")
		{
			places.GET("", h.Place.List)
			places.GET("/kind/:kind", h.Place.ListByKind)
			places.GET("/:id", h.Place.ByID)
			places.GET("/image/:id", h.Place.ImageByID)
		}

	}

	return gHttp
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
