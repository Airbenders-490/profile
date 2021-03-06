package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/driftprogramming/pgxpoolmock"
	_ "github.com/jackc/pgx/v4/pgxpool"
)

type reviewRepository struct {
	db pgxpoolmock.PgxPool
}

// NewReviewRepository is the constructor
func NewReviewRepository(db pgxpoolmock.PgxPool) domain.ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

const (
	insertReview       = `INSERT INTO review (id, reviewed, reviewer, created_at) VALUES ($1, $2, $3, $4);`
	joinWithTags       = `INSERT INTO review_tag (review_id, tag_name) VALUES ($1, $2)`
	getReviewForAndBy  = `SELECT * FROM review WHERE reviewed=$1 and reviewer=$2`
	getReviewsFor      = `SELECT * FROM review WHERE reviewed=$1`
	getReviewsBy       = `SELECT * FROM review WHERE reviewer=$1`
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

func (r *reviewRepository) getTagsFor(ctx context.Context, review *domain.Review) error {
	rows, err := r.db.Query(ctx, getTagsFor, review.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.Name)
		if err != nil {
			return errors.NewInternalServerError(err.Error())
		}



		tags = append(tags, &tag)
	}

	review.Tags = tags
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



		err = rows.Scan(&review.ID, &review.Reviewer.ID,&review.Reviewed.ID, &review.CreatedAt)


		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}
	}
	err = r.getTagsFor(ctx, &review)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
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

		var review domain.Review
		err = rows.Scan(&review.ID, &review.Reviewer.ID,&review.Reviewed.ID, &review.CreatedAt)
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		err = r.getTagsFor(ctx, &review)
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *reviewRepository) GetReviewsBy(ctx context.Context, reviewer string) ([]domain.Review, error) {
	rows, err := r.db.Query(ctx, getReviewsBy, reviewer)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var reviews []domain.Review
	for rows.Next() {

		var review domain.Review
		err = rows.Scan(&review.ID, &review.Reviewer.ID,&review.Reviewed.ID, &review.CreatedAt)
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}




		err = r.getTagsFor(ctx, &review)
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}
