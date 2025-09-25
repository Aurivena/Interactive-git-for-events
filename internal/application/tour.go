package application

import "arch/internal/domain/entity"

func (a *Application) TourAll(sessionID string) ([]entity.TourOutput, error) {
	output, err := a.post.TourReader.Reader(sessionID)
	if err != nil {
		return nil, err
	}

	return output, err
}

func (a *Application) TourByID(id entity.UUID) (*entity.TourOutput, error) {
	output, err := a.post.TourReader.ReaderByID(id)
	if err != nil {
		return nil, err
	}

	return output, err
}
