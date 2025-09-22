package initialization

import (
	"arch/internal/application"
	"arch/internal/delivery/http"
	"arch/internal/delivery/middleware"
	"arch/internal/infrastructure"
	"arch/internal/migrations"

	"github.com/Aurivena/spond/v2/core"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func InitLayers() (delivery *http.Http, businessDatabase *sqlx.DB) {
	spond := core.NewSpond()
	businessDatabase = infrastructure.NewBusinessDatabase(ConfigService)
	sources := infrastructure.Sources{
		BusinessDB: businessDatabase,
	}

	s3Minio := NewMinioStorage(ConfigService.Minio)
	infrastructures := infrastructure.New(&sources, s3Minio, ConfigService.Minio)

	mgr := migrations.New(infrastructures.PlaceWriter, infrastructures.PlaceBinding, infrastructures.MinioWriter)
	if err := mgr.DownloadImages(); err != nil {
		logrus.Warnf("Failed to download images: %v", err)
		logrus.Errorf("Failed to download images: %v", err)
	}

	app := application.New(infrastructures, &ConfigService.Ai)
	middleware := middleware.New(spond)
	delivery = http.NewHttp(app, spond, middleware)
	return delivery, businessDatabase
}
