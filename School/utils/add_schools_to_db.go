// Package utils adds all the schools to the db. Import this package to add all the schools
// WARNING: DO THIS ONLY ONCE EVER FOR AN APPLICATION
// OTHERWISE ALL THE IDs OF SCHOOLS WILL BE RANDOM
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		addSchoolsToDB()
	}
}

func addSchoolsToDB() {
	fmt.Println("creating schools")
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
	createJSONFile()
	jsonFile, _ := os.Open("school.json")
	schoolJSON, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(schoolJSON, &schools)

	for _, school := range schools {
		_, err = tx.Exec(context.Background(), insert, school.ID, school.Name, school.Country, school.Domains)
		if err != nil {
			log.Fatalln(err)
		}

		tx.Commit(context.Background())
	}
}
