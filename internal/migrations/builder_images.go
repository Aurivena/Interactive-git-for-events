package migrations

import (
	"arch/internal/domain"
	"arch/internal/ports"
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var rootDir = "internal/migrations/object"

type Migrations struct {
	writer      ports.PlaceWriter
	binding     ports.PlaceBinding
	minioWriter ports.MinioWrite
}

type imageBlob struct {
	Name string
	Data []byte
}

type objectInfo struct {
	Dir    string
	TxtSQL string
	Images []imageBlob
}

func New(writer ports.PlaceWriter, binding ports.PlaceBinding, minioWriter ports.MinioWrite) *Migrations {
	return &Migrations{
		writer:      writer,
		binding:     binding,
		minioWriter: minioWriter,
	}
}

func (m *Migrations) DownloadImages() error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return fmt.Errorf("read root dir %q: %w", rootDir, err)
	}

	var firstErr error

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".zip") {
			continue
		}

		zipPath := filepath.Join(rootDir, e.Name())
		zr, err := zip.OpenReader(zipPath)
		if err != nil {
			logrus.WithError(err).WithField("zip", zipPath).Error("open zip")
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		objs, err := reader(zr.File)
		_ = zr.Close()
		if err != nil {
			logrus.WithError(err).WithField("zip", zipPath).Error("read zip entries")
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err = m.execute(objs); err != nil {
			logrus.WithError(err).WithField("zip", zipPath).Error("execute")
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
	}

	return firstErr
}

func (m *Migrations) execute(objs []objectInfo) error {
	ctx := context.Background()
	for _, obj := range objs {
		txtUUid, err := uuid.NewV7()
		if err != nil {
			logrus.Error("Error at creating UUID")
			return err
		}

		err = m.writer.Write(txtUUid, obj.TxtSQL)
		if errors.Is(err, domain.FileDuplicate) {
			continue
		}
		if err != nil {
			logrus.Error("Error at creating UUID")
			return err
		}

		for _, image := range obj.Images {
			err = m.minioWriter.Write(ctx, image.Data, image.Name)
			if errors.Is(err, domain.FileDuplicate) {
				logrus.Warn("File already exists")
				continue
			}
			if err != nil {
				logrus.Error("Error at write file")
				return err
			}

			if err = m.binding.Bind(txtUUid, image.Name); err != nil {
				logrus.Error("Error at bind file")
				return err
			}
		}
	}
	return nil
}

func reader(files []*zip.File) ([]objectInfo, error) {
	byDir := make(map[string]*objectInfo)

	for _, zf := range files {
		if zf.FileInfo().IsDir() {
			continue
		}

		if strings.HasPrefix(zf.Name, "__MACOSX/") || strings.HasPrefix(path.Base(zf.Name), "._") {
			continue
		}

		dir := path.Dir(zf.Name)
		base := path.Base(zf.Name)
		ext := strings.ToLower(path.Ext(base))

		r, err := zf.Open()
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", zf.Name, err)
		}

		switch ext {
		case ".txt":
			b, err := io.ReadAll(r)
			_ = r.Close()
			if err != nil {
				return nil, fmt.Errorf("read txt %q: %w", zf.Name, err)
			}

			if bytes.IndexByte(b, 0x00) >= 0 {
				logrus.WithField("file", zf.Name).Warn("skip binary txt (contains NUL)")
				continue
			}

			get(byDir, dir).TxtSQL = string(b)

		default:
			if isImageExt(ext) {
				if err = appendFile(r, byDir, base, ext, dir); err != nil {
					return nil, fmt.Errorf("append image %q: %w", zf.Name, err)
				}
			} else {
				_ = r.Close()
			}
		}
	}

	out := make([]objectInfo, 0, len(byDir))
	for _, v := range byDir {
		out = append(out, *v)
	}
	return out, nil
}

func get(byDir map[string]*objectInfo, dir string) *objectInfo {
	if o := byDir[dir]; o != nil {
		return o
	}
	o := &objectInfo{Dir: dir}
	byDir[dir] = o
	return o
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

func appendFile(r io.ReadCloser, byDir map[string]*objectInfo, base, ext, dir string) error {
	b, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil {
		return err
	}
	oi := get(byDir, dir)
	name := strings.TrimSuffix(base, ext)
	oi.Images = append(oi.Images, imageBlob{Name: name, Data: b})
	return nil
}
