package tour

import (
	"arch/internal/delivery/middleware"
	"arch/internal/domain/entity"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

func (h *Handler) All(c *gin.Context) {
	output, err := h.application.TourAll(c.GetHeader(middleware.Session))
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

func (h *Handler) ByID(c *gin.Context) {
	id := entity.UUID(c.Param("id"))
	if ok := id.Valid(); !ok {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: envelope.BadRequest,
			Detail: envelope.ErrorDetail{
				Title:   "Некорректный идентификатор",
				Message: "ID должен быть в формате UUID (RFC 4122)",
				Solution: "Проверьте формат: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\n" +
					"Пример: 11111111-1111-1111-1111-111111111111",
			},
		})
		return
	}
	output, err := h.application.TourByID(id)
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
