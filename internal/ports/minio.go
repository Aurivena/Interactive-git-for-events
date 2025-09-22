package ports

import (
	"context"
)

type MinioWrite interface {
	Write(ctx context.Context, data []byte, fileID string) error
}
