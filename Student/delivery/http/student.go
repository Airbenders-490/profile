package http

import (
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/airbenders/profile/utils/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StudentHandler struct
type StudentHandler struct {
	UseCase domain.StudentUseCase
}

// NewStudentHandler is the constructor
func NewStudentHandler(u domain.StudentUseCase) *StudentHandler {
	return &StudentHandler{UseCase: u}
}

const errorMessage = "id must be provided"

// GetByID returns the student's profile with that ID. If it doesn't exist, returns 404
func (h *StudentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {

		c.JSON(http.StatusBadRequest, errors.NewBadRequestError(errorMessage))
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

// Create is hit when the student first creates his account and is asked to set it up.
func (h *StudentHandler) Create(c *gin.Context) {
	key, _ := c.Get("loggedID")
	loggedID, _ := key.(string)

	var student domain.Student
	err := c.ShouldBindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid data"))
		return
	}
	student.ID = loggedID

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

// Update changes the student record. Ensures the student is the same as logged in, and then makes changes as requested
func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError(errorMessage))
		return
	}

	key, _ := c.Get("loggedID")
	loggedID, _ := key.(string)

	if loggedID != id {
		c.JSON(http.StatusBadRequest, errors.NewUnauthorizedError("Can only update for self"))
		return
	}

	var student domain.Student
	err := c.ShouldBindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid data"))
		return
	}

	ctx := c.Request.Context()
	updatedStudent, err := h.UseCase.Update(ctx, id, &student)
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

	c.JSON(http.StatusOK, updatedStudent)
}

// Delete simply deletes the profile as requested
func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError(errorMessage))
		return
	}

	key, _ := c.Get("loggedID")
	loggedID, _ := key.(string)

	if loggedID != id {
		c.JSON(http.StatusBadRequest, errors.NewUnauthorizedError("Can only delete for self"))
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
