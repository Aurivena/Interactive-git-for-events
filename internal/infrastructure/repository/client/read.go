package client

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func (r *Client) Read(sessionID string) (json.RawMessage, error) {
	var output json.RawMessage

	if err := r.db.Get(&output, `SELECT survey FROM client WHERE session_id = $1`, sessionID); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return output, nil
}
