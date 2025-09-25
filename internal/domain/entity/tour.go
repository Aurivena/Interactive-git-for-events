package entity

type TourInput struct {
	DateTour
	Coordinates
	IsTest bool `json:"is_test"`
}
type RouteParams struct {
	DateTour
	PerDayLimit  int      `json:"per_day_limit"`
	Tier         []string `json:"tiers"`
	KindPriority []Kind   `json:"kind_priority"`
	DayStart     string   `json:"day_start"`
	DayEnd       string   `json:"day_end"`
}

type DateTour struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type Coordinates struct {
	Lat *float64 `json:"lat,omitempty" example:"55.434630"`
	Lon *float64 `json:"lon,omitempty" example:"65.353470"`
}

type DayPlan struct {
	Day    string      `json:"day"`
	Places []PlaceInfo `json:"places"`
}

type Tour struct {
	DateTour `json:"date_tour"`
	Days     []DayPlan `json:"placesInfo"`
}

type TourOutput struct {
	ID   UUID `json:"id"`
	Tour Tour `json:"tour"`
}
