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

	loggedID, _ := c.Get("loggedID")
	reviewer, _ := loggedID.(string)

	ctx := c.Request.Context()
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

// EditReview alters the tags for the review
func (h *ReviewHandler) EditReview(c *gin.Context) {
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

	loggedID, _ := c.Get("loggedID")
	reviewer, _ := loggedID.(string)

	updatedReview, err := h.u.EditReview(ctx, &review, reviewer)
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

	c.JSON(http.StatusCreated, updatedReview)
}

// GetReviewsBy returns the reviews made by that student
func (h *ReviewHandler) GetReviewsBy(c *gin.Context) {
	reviewer := c.Param("reviewer")
	if reviewer == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("must provide reviewer id"))
		return
	}

	loggedID, _ := c.Get("loggedID")
	logged, _ := loggedID.(string)
	if logged != reviewer {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not allowed to get reviews by others"))
	}

	if loggedID != reviewer {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized to edit this review"))
	}
	ctx := c.Request.Context()
	reviews, err := h.u.GetReviewsBy(ctx, reviewer)
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

	c.JSON(http.StatusOK, reviews)
}
