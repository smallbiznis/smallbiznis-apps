package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/smallbiznis/smallbiznis-apps/pkg/errutil"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Errors.Last()

		if v, ok := err.Err.(errutil.BaseError); ok {
			c.JSON(v.Code.HTTPStatus(), v)
			return
		}
	}
}
