package ai

import (
	"arch/internal/domain/entity"
	"context"
	"net/http"
	"time"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Send(c *gin.Context) {
	var input entity.UserSend
	if err := c.ShouldBindJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:    "Bad Request",
				Message:  "Bad Request",
				Solution: "",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	output, err := h.application.SendAi(ctx, input)
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusInternalServerError,
			Detail: envelope.ErrorDetail{
				Title:   "Internal Server Error",
				Message: err.Error(),
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}
