package utils// Package utils adds all the schools to the db. Import this package to add all the schools
import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// WARNING: DO THIS ONLY ONCE EVER FOR AN APPLICATION
// OTHERWISE ALL THE IDs OF SCHOOLS WILL BE RANDOM


func createJSONFile() {
	resp, _ := http.Get("https://raw.githubusercontent.com/Hipo/university-domains-list/master/world_universities_and_domains.json%22")
	//schoolJSON, _ := ioutil.ReadAll(resp.Body)
	//// The os.Create creates or truncates the named file. If the file already exists, it is truncated.
	//jsonFile,_ := os.Create("school.json")
	//data := []byte(schoolJSON)
	//_, err2 := jsonFile.Write(data)
	//if err2 != nil {
	//	log.Fatal(err2)
	//}

}
