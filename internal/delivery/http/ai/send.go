package ai

import (
	"arch/internal/delivery/middleware"
	"arch/internal/domain/entity"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

// Send
// @Tags         AI
// @Summary      Отправить запрос на полуение списка мест
// @Description  Обрабатывает запрос и на основе него выдает список мест, которые можно посетить.
// @Accept json
// @Produce json
// @Param input body entity.UserSend true "Данные для создания генерации ответа"
// @Success      200           {object} []entity.ChatOutput  "Ответ сгенерирован"
// @Failure      400           {object} entity.AppErrorDoc         "Некорректные данные (Spond error)"
// @Failure      500           {object} entity.AppErrorDoc         "Внутренняя ошибка сервера (Spond error)"
// @Router       /ai/send [post]
func (h *Handler) Send(c *gin.Context) {
	var input entity.UserSend
	if err := c.ShouldBindJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: envelope.BadRequest,
			Detail: envelope.ErrorDetail{
				Title:   "Ошибка при запросе",
				Message: "Не удалось обработать ваш запрос",
				Solution: "1. Перепроверьте веденные вами данные.\n" +
					"2. Обратитесь к администратору, если не смогли решить проблема.",
			},
		})
		return
	}
	output, err := h.application.SendAi(input, c.GetHeader(middleware.Session))
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: envelope.InternalServerError,
			Detail: envelope.ErrorDetail{
				Title:    "Ошибка сервера",
				Message:  "Мы 100% устраняем эту ошибку!",
				Solution: "1. Пожалуйста подождите и сообщите в тех-поддержку!!!",
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}

// GenerateTour
// @Summary     Сгенерировать маршрут тура
// @Description Формирует новый тур на основе введённых параметров (даты, координаты, типы мест, лимиты).
// @Tags        Tour
// @Accept      json
// @Produce     json
// @Param       input  body      entity.TourInput  true  "Параметры генерации тура"
// @Success     200    {object}  entity.Tour "Сгенерированный тур"
// @Failure     400    {object}  entity.AppErrorDoc "Ошибка валидации входных данных"
// @Failure     500    {object}  entity.AppErrorDoc"Внутренняя ошибка сервера"
// @Router      /ai/generate/tour [post]
func (h *Handler) GenerateTour(c *gin.Context) {
	var input entity.TourInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: envelope.BadRequest,
			Detail: envelope.ErrorDetail{
				Title:   "Ошибка при запросе",
				Message: "Не удалось обработать ваш запрос",
				Solution: "1. Перепроверьте веденные вами данные.\n" +
					"2. Обратитесь к администратору, если не смогли решить проблема.",
			},
		})
		return
	}

	output, err := h.application.GenerateTour(&input, c.GetHeader(middleware.Session))
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: envelope.InternalServerError,
			Detail: envelope.ErrorDetail{
				Title:    "Ошибка сервера",
				Message:  "Мы 100% устраняем эту ошибку!",
				Solution: "1. Пожалуйста подождите и сообщите в тех-поддержку!!!",
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}
