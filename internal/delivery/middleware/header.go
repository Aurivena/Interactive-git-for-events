package middleware

import (
	"net/http"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) Session(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		m.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:    "Не указан идентификатор сессии",
				Message:  "Заголовок X-Session-ID обязателен для получения истории.",
				Solution: "Добавьте X-Session-ID в заголовок запроса и повторите попытку.",
			},
		})
		return
	}
}
