package s3

import (
	"arch/internal/domain"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (m *S3) Write(ctx context.Context, data []byte, fileID string) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}

	_, err := m.minioClient.StatObject(ctx, m.cfg.MinioBucketName, fileID, minio.StatObjectOptions{})
	if err == nil {
		logrus.Info("File already exists")
		return domain.FileDuplicate
	}
	if minio.ToErrorResponse(err).Code != "NoSuchKey" {
		return fmt.Errorf("ошибка при проверке существования файла: %w", err)
	}

	_, err = m.minioClient.PutObject(ctx, m.cfg.MinioBucketName, fileID, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: http.DetectContentType(data),
	})
	if err != nil {
		return err
	}

	return nil
}
