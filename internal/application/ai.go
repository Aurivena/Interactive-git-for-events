package application

import (
	"arch/internal/domain/ai"
	"arch/internal/domain/entity"
	"context"
)

func (a *Application) SendAi(ctx context.Context, input entity.UserSend) ([]entity.ChatOutput, error) {
	q := ai.New(*a.qwqConfig)
	params, err := q.Send(ctx, input.Message)
	if err != nil {
		return nil, err
	}

	output := make([]entity.ChatOutput, len(params))

	for i := range params {
		if params[i].Count == 0 {
			params[i].Count = ai.DefaultCount
		}
		out, err := a.post.PlaceGet.Get(&params[i], input.Lon, input.Lat)
		if err != nil {
			return nil, err
		}

		if out != nil {
			output[i].PlaceInfo = out
			output[i].Message = params[i].Message
		}
	}

	return output, nil
}
