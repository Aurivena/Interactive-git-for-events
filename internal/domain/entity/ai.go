package entity

type UserSend struct {
	Message string   `json:"message"`
	Lat     *float64 `json:"lat,omitempty"`
	Lon     *float64 `json:"lon,omitempty"`
}

type PlaceInfo struct {
	ID          string  `json:"id" db:"id"`
	Title       string  `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	Address     string  `json:"address" db:"address"`
	Lon         float64 `json:"lon" db:"lon"`
	Lat         float64 `json:"lat" db:"lat"`
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
	PlaceInfo []PlaceInfo
	Message   string `json:"message"`
}
