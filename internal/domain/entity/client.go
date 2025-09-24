package entity

type Survey struct {
	ComfortService int `json:"comfortService"`
	Culture        int `json:"culture"`
	Sport          int `json:"sport"`
	History        int `json:"history"`
	FreshAir       int `json:"freshAir"`
	Religion       int `json:"religion"`
	Shop           int `json:"shop"`
	Restaurant     int `json:"restaurant"`
}
