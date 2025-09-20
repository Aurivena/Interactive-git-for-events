package initialization

import (
	"arch/internal/application"
	"arch/internal/delivery/http"
	"arch/internal/delivery/middleware"
	"arch/internal/infrastructure"

	"github.com/Aurivena/spond/v2/core"
	"github.com/jmoiron/sqlx"
)

func InitLayers() (delivery *http.Http, businessDatabase *sqlx.DB) {
	spond := core.NewSpond()
	businessDatabase = infrastructure.NewBusinessDatabase(ConfigService)
	sources := infrastructure.Sources{
		BusinessDB: businessDatabase,
	}
	infrastructures := infrastructure.New(&sources)
	client := newVpnClient()
	app := application.New(infrastructures, &ConfigService.QwQ, client)
	middleware := middleware.New()
	delivery = http.NewHttp(app, spond, middleware)
	return delivery, businessDatabase
}
