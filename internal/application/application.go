package application

import (
	"arch/internal/domain/entity"
	"arch/internal/infrastructure"
)

type Application struct {
	qwqConfig *entity.AiConfig
	post      *infrastructure.Infrastructure
}

func New(post *infrastructure.Infrastructure, qwqConfig *entity.AiConfig) *Application {
	return &Application{
		qwqConfig: qwqConfig,
		post:      post,
	}
}
