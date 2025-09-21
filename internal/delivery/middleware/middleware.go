package middleware

import "github.com/Aurivena/spond/v2/core"

type Middleware struct {
	spond *core.Spond
}

const (
	Session = "X-Session-ID"
)

func New(spond *core.Spond) *Middleware {
	return &Middleware{
		spond: spond,
	}
}
