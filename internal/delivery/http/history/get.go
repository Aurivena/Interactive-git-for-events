package history

import (
	"arch/internal/delivery/middleware"
	"arch/internal/domain/entity"
	"arch/internal/domain/parse"
	"net/http"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

// ListHistory
// @Tags         History
// @Summary      Список сообщений сессии
// @Description  Возвращает историю чата по X-Session-ID с пагинацией (page — с 1).
// @Accept       json
// @Produce      json
// @Param        X-Session-ID  header  string  true  "ID сессии"
// @Param        page          query   int     false "Номер страницы (>=1)" minimum(1) default(1)
// @Param        limit         query   int     false "Размер страницы (1..100)" minimum(1) maximum(100) default(20)
// @Success      200  {array}   entity.HistoryDoc            "История сообщений"
// @Failure      400  {object}  entity.AppErrorDoc        "Ошибочные параметры/заголовки"
// @Failure      404  {object}  entity.AppErrorDoc        "Не найдено"
// @Failure      500  {object}  entity.AppErrorDoc        "Внутренняя ошибка сервера"
// @Router       /history [get]
func (h *Handler) ListHistory(c *gin.Context) {
	var query entity.Query
	if err := parse.Parse(&query, c); err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:    "Параметры запроса некорректны",
				Message:  "Не удалось разобрать query params (page, limit).",
				Solution: "Передайте целые числа: page >= 1, limit >= 10",
			},
		})
		return
	}

	output, err := h.application.ListHistory(&query, c.GetHeader(middleware.Session))
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusInternalServerError,
			Detail: envelope.ErrorDetail{
				Title:    "Не удалось получить историю",
				Message:  "Произошла ошибка при обращении к хранилищу истории.",
				Solution: "Попробуйте позже. Если ошибка повторяется — свяжитесь с поддержкой.",
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}
