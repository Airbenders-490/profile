package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/jackc/pgx/v4/pgxpool"
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
	insertReview       = `INSERT INTO review (id, reviewed, reviewer, created_at) VALUES ($1, $2, $3, $4);`
	joinWithTags       = `INSERT INTO review_tag (review_id, tag_name) VALUES ($1, $2)`
	getReviewForAndBy  = `SELECT * FROM review WHERE reviewed=$1 and reviewer=$2`
	getReviewsFor      = `SELECT * FROM review WHERE reviewed=$1`
	deleteExistingTags = `DELETE FROM review_tag WHERE review_id=$1`
	getTagsFor         = `SELECT tag_name FROM review_tag WHERE review_id=$1`
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

	err = r.addTags(ctx, review)

	return err
}

func (r *reviewRepository) addTags(ctx context.Context, review *domain.Review) error {
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

func (r *reviewRepository) getTagsFor(ctx context.Context, review *domain.Review) ([]domain.Tag, error) {
	rows, err := r.db.Query(ctx, getTagsFor, review.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var tags []domain.Tag
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		var tag domain.Tag
		tag.Name = values[0].(string)

		tags = append(tags, tag)
	}

	return tags, nil
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

		review.ID = values[0].(string)
		review.Reviewed.ID = values[1].(string)
		review.Reviewer.ID = values[2].(string)
		review.CreatedAt = values[3].(time.Time)
	}

	return &review, nil
}

// UpdateReviewTags only changes the tags in a review. Nothing else
func (r *reviewRepository) UpdateReviewTags(ctx context.Context, review *domain.Review) error {
	// delete existing tags
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	_, err = tx.Exec(ctx, deleteExistingTags, review.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	// add the new ones
	err = r.addTags(ctx, review)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *reviewRepository) GetReviewsFor(ctx context.Context, reviewed string) ([]domain.Review, error) {
	rows, err := r.db.Query(ctx, getReviewsFor, reviewed)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var reviews []domain.Review
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		var review domain.Review
		review.ID = value[0].(string)
		review.Reviewed.ID = value[1].(string)
		review.Reviewer.ID = value[2].(string)
		review.CreatedAt = value[3].(time.Time)

		tags, err := r.getTagsFor(ctx, &review)
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}
		review.Tags = tags

		reviews = append(reviews, review)
	}

	return reviews, nil
}
