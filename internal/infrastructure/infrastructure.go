package infrastructure

import (
	"arch/internal/domain/entity"
	"arch/internal/infrastructure/repository/history"
	"arch/internal/infrastructure/repository/place"
	"arch/internal/infrastructure/s3"
	"arch/internal/ports"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type Infrastructure struct {
	PlaceReader  ports.PlaceReader
	PlaceWriter  ports.PlaceWriter
	PlaceBinding ports.PlaceBinding

	HistoryWriter ports.HistoryWriter
	HistoryReader ports.HistoryReader
	MinioWriter   ports.MinioWrite
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(sources *Sources, client *minio.Client, cfg entity.MinioConfig) *Infrastructure {
	return &Infrastructure{
		PlaceReader:  place.NewPlace(sources.BusinessDB),
		PlaceWriter:  place.NewPlace(sources.BusinessDB),
		PlaceBinding: place.NewPlace(sources.BusinessDB),

		HistoryWriter: history.NewHistory(sources.BusinessDB),
		HistoryReader: history.NewHistory(sources.BusinessDB),
		MinioWriter:   s3.New(client, cfg),
	}
}
