package place

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (r *Place) Write(id uuid.UUID, sql string) error {
	_, err := r.db.Exec(sql, id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
