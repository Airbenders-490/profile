// Package utils adds all the schools to the db. Import this package to add all the schools
// WARNING: DO THIS ONLY ONCE EVER FOR AN APPLICATION
// OTHERWISE ALL THE IDs OF SCHOOLS WILL BE RANDOM
package utils

import (
	"context"
	"encoding/json"
	"github.com/airbenders/profile/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		addSchoolsToDB()
	}
}

func addSchoolsToDB() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	insert := `INSERT INTO school (id, name, country, domains) VALUES ($1, $2, $3, $4)`
	tx, err := pool.Begin(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback(context.Background())

	var schools []domain.School
	resp, _ := http.Get("https://raw.githubusercontent.com/Hipo/university-domains-list/master/world_universities_and_domains.json")
	schoolJson, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(schoolJson, &schools)

	for _, school := range schools {
		school.ID = uuid.New().String()
		_, err = tx.Exec(context.Background(), insert, school.ID, school.Name, school.Country, school.Domains)
		if err != nil {
			log.Fatalln(err)
		}

		tx.Commit(context.Background())
	}
}