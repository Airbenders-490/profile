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
	selectByID = `SELECT id, first_name, last_name, email, general_info, school, created_at, updated_at
	FROM public.student WHERE id=$1;`
	update = `UPDATE public.student
	SET first_name=$2, last_name=$3, email=$4, general_info=$5, created_at=$6, updated_at=$7
	WHERE id=$1;`
	deleteStudent = `DELETE FROM public.student
	WHERE id=$1;`
	getSchoolName = `SELECT name FROM school WHERE ID=$1`
	//addEnrolledClass = `UPDATE public.student SET current_classes=array_append(current_classes, $1), updated_at=$2 WHERE id=$3;`
	updateCurrentClasses = `UPDATE public.student SET current_classes=$1, updated_at=$2 WHERE id = $3;`
	updateClassesTaken = `UPDATE public.student SET classes_taken=$1, updated_at=$2 WHERE id=$3;`
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
		// todo: remove dead code after confirming new one works
		//values, err := rows.Values()
		//if err != nil {
		//	err = errors.NewInternalServerError(err.Error())
		//	return nil, err
		//}
		//
		//student.ID = values[0].(string)
		//student.FirstName = values[1].(string)
		//student.LastName = values[2].(string)
		//student.Email = values[3].(string)
		//student.GeneralInfo = values[4].(string)
		//if values[5] != nil {
		//	student.School = &domain.School{ID: values[5].(string)}
		//	row := r.db.QueryRow(ctx, getSchoolName, student.School.ID)
		//	var name string
		//	err = row.Scan(&name)
		//	if err != nil {
		//		log.Println("unable to get the school name")
		//	}
		//	student.School.Name = name
		//}
		//student.CreatedAt = values[6].(time.Time)
		//student.UpdatedAt = values[7].(time.Time)
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

func (r *studentRepository) UpdateCurrentClass(ctx context.Context, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updateCurrentClasses, st.CurrentClasses, time.Now(), st.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *studentRepository) UpdateClassesTaken(ctx context.Context, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updateClassesTaken, st.ID, time.Now(), st.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *studentRepository) CompleteClass(ctx context.Context, st *domain.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updateClassesTaken, st.ID, time.Now(), st.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	_, err = tx.Exec(ctx, updateCurrentClasses, st.ID, time.Now(), st.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}