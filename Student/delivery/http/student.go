package http

import (
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/airbenders/profile/utils/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StudentHandler struct {
	UseCase domain.StudentUseCase
}

func NewStudentHandler(u domain.StudentUseCase) *StudentHandler {
	return &StudentHandler{UseCase: u}
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id must be provided"))
		return
	}

	ctx := c.Request.Context()
	student, err := h.UseCase.GetByID(ctx, id)
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
	c.JSON(200, student)
}

func (h *StudentHandler) Create(c *gin.Context) {
	var student domain.Student
	err := c.ShouldBindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid data"))
		return
	}

	ctx := c.Request.Context()
	err = h.UseCase.Create(ctx, &student)
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

	c.JSON(http.StatusCreated, httputils.NewResponse("student created"))
}

func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id must be provided"))
		return
	}

	var student domain.Student
	err := c.ShouldBindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid data"))
		return
	}

	ctx := c.Request.Context()
	err = h.UseCase.Update(ctx, id, &student)
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

	c.JSON(http.StatusOK, httputils.NewResponse("student updated"))
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id must be provided"))
		return
	}

	ctx := c.Request.Context()
	err := h.UseCase.Delete(ctx, id)
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

	c.JSON(http.StatusOK, httputils.NewResponse("student deleted"))
}
