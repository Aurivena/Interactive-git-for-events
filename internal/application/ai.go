package application

import (
	"arch/internal/domain/ai"
	"arch/internal/domain/entity"
	"encoding/json"
	"time"
)

func (a *Application) SendAi(input entity.UserSend, sessionID string) ([]entity.ChatOutput, error) {
	q := ai.New(*a.qwqConfig)

	if input.Istest == true {
		time.Sleep(1 * time.Second)
		for _, val := range ExampleChatOutputs {
			if err := a.post.History.Save(val, input.Message, sessionID); err != nil {
				return nil, err
			}
		}
		return ExampleChatOutputs, nil
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
		out, err := a.post.PlaceGet.Get(&params[i], input.Lon, input.Lat)
		if err != nil {
			return nil, err
		}

		if out != nil {
			output[i].PlaceInfo = out
			output[i].Message = params[i].Message
		}

		if err = a.post.History.Save(output[i], input.Message, sessionID); err != nil {
			return nil, err
		}
	}

	return output, nil
}

var ExampleChatOutputs = []entity.ChatOutput{
	{
		PlaceInfo: []entity.PlaceInfo{
			{
				ID:      "11111111-1111-1111-1111-111111111111",
				Title:   "Россия",
				Kind:    "cinema",
				Address: "Курган, ул. Володарского, 75",
				Lon:     65.339361,
				Lat:     55.439371,
				Tags: json.RawMessage(`{
					"phone": "+7 (3522) 60-52-50",
					"website": "https://россия45.рф",
					"schedule": [
						{"end": "21:00", "week": "monday", "start": "09:00", "spans_midnight": false},
						{"end": "21:00", "week": "friday", "start": "09:00", "spans_midnight": false}
					]
				}`),
			},
			{
				ID:      "11111111-1111-1111-1111-111111111112",
				Title:   "Pushka",
				Kind:    "cinema",
				Address: "Курган, ул. Пушкина, 25, ТРЦ «Пушкинский», 3 этаж",
				Lon:     65.318954,
				Lat:     55.432190,
				Tags: json.RawMessage(`{
					"phone": "+7 (3522) 60-70-55",
					"website": "https://cinema.pushka.club/kurgan/pushka",
					"schedule": [
						{"end": "23:00", "week": "saturday", "start": "09:00", "spans_midnight": false},
						{"end": "23:00", "week": "sunday", "start": "09:00", "spans_midnight": false}
					]
				}`),
			},
		},
		Message: "Нашёл для тебя пару кинотеатров 🎬 — «Россия» для классического вечера и «Pushka» для современного киношного вайба.",
	},
	{
		PlaceInfo: []entity.PlaceInfo{
			{
				ID:      "11111111-1111-1111-1111-111111111113",
				Title:   "Ultra Cinema Курган",
				Kind:    "cinema",
				Address: "Курган, ул. Коли Мяготина, 8, ТРЦ «Hyper City», 2 этаж",
				Lon:     65.280027,
				Lat:     55.426618,
				Tags: json.RawMessage(`{
					"phone": "+7 (3522) 22-89-87",
					"website": "https://kurgan.ultra-cinema.ru",
					"schedule": [
						{"end": "00:00", "week": "friday", "start": "10:00", "spans_midnight": true}
					]
				}`),
			},
			{
				ID:      "11111111-1111-1111-1111-111111111114",
				Title:   "Клумба Синема",
				Kind:    "cinema",
				Address: "Курган, 2-й микрорайон, 17, ТРЦ «Стрекоза», 2 этаж",
				Lon:     65.264725,
				Lat:     55.464850,
				Tags: json.RawMessage(`{
					"phone": "+7 (963) 869-80-49",
					"website": "https://klumba-cinema.ru",
					"schedule": [
						{"end": "02:00", "week": "saturday", "start": "09:00", "spans_midnight": true}
					]
				}`),
			},
		},
		Message: "Ещё два варианта 🎥 — «Ultra Cinema Курган» для ночных сеансов и «Клумба Синема» с поздними показами до 2-х ночи.",
	},
}
