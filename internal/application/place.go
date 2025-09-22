package application

import "arch/internal/domain/entity"

func (a *Application) List() ([]entity.PlaceInfo, error) {
	output, err := a.post.PlaceReader.List()
	if err != nil {
		return nil, err
	}

	for i := range output {
		images, err := a.post.PlaceReader.ImagesByPlaceID(output[i].ID)
		if err != nil {
			return nil, err
		}
		output[i].Images = images
	}
	return output, nil
}

func (a *Application) ListByKind(kind entity.Kind) ([]entity.PlaceInfo, error) {
	output, err := a.post.PlaceReader.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	for i := range output {
		images, err := a.post.PlaceReader.ImagesByPlaceID(output[i].ID)
		if err != nil {
			return nil, err
		}
		output[i].Images = images
	}
	return output, nil
}

func (a *Application) ByID(id entity.UUID) (*entity.PlaceInfo, error) {
	output, err := a.post.PlaceReader.ByID(id)
	if err != nil {
		return nil, err
	}
	images, err := a.post.PlaceReader.ImagesByPlaceID(output.ID)
	if err != nil {
		return nil, err
	}
	output.Images = images
	return output, nil
}

func (a *Application) ImageByID(id entity.UUID) ([]byte, string, error) {
	data, contentType, err := a.post.MinioReader.GetImage(id)
	if err != nil {
		return nil, "", err
	}

	return data, contentType, nil
}
