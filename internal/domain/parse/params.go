package parse

import (
	"arch/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Parse(query *entity.Query, c *gin.Context) error {
	page := c.Query("page")
	limit := c.DefaultQuery("limit", "10")

	var err error
	query.Page, err = strconv.Atoi(page)
	if err != nil {
		logrus.Error(err)
		return err
	}
	query.Limit, err = strconv.Atoi(limit)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
