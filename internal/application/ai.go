package application

import (
	"arch/internal/domain/ai"
	"arch/internal/domain/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	centerLon = 65.340956
	centerLat = 55.439635
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
		for i := range examples {
			images, err := a.post.PlaceReader.ImagesByPlaceID(examples[i].ID)
			if err != nil {
				return nil, err
			}
			examples[i].Images = images
		}
		ot := []entity.ChatOutput{
			{
				PlaceInfo: examples,
				Message:   input.Message,
			},
		}

		if err := a.post.HistoryWriter.Write(ot, input.Message, sessionID); err != nil {
			return nil, err
		}
		return ot, nil
	}
	survey, err := a.post.ClientReader.Read(sessionID)
	if err != nil {
		return nil, err
	}
	params, err := q.Send(input.Message, ai.SendPrompt, survey)
	if err != nil {
		return nil, err
	}

	aiOutput, err := ai.BuildOutput[[]entity.RequestPayload](bytes.NewReader(params))
	if err != nil {
		return nil, err
	}

	output := make([]entity.ChatOutput, len(params))

	for i := range aiOutput {

		if aiOutput[i].Count == 0 {
			aiOutput[i].Count = ai.DefaultCount
		}
		out, err := a.post.PlaceReader.Get(&aiOutput[i], input.Lon, input.Lat)
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
			output[i].Message = aiOutput[i].Message
		}
	}
	if err = a.post.HistoryWriter.Write(output, input.Message, sessionID); err != nil {
		return nil, err
	}

	return output, nil
}

func (a *Application) GenerateTour(input *entity.TourInput, sessionID string) (*entity.TourOutput, error) {
	survey, err := a.post.ClientReader.Read(sessionID)
	if err != nil {
		return nil, err
	}

	q := ai.New(*a.aiConfig)
	message := fmt.Sprintf("Составь тут на даты от %s до %s", input.DateFrom, input.DateTo)

	var aiOutput entity.RouteParams
	if input.IsTest {
		aiOutput = entity.RouteParams{
			DateTour: entity.DateTour{
				DateFrom: "2025-09-25",
				DateTo:   "2025-09-27",
			},
			PerDayLimit:  5,
			Tier:         []string{"standard", "economy", "value", "premium"},
			KindPriority: []entity.Kind{"cinema", "historic", "museum", "park", "restaurant"},
			DayStart:     "10:00",
			DayEnd:       "22:00",
		}

	} else {
		params, err := q.Send(message, ai.RouteParamsFromSurveyPrompt, survey)
		if err != nil {
			return nil, err
		}
		aiOutput, err = ai.BuildOutput[entity.RouteParams](bytes.NewReader(params))
		if err != nil {
			return nil, err
		}
	}

	checkCoordinates(&input.Lat, &input.Lon)
	aiOutput.DateFrom = input.DateFrom
	aiOutput.DateTo = input.DateTo
	raw, err := a.post.TourGenerates.GenerateTour(aiOutput, *input.Lon, *input.Lat)
	if err != nil {
		return nil, err
	}

	var tour entity.Tour
	if err = json.Unmarshal(raw, &tour); err != nil {
		logrus.Error(err)
		return nil, err
	}

	id, err := a.post.TourWriter.Write(aiOutput.DateFrom, aiOutput.DateTo, sessionID, tour)
	if err != nil {
		return nil, err
	}

	return &entity.TourOutput{
		ID:   *id,
		Tour: tour,
	}, nil
}

func checkCoordinates(lat, lon **float64) {
	if *lat == nil {
		v := centerLat
		*lat = &v
	}
	if *lon == nil {
		v := centerLon
		*lon = &v
	}
	if **lat < -90 || **lat > 90 {
		if **lon >= -90 && **lon <= 90 {
			tmp := **lat
			**lat = **lon
			**lon = tmp
		}
	}
}
