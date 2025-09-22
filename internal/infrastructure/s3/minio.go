package s3

import (
	"arch/internal/domain/entity"

	"github.com/minio/minio-go/v7"
)

type S3 struct {
	minioClient *minio.Client
	cfg         entity.MinioConfig
}

func New(minioClient *minio.Client, cfg entity.MinioConfig) *S3 {
	return &S3{
		minioClient: minioClient,
		cfg:         cfg,
	}
}
