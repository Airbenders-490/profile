package http

import (
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TagHandler struct
type TagHandler struct {
	u domain.TagUseCase
}

// NewTagHandler is a constructor
func NewTagHandler(u domain.TagUseCase) *TagHandler {
	return &TagHandler{u: u}
}

// GetAllTags returns all the tags currently available
func (h *TagHandler) GetAllTags(c *gin.Context) {
	ctx := c.Request.Context()

	tags, err := h.u.GetAllTags(ctx)
	if err != nil {
		switch v := err.(type) {
		case *errors.RestError:
			c.JSON(v.Code, v)
			return
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError(err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, tags)
}
