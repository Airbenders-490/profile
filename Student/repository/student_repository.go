package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/driftprogramming/pgxpoolmock"
	"time"
)

type studentRepository struct {
	db pgxpoolmock.PgxPool
}

// NewStudentRepository is the constructor
func NewStudentRepository(db pgxpoolmock.PgxPool) domain.StudentRepository {
	return &studentRepository{
		db: db,
	}
}

const (
	insert = `INSERT INTO public.student(
	id, first_name, last_name, email, general_info, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`
	selectByID = `SELECT id, first_name, last_name, email, general_info, school, current_classes, classes_taken, created_at, updated_at
	FROM public.student WHERE id=$1;`
	update = `UPDATE public.student
	SET first_name=$2, last_name=$3, email=$4, general_info=$5, created_at=$6, updated_at=$7
	WHERE id=$1;`
	deleteStudent = `DELETE FROM public.student
	WHERE id=$1;`
	getSchoolName = `SELECT name FROM school WHERE ID=$1`
	updateClasses = `UPDATE public.student SET current_classes=$1, classes_taken=$2, updated_at=$3 WHERE id = $4;`
	search = `SELECT id, first_name, last_name, email, general_info, school, current_classes, classes_taken, created_at, updated_at
	FROM public.student WHERE first_name LIKE $1 and last_name like $2 and current_classes @> $3`
)

// Create stores the student in the db. Returns err if unable to
func (r *studentRepository) Create(ctx context.Context, id string, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, insert, id, st.FirstName, st.LastName, st.Email, st.GeneralInfo, st.CreatedAt, st.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

// GetByID returns either an empty student or valid student if it exsits. Returns nil and error if there's some db error
func (r *studentRepository) GetByID(ctx context.Context, id string) (*domain.Student, error) {
	rows, err := r.db.Query(ctx, selectByID, id)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var student domain.Student
	for rows.Next() {
		var schoolID *string
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.GeneralInfo,
			&schoolID, &student.CurrentClasses, &student.ClassesTaken, &student.CreatedAt, &student.UpdatedAt)
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}
		if schoolID != nil {
			student.School = &domain.School{
				ID: *schoolID,
			}
		}
	}

	return &student, nil
}

// Update changes the record in the db. Returns err if isn't able to
func (r *studentRepository) Update(ctx context.Context, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, update, st.ID, st.FirstName, st.LastName, st.Email, st.GeneralInfo, st.CreatedAt, st.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

// Delete deletes the record in the db. Returns err if isn't able to
func (r *studentRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteStudent, id)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *studentRepository) UpdateClasses(ctx context.Context, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, updateClasses, st.CurrentClasses, st.ClassesTaken, time.Now(), st.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}
func (r *studentRepository) SearchStudents(ctx context.Context, st *domain.Student) ([]domain.Student, error) {
	rows, err := r.db.Query(ctx, search, st.FirstName + "%", st.LastName + "%", st.CurrentClasses)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var students []domain.Student
	for rows.Next() {
		var student domain.Student
		var school domain.School
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.GeneralInfo,
			&school.ID, &student.CurrentClasses, &student.ClassesTaken, &student.CreatedAt, &student.UpdatedAt)
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}
		if school.ID != "" {
			student.School = &school
		}
		students = append(students, student)
	}

	return students, nil
}