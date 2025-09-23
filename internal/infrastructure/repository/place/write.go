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

	if !(strings.Contains(s, "23505") ||
		strings.Contains(s, "duplicate key value violates unique constraint") ||
		strings.Contains(s, "unique constraint")) {
		return false
	}

	if len(constraints) > 0 {
		for _, c := range constraints {
			if c == "" {
				continue
			}
			if strings.Contains(s, strings.ToLower(c)) {
				return true
			}
		}
		return false
	}
	return true
}

func (r *Place) Write(id uuid.UUID, rawSQL string) error {
	_, err := r.db.Exec(rawSQL, id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.FileDuplicate
		}
		logrus.Error(err)
		return err
	}
	return nil
}
