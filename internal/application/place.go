package application

import "arch/internal/domain/entity"

func (a *Application) List() ([]entity.PlaceInfo, error) {
	output, err := a.post.PlaceReader.List()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (a *Application) ListByKind(kind entity.Kind) ([]entity.PlaceInfo, error) {
	output, err := a.post.PlaceReader.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (a *Application) ByID(id entity.UUID) (*entity.PlaceInfo, error) {
	output, err := a.post.PlaceReader.ByID(id)
	if err != nil {
		return nil, err
	}
	return output, nil
}
