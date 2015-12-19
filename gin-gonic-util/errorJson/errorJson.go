package errorJson

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsaikd/KDGoLib/errutil"
)

// RenderErrorJSON is a gin middleware to render error in json format
func RenderErrorJSON(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		errs := []error{}
		for _, err := range c.Errors {
			errs = append(errs, err)
		}

		c.JSON(
			http.StatusNotFound,
			errutil.NewErrorSlice(errs...),
		)
	}
}
