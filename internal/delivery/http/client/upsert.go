package client

import (
	"arch/internal/delivery/middleware"
	"arch/internal/domain/entity"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Upsert
// @Tags        Client
// @Summary     Сохранить/обновить анкету клиента
// @Description Принимает ответы анкеты (survey) и сохраняет их в профиле клиента.
// @Accept      json
// @Produce     json
// @Param       X-Session-ID header string true "ID сессии клиента"
// @Param       input body entity.Survey true "Ответы анкеты"
// @Success     204 {object} nil "Успешное сохранение"
// @Failure     400 {object} entity.AppErrorDoc "Некорректный запрос (невалидный JSON)"
// @Failure     500 {object} entity.AppErrorDoc "Внутренняя ошибка сервера"
// @Router      /client/upsert [post]
func (h *Handler) Upsert(c *gin.Context) {
	var input entity.Survey
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		logrus.Error(err)
		return
	}

	if err := h.application.UpsertClientSurvey(c.GetHeader(middleware.Session), input); err != nil {
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

	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}
