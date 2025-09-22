package entity

type AppErrorDoc struct {
	Code   int         `json:"code"`
	Detail ErrorDetail `json:"detail"`
}

type ErrorDetail struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Solution string `json:"solution"`
}

// --- Вложенные типы для tags ---

type ScheduleItem struct {
	End           string `json:"end" example:"23:00"`
	Week          string `json:"week" example:"monday"`
	Start         string `json:"start" example:"09:00"`
	SpansMidnight bool   `json:"spans_midnight" example:"false"`
}

type PlaceTags struct {
	Phone    string         `json:"phone"   example:"+7 (3522) 60-70-55"`
	Website  string         `json:"website" example:"https://cinema.pushka.club/kurgan/pushka"`
	Schedule []ScheduleItem `json:"schedule"`
}

// --- Док-версия PlaceInfo ---

// @name PlaceInfoDoc
type PlaceInfoDoc struct {
	ID          string    `json:"id"          example:"0199574c-c996-7301-95af-4b76e8b6088a"`
	Title       string    `json:"title"       example:"Pushka"`
	Kind        string    `json:"kind" example:"cinema"`
	Description *string   `json:"description" example:"Современный кинотеатр в центре города"`
	Address     string    `json:"address"     example:"Курган, ул. Пушкина, 25, ТРЦ «Пушкинский», 3 этаж"`
	Lon         float64   `json:"lon"         example:"65.318954"`
	Lat         float64   `json:"lat"         example:"55.432190"`
	Tags        PlaceTags `json:"tags"`
}

type ChatOutputDoc struct {
	PlaceInfo []PlaceInfoDoc `json:"placesInfo"`
	Message   string         `json:"message"`
}

type HistoryDoc struct {
	ID        string          `json:"id" db:"id"`
	Message   string          `json:"message" db:"message"`
	AiMessage []ChatOutputDoc `json:"ai_message" db:"ai_message"`
}
