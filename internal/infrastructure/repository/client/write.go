package client

import (
	"arch/internal/domain/entity"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func (r *Client) Write(sessionID string, survey entity.Survey) error {
	jsonb, err := json.Marshal(survey)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = r.db.Exec(`INSERT INTO client(session_id,survey) VALUES ($1,$2)`, sessionID, jsonb)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
