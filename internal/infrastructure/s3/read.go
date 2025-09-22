package s3

import (
	"arch/internal/domain/entity"
	"context"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (m *S3) GetImage(id entity.UUID) ([]byte, string, error) {
	var data []byte
	var contentType string

	ctx := context.Background()
	obj, err := m.minioClient.GetObject(ctx, m.cfg.MinioBucketName, string(id), minio.GetObjectOptions{})
	if err != nil {
		logrus.Error("minioClient.GetObject failed", err)
		return nil, "", err
	}
	defer obj.Close()

	if st, statErr := obj.Stat(); statErr == nil && st.ContentType != "" {
		contentType = st.ContentType
	}

	data, err = io.ReadAll(obj)
	if err != nil {
		return nil, "", err
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	return data, contentType, nil
}
