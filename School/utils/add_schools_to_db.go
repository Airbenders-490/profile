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
	"path"
	"strings"
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		addSchoolsToDB()
	}
}

const (
	insert = `INSERT INTO school (id, name, country, domains) VALUES ($1, $2, $3, $4)`
	get = `SELECT * FROM school LIMIT 1`
)

func addSchoolsToDB() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := pool.Query(context.Background(), get)
	if err != nil {
		log.Fatalln("can't get the rows in school db")
	}

	// if a row already exist means db is already populated or being populated. So return
	for rows.Next() {
		return
	}
	// if reached here, then we can create schools
	fmt.Println("creating schools")

	tx, err := pool.Begin(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback(context.Background())

	var schools []domain.School

	pathToFile := getPathToSchoolsFile(err)

	schoolJSON, err := ioutil.ReadFile(pathToFile)
	fatalOnError(err)

	_ = json.Unmarshal(schoolJSON, &schools)
	for _, school := range schools {
		_, err = tx.Exec(context.Background(), insert, school.ID, school.Name, school.Country, school.Domains)
		if err != nil {
			log.Fatalln(err)
		}

		tx.Commit(context.Background())
	}
}

func getPathToSchoolsFile(err error) string {
	cwd, err := os.Getwd()
	fatalOnError(err)
	pathToFile := path.Join(cwd, "School", "utils", "schools.json")
	log.Println("the path to file is", pathToFile)
	return pathToFile
}

func fatalOnError(err error) {
	if err != nil {
		log.Fatalln("can't find the file!")
	}
}
