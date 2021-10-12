package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/jackc/pgx/v4/pgxpool"
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
