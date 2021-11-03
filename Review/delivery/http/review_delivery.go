package http

import (
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ReviewHandler struct
type ReviewHandler struct {
	u domain.ReviewUseCase
}

// NewReviewHandler is the constructor
func NewReviewHandler(ru domain.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{u: ru}
}

// AddReview binds review body and forwards it to the useCase for processing
func (h *ReviewHandler) AddReview(c *gin.Context) {
	reviewed := c.Param("reviewed")
	if reviewed == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("must provide reviewed id"))
		return
	}

	var review domain.Review
	err := c.ShouldBindJSON(&review)
	if err != nil || review.Reviewed.ID == "" || review.Tags == nil {
		fmt.Println(review, err.Error())
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid review body"))
		return
	}

	ctx := c.Request.Context()
	// todo: get this from auth!
	reviewer := "asda"
	createdReview, err := h.u.AddReview(ctx, &review, reviewer)
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

	c.JSON(http.StatusCreated, createdReview)
}
