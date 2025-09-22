package history

import (
	"arch/internal/domain/entity"

	"github.com/sirupsen/logrus"
)

func (r *History) Write(aiMessage []entity.ChatOutput, message string, sessionID string) error {
	_, err := r.db.Exec(`INSERT INTO  history (session,message,ai_message) VALUES ($1,$2,$3)`, sessionID, message, aiMessage)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
