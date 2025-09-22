package s3

import (
	"arch/internal/domain"
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/minio/minio-go/v7"
)

func (m *S3) Write(ctx context.Context, data []byte, fileID string) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}

	_, err := m.minioClient.PutObject(ctx, m.cfg.MinioBucketName, fileID, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: http.DetectContentType(data),
	})
	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == 412 {
			return domain.FileDuplicate
		}
		return err
	}

	return nil
}
