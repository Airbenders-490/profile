package http

import (
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/airbenders/profile/utils/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

// SchoolHandler struct
type SchoolHandler struct {
	u domain.SchoolUseCase
}

// NewSchoolHandler returns a new SchoolHandler
func NewSchoolHandler(u domain.SchoolUseCase) *SchoolHandler {
	return &SchoolHandler{u}
}

// SearchStudentSchool is an endpoint that returns schools that matches the name
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

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func extractSchool(email string) (bool, string) {
	if !isEmailValid(email) {
		return false, ""
	}
	return true, strings.Split(email, "@")[1]
}

// SendConfirmationMail sends email to the client for school confirmation
func (h *SchoolHandler) SendConfirmationMail(c *gin.Context) {
	ctx := c.Request.Context()
	key, _ := c.Get("loggedID")
	loggedID, _ := key.(string)

	email := c.Query("email")
	ok, domainName := extractSchool(email)
	if !ok {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("please provide a valid email"))
		return
	}

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

	if len(schools) == 0 {
		c.JSON(http.StatusBadRequest, errors.NewNotFoundError("no school matching the domain found"))
		return
	}

	school := schools[0]
	err = h.u.SendConfirmation(ctx, &domain.Student{ID: loggedID}, email, &school)
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

	c.JSON(200, httputils.NewResponse("email sent"))
}

// ConfirmSchoolRegistration is an internal endpoint (not accessible from the app) that is embedded in the email
func (h *SchoolHandler) ConfirmSchoolRegistration(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("must provide a valid email"))
		return
	}

	err := h.u.ConfirmSchoolEnrollment(ctx, token)
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

	c.JSON(http.StatusOK, httputils.NewResponse("school confirmed"))
}
