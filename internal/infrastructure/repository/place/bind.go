package place

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (r *Place) Bind(placeID uuid.UUID, imageID string) error {
	_, err := r.db.Exec(`INSERT INTO place_image(place_id,image_id) VALUES ($1,$2)`, placeID, imageID)
	if err != nil {
		logrus.Error("error on insert place_image", err)
		return err
	}
	return nil
}
