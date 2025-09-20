package place

import (
	"arch/internal/domain/entity"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

// List
// @Tags        Places
// @Summary     Получить список мест
// @Description Возвращает полный список мест.
// @Produce     json
// @Success     200 {array}  []entity.PlaceInfo         "Список мест"
// @Failure     500 {object} entity.AppError      "Внутренняя ошибка сервера (Spond error)"
// @Router      /place/list [get]
func (h *Handler) List(c *gin.Context) {
	output, err := h.application.List()
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusInternalServerError,
			Detail: envelope.ErrorDetail{
				Title:   "Не удалось получить список мест",
				Message: "Сервис мест вернул ошибку: " + err.Error(),
				Solution: "1) Повторите запрос через 1–2 минуты\n" +
					"2) Если не помогает — обратитесь в поддержку и укажите код: PLACE_LIST_FAILED",
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}

// ByID
// @Tags        Places
// @Summary     Получить место по ID
// @Description Возвращает место по идентификатору.
// @Produce     json
// @Param       id   path   string true "ID места (UUID, формат: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"
// @Success     200  {object} entity.PlaceInfo        "Место найдено"
// @Failure     400  {object} entity.AppError     "Некорректный формат ID (Spond error)"
// @Failure     404  {object} entity.AppError     "Место не найдено (Spond error)"
// @Failure     500  {object} entity.AppError     "Внутренняя ошибка сервера (Spond error)"
// @Router      /place/{id} [get]
func (h *Handler) ByID(c *gin.Context) {
	id := entity.UUID(c.Param("id"))
	if ok := id.Valid(); !ok {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:   "Некорректный идентификатор",
				Message: "ID должен быть в формате UUID (RFC 4122)",
				Solution: "Проверьте формат: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\n" +
					"Пример: 11111111-1111-1111-1111-111111111111",
			},
		})
		return
	}

	output, err := h.application.ByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.spond.SendResponseError(c.Writer, &envelope.AppError{
				Code: http.StatusNotFound,
				Detail: envelope.ErrorDetail{
					Title:   "Место не найдено",
					Message: "Место с указанным ID отсутствует",
					Solution: "1) Проверьте корректность ID\n" +
						"2) Получите список доступных мест через /place и выберите существующий ID",
				},
			})
			return
		}

		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusInternalServerError,
			Detail: envelope.ErrorDetail{
				Title:   "Ошибка при получении места",
				Message: "Внутренняя ошибка сервера: " + err.Error(),
				Solution: "1) Повторите запрос позже\n" +
					"2) Обратитесь в поддержку и укажите код: PLACE_GET_BY_ID_FAILED",
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}

// ListByKind
// @Tags        Places
// @Summary     Получить места по типу (kind)
// @Description Возвращает список мест, отфильтрованный по типу.
// @Produce     json
// @Param       kind path string true "Тип места (enum)" Enums(cinema,theatre,concert_hall,stadium,sport,museum,art_gallery,historic,memorial,park,zoo,aquapark,attraction,church,monastery,mosque,synagogue,mall,market,monument,restaurant)
// @Success     200  {array}  []entity.PlaceInfo        "Список мест по типу"
// @Failure     400  {object} entity.AppError     "Некорректное значение kind (Spond error)"
// @Failure     500  {object} entity.AppError     "Внутренняя ошибка сервера (Spond error)"
// @Router      /place/list/{kind} [get]
func (h *Handler) ListByKind(c *gin.Context) {
	kind := entity.Kind(c.Param("kind"))
	if ok := kind.Valid(); !ok {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusBadRequest,
			Detail: envelope.ErrorDetail{
				Title:   "Некорректное значение kind",
				Message: "Передан тип, отсутствующий в справочнике",
				Solution: "Используйте одно из значений enum (см. /docs): " +
					"cinema, theatre, concert_hall, stadium, sport, museum, art_gallery, historic, memorial, park, zoo, aquapark, attraction, church, monastery, mosque, synagogue, mall, market, monument, restaurant",
			},
		})
		return
	}

	output, err := h.application.ListByKind(kind)
	if err != nil {
		h.spond.SendResponseError(c.Writer, &envelope.AppError{
			Code: http.StatusInternalServerError,
			Detail: envelope.ErrorDetail{
				Title:   "Ошибка при получении списка по типу",
				Message: "Внутренняя ошибка сервера: " + err.Error(),
				Solution: "1) Повторите запрос позже\n" +
					"2) Обратитесь в поддержку и укажите код: PLACE_LIST_BY_KIND_FAILED",
			},
		})
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}
