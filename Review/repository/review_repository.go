package repository

import (
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type reviewRepository struct {
	db *pgxpool.Pool
}

// NewReviewRepository is the constructor
func NewReviewRepository(db *pgxpool.Pool) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

const (
	insertReview      = `INSERT INTO review (id, reviewed, reviewer, created_at) VALUES ($1, $2, $3, $4);`
	joinWithTags      = `INSERT INTO review_tag (review_id, tag_name) VALUES ($1, $2)`
	getReviewForAndBy = `SELECT * FROM review WHERE reviewed=$1 and reviewer=$2`
)

// AddReview adds the review to the review table as well as joins the tags
func (r *reviewRepository) AddReview(ctx context.Context, review *domain.Review) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, insertReview, review.ID, review.Reviewed.ID, review.Reviewer.ID, review.CreatedAt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	for _, tag := range review.Tags {
		txTag, err := r.db.Begin(ctx)
		if err != nil {
			return errors.NewInternalServerError(err.Error())
		}
		_, err = txTag.Exec(ctx, joinWithTags, review.ID, tag.Name)
		if err != nil {
			return errors.NewInternalServerError(err.Error())
		}
		err = txTag.Commit(ctx)
		if err != nil {
			return errors.NewInternalServerError(err.Error())
		}
		_ = txTag.Rollback(ctx)
	}

	return nil
}

// GetReviewByAndFor returns a review if it exists. Otherwise, returns nil. CAN RETURN nil, nil if no error and no
// review is found!
func (r *reviewRepository) GetReviewByAndFor(ctx context.Context, reviewer string, reviewed string) (*domain.Review, error) {
	rows, err := r.db.Query(ctx, getReviewForAndBy, reviewed, reviewer)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var review domain.Review
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}
		fmt.Println(values)
		review.ID = values[0].(string)
		review.Reviewed.ID = values[1].(string)
		review.Reviewer.ID = values[2].(string)
		review.CreatedAt = values[3].(time.Time)
	}

	log.Println(review, err)
	return &review, nil
}

// EditReview only changes the tags in a review. Nothing else
func (r *reviewRepository) EditReview(ctx context.Context, review *domain.Review) error {
	panic("implement me")
}
