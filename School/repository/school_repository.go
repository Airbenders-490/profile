package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type schoolRepository struct {
	db *pgxpool.Pool
}

func NewSchoolRepository(db *pgxpool.Pool) domain.SchoolRepository {
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
	getConfirmationByToken = `SELECT token, sc_id, st_id, created_at FROM confirmation WHERE token=$1`
	updateStudentWithSchool = `UPDATE public.student
	SET school=$1 WHERE id=$2;`
)

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
		values, err := rows.Values()
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		school.ID = values[0].(string)
		school.Name = values[1].(string)
		school.Country = values[2].(string)

		schools = append(schools, school)
	}

	return schools, nil
}


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
	 	return errors.NewInternalServerError(err.Error())
	 }

	 return nil
}

func (r *schoolRepository) GetConfirmationByToken(ctx context.Context, token string) (*domain.Confirmation, error) {
	rows, err := r.db.Query(ctx, getConfirmationByToken, token)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var confirmation domain.Confirmation
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		confirmation.Token = values[0].(string)
		confirmation.School.ID = values[1].(string)
		confirmation.Student.ID = values[2].(string)
		confirmation.CreatedAt = values[3].(time.Time)
	}

	return &confirmation, nil
}

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