package tour

import (
	"arch/internal/delivery/middleware"
	"arch/internal/domain/entity"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

// All
// @Summary     Получить все туры пользователя
// @Description Возвращает список всех туров, сохранённых за конкретную сессию (X-Session-ID).
// @Tags        Tour
// @Produce     json
// @Success     200  {array}   entity.Tour  "Список туров"
// @Failure     500  {object}  entity.AppErrorDoc  "Внутренняя ошибка сервера"
// @Router      /tours [get]
// @Param       X-Session-ID  header  string  true  "Идентификатор сессии"
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

// ByID
// @Summary     Получить тур по ID
// @Description Возвращает тур по его уникальному идентификатору.
// @Tags        Tour
// @Produce     json
// @Param       id   path     string           true  "UUID тура"
// @Success     200  {object} entity.Tour "Найденный тур"
// @Failure     400  {object} entity.AppErrorDoc "Некорректный ID"
// @Failure     500  {object} entity.AppErrorDoc "Внутренняя ошибка сервера"
// @Router      /tours/{id} [get]
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
