package application

import (
	"arch/internal/domain/entity"
	"arch/internal/infrastructure"
	"net/http"
)

type Application struct {
	qwqConfig   *entity.AiConfig
	proxyClient *http.Client
	post        *infrastructure.Infrastructure
}

func New(post *infrastructure.Infrastructure, qwqConfig *entity.AiConfig, proxyClient *http.Client) *Application {
	return &Application{
		qwqConfig:   qwqConfig,
		proxyClient: proxyClient,
		post:        post,
	}
}
