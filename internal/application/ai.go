package application

import (
	"arch/internal/domain/ai"
	"arch/internal/domain/entity"
	"time"
)

func (a *Application) SendAi(input entity.UserSend, sessionID string) ([]entity.ChatOutput, error) {
	q := ai.New(*a.aiConfig)

	if len(input.Message) == 0 {
		return []entity.ChatOutput{
			{
				Message:   "Извините, но ваше сообщение пустое 👀!",
				PlaceInfo: []entity.PlaceInfo{},
			},
		}, nil
	}
	if input.Istest == true {
		time.Sleep(1 * time.Second)
		examples, err := a.post.PlaceReader.ListByKind("cinema")
		if err != nil {
			return nil, err
		}
		ot := []entity.ChatOutput{
			entity.ChatOutput{
				PlaceInfo: examples,
				Message:   input.Message,
			},
		}

		if err := a.post.HistoryWriter.Write(ot, input.Message, sessionID); err != nil {
			return nil, err
		}
		return ot, nil
	}
	params, err := q.Send(input.Message)
	if err != nil {
		return nil, err
	}

	output := make([]entity.ChatOutput, len(params))

	for i := range params {

		if params[i].Count == 0 {
			params[i].Count = ai.DefaultCount
		}
		out, err := a.post.PlaceReader.Get(&params[i], input.Lon, input.Lat)
		if err != nil {
			return nil, err
		}

		for i := range out {
			images, err := a.post.PlaceReader.ImagesByPlaceID(out[i].ID)
			if err != nil {
				return nil, err
			}
			out[i].Images = images
		}

		if out != nil {
			output[i].PlaceInfo = out
			output[i].Message = params[i].Message
		}
	}
	if err = a.post.HistoryWriter.Write(output, input.Message, sessionID); err != nil {
		return nil, err
	}

	return output, nil
}
