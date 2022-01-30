package utils// Package utils adds all the schools to the db. Import this package to add all the schools
import (
	"encoding/json"
	"github.com/airbenders/profile/domain"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// WARNING: DO THIS ONLY ONCE EVER FOR AN APPLICATION
// OTHERWISE ALL THE IDs OF SCHOOLS WILL BE RANDOM


func createJSONFile() {
	var schools []domain.School
	resp, _ := http.Get("https://raw.githubusercontent.com/Hipo/university-domains-list/master/world_universities_and_domains.json")
	schoolJSON, _ := ioutil.ReadAll(resp.Body)
	// The os.Create creates or truncates the named file. If the file already exists, it is truncated.
	jsonFile,_ := os.Create("school.json")
	data := []byte(schoolJSON)
	_, err2 := jsonFile.Write(data)
	if err2 != nil {
		log.Fatal(err2)
	}
	jsonFile, _ = os.Open("school.json")
	schoolJSON, _ = ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(schoolJSON, &schools)

}
