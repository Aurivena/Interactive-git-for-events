package place

import (
	"arch/internal/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func isUniqueViolation(err error, constraints ...string) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())

	// Ловим любые варианты unique-ошибок
	if !(strings.Contains(s, "23505") ||
		strings.Contains(s, "duplicate key value violates unique constraint") ||
		strings.Contains(s, "unique constraint")) {
		return false
	}

	// Опционально — проверим имя констрейнта, если передали
	if len(constraints) > 0 {
		for _, c := range constraints {
			if c == "" {
				continue
			}
			if strings.Contains(s, strings.ToLower(c)) {
				return true
			}
		}
		// 23505 есть, но не наш констрейнт — считаем не нашим кейсом
		return false
	}

	// 23505 без уточнений — ок, это дубль
	return true
}

func (r *Place) Write(id uuid.UUID, rawSQL string) error {
	_, err := r.db.Exec(rawSQL, id)
	if err != nil {
		// Хотим именно ваш индекс — добавьте сюда имя, если нужно:
		// if isUniqueViolation(err, "uniq_place_title_address") { ... }
		if isUniqueViolation(err) {
			return domain.FileDuplicate
		}
		logrus.Error(err)
		return err
	}
	return nil
}
