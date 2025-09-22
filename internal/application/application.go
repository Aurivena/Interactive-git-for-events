package application

import (
	"arch/internal/domain/entity"
	"arch/internal/infrastructure"
)

type Application struct {
	aiConfig *entity.AiConfig
	post     *infrastructure.Infrastructure
}

func New(post *infrastructure.Infrastructure, aiConfig *entity.AiConfig) *Application {
	return &Application{
		aiConfig: aiConfig,
		post:     post,
	}
}
