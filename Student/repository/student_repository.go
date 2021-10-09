package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type studentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) domain.StudentRepository {
	return &studentRepository{
		db: db,
	}
}

const (
	insert = `INSERT INTO public.student(
	id, first_name, last_name, email, general_info, school, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	selectByID = `SELECT id, first_name, last_name, email, general_info, school, created_at, updated_at
	FROM public.student WHERE id=$1;`
	update = `UPDATE public.student
	SET first_name=$2, last_name=$3, email=$4, general_info=$5, school=$6, created_at=$7, updated_at=$8
	WHERE id=$1;`
	delete = `DELETE FROM public.student
	WHERE id=$1;`
)

func (r *studentRepository) Create(ctx context.Context, id string, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, insert, id, st.FirstName, st.LastName, st.Email, st.GeneralInfo, st.School, st.CreatedAt, st.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *studentRepository) GetByID(ctx context.Context, id string) (*domain.Student, error) {
	rows, err := r.db.Query(ctx, selectByID, id)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var student domain.Student
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		student.ID = values[0].(string)
		student.FirstName = values[1].(string)
		student.LastName = values[2].(string)
		student.Email = values[3].(string)
		student.GeneralInfo = values[4].(string)
		student.School = values[5].(string)
		student.CreatedAt = values[6].(time.Time)
		student.UpdatedAt = values[7].(time.Time)
	}

	return &student, nil
}

func (r *studentRepository) Update(ctx context.Context, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, update, st.ID, st.FirstName, st.LastName, st.Email, st.GeneralInfo, st.School, st.CreatedAt, st.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *studentRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, delete, id)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
