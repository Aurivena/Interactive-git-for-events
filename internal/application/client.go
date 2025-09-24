package application

import "arch/internal/domain/entity"

func (a *Application) UpsertClientSurvey(sessionID string, survey entity.Survey) error {
	if err := a.post.ClientWriter.Write(sessionID, survey); err != nil {
		return err
	}

	return nil
}
