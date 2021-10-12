package http

import (
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SchoolHandler struct {
	u domain.SchoolUseCase
}

func NewSchoolHandler(u domain.SchoolUseCase) *SchoolHandler {
	return &SchoolHandler{u}
}

func (h *SchoolHandler) SearchStudentSchool(c *gin.Context) {
	domainName, ok := c.GetQuery("domain")
	if !ok || domainName == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("please provide a query"))
		return
	}

	ctx := c.Request.Context()
	schools, err := h.u.SearchSchoolByDomain(ctx, domainName)
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

	c.JSON(http.StatusOK, schools)
}
