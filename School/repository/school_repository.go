package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/driftprogramming/pgxpoolmock"
	"log"
)

type schoolRepository struct {
	db pgxpoolmock.PgxPool
}

// NewSchoolRepository returns an instance of school repository
func NewSchoolRepository(db pgxpoolmock.PgxPool) domain.SchoolRepository {
	return &schoolRepository{
		db: db,
	}
}

const (
	findByDomain = `SELECT s.id, s.name, s.country FROM (select id, name, country, unnest(domains) as domain from 
	school) as s WHERE s.domain SIMILAR TO ($1)`
	insertConfirmation = `INSERT INTO public.confirmation(
	token, sc_id, st_id, created_at)
	VALUES ($1, $2, $3, $4);`
	getConfirmationByToken  = `SELECT token, sc_id, st_id, created_at FROM confirmation WHERE token=$1`
	updateStudentWithSchool = `UPDATE public.student
	SET school=$1 WHERE id=$2;`
)

// SearchByDomain finds the schools matching the domain name pattern. Otherwise, returns an empty slice
func (r *schoolRepository) SearchByDomain(ctx context.Context, domainName string) ([]domain.School, error) {
	rows, err := r.db.Query(ctx, findByDomain, domainName)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var schools []domain.School
	for rows.Next() {
		var school domain.School
		err = rows.Scan(&school.ID, &school.Name, &school.Country)
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		schools = append(schools, school)
	}

	return schools, nil
}

// SaveConfirmationToken saves the token which will be used to confirm student's school. Returns nil if can't save
func (r *schoolRepository) SaveConfirmationToken(ctx context.Context, confirmation *domain.Confirmation) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, insertConfirmation, confirmation.Token, confirmation.School.ID, confirmation.Student.ID, confirmation.CreatedAt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("error in Save confirmation repo while saving the token")
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

// GetConfirmationByToken returns a Confirmation with verifiable token and the student and school info.
// Return empty confirmation if no such record found
func (r *schoolRepository) GetConfirmationByToken(ctx context.Context, token string) (*domain.Confirmation, error) {
	rows, err := r.db.Query(ctx, getConfirmationByToken, token)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var confirmation domain.Confirmation
	for rows.Next() {
		err = rows.Scan(&confirmation.Token, &confirmation.School.ID, &confirmation.Student.ID, &confirmation.CreatedAt)
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}
	}

	return &confirmation, nil
}

// AddSchoolForStudent is invoked when the student's school is confirmed. Student table is altered to store the school
// ID now
func (r *schoolRepository) AddSchoolForStudent(ctx context.Context, stID string, scID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, updateStudentWithSchool, scID, stID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
