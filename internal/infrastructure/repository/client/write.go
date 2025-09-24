package client

import (
	"arch/internal/domain/entity"

	"github.com/sirupsen/logrus"
)

func (r *Client) Write(sessionID string, survey entity.Survey) error {
	_, err := r.db.Exec(`INSERT INTO client(session_id,survey) VALUES ($1,$2)`, sessionID, survey)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
