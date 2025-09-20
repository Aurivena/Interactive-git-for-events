package entity

import "encoding/json"

type UserSend struct {
	Message string   `json:"message" example:"Хочу сходить в кино и ресторан"`
	Lat     *float64 `json:"lat,omitempty" example:"55.434630"`
	Lon     *float64 `json:"lon,omitempty" example:"65.353470"`
}

type PlaceInfo struct {
	ID          string          `json:"id" db:"id"`
	Title       string          `json:"title" db:"title"`
	Kind        string          `json:"kind" db:"kind"`
	Description *string         `json:"description" db:"description"`
	Address     string          `json:"address" db:"address"`
	Lon         float64         `json:"lon" db:"lon"`
	Lat         float64         `json:"lat" db:"lat"`
	Tags        json.RawMessage `json:"tags" db:"tags" swaggertype:"object"`
}

type RequestPayload struct {
	Radius       int       `json:"radius"`
	Count        int       `json:"count"`
	Kind         string    `json:"kind"`
	Tier         string    `json:"tier"`
	DayOfTheWeek []Weekday `json:"dayOfTheWeek"`
	Time         *TimeOnly `json:"time"`
	Message      string    `json:"message"`
}

type ChatOutput struct {
	PlaceInfo []PlaceInfo `json:"placesInfo"`
	Message   string      `json:"message"`
}
