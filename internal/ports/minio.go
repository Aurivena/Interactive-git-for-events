package ports

import (
	"arch/internal/domain/entity"
	"context"
)

type MinioWrite interface {
	Write(ctx context.Context, data []byte, fileID string) error
}

type MinioReader interface {
	GetImage(id entity.UUID) ([]byte, string, error)
}
