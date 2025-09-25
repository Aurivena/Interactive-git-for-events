package infrastructure

import (
	"arch/internal/domain/entity"
	"arch/internal/infrastructure/repository/history"
	"arch/internal/infrastructure/repository/place"
	"arch/internal/infrastructure/s3"
	"arch/internal/infrastructure/tour"
	"arch/internal/ports"

	client_app "arch/internal/infrastructure/repository/client"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type Infrastructure struct {
	PlaceReader  ports.PlaceReader
	PlaceWriter  ports.PlaceWriter
	PlaceBinding ports.PlaceBinding

	HistoryWriter ports.HistoryWriter
	HistoryReader ports.HistoryReader

	MinioWriter ports.MinioWrite
	MinioReader ports.MinioReader

	ClientWriter ports.ClientWrite
	ClientReader ports.ClientReade

	TourWriter    ports.TourWrite
	TourReader    ports.TourReader
	TourGenerates ports.TourGenerate
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(sources *Sources, client *minio.Client, cfg entity.MinioConfig) *Infrastructure {
	return &Infrastructure{
		PlaceReader:  place.New(sources.BusinessDB),
		PlaceWriter:  place.New(sources.BusinessDB),
		PlaceBinding: place.New(sources.BusinessDB),

		HistoryWriter: history.New(sources.BusinessDB),
		HistoryReader: history.New(sources.BusinessDB),

		MinioWriter: s3.New(client, cfg),
		MinioReader: s3.New(client, cfg),

		ClientWriter: client_app.New(sources.BusinessDB),
		ClientReader: client_app.New(sources.BusinessDB),

		TourWriter:    tour.New(sources.BusinessDB),
		TourReader:    tour.New(sources.BusinessDB),
		TourGenerates: tour.New(sources.BusinessDB),
	}
}
