package entity

type TimeOnly struct {
	Hour   int
	Minute int
	Second int
}

func (t *TimeOnly) Valid() bool {
	return t.Hour >= 0 && t.Hour < 24 &&
		t.Minute >= 0 && t.Minute < 60 &&
		t.Second >= 0 && t.Second < 60
}

type Weekday string

const (
	monday    Weekday = "monday"
	tuesday   Weekday = "tuesday"
	wednesday Weekday = "wednesday"
	thursday  Weekday = "thursday"
	friday    Weekday = "friday"
	saturday  Weekday = "saturday"
	sunday    Weekday = "sunday"
)

func (w *Weekday) Valid() bool {
	switch *w {
	case monday, tuesday, wednesday, thursday, friday, saturday, sunday:
		return true
	default:
		return false
	}
}

func (w *Weekday) Convert() string {
	return string(*w)
}
