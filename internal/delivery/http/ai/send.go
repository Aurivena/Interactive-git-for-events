package ai

import (
	"arch/internal/domain/entity"
	"net/http"

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
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:   "Ошибка при запросе",
				Message: "Не удалось обработать ваш запрос",
				Solution: "1. Перепроверьте веденные вами данные.\n" +
					"2. Обратитесь к администратору, если не смогли решить проблема.",
			},
		})
		return
	}
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:    "Не указан идентификатор сессии",
				Message:  "Заголовок X-Session-ID обязателен для получения истории.",
				Solution: "Добавьте X-Session-ID в заголовок запроса и повторите попытку.",
			},
		})
		return
	}
	output, err := h.application.SendAi(input, sessionID)
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusInternalServerError,
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
