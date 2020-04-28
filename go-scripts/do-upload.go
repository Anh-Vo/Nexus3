package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type payload struct {
	url      string
	username string
	password string
	filename string
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	host := flag.String("h", "http://localhost:8080", "The host of interest")
	path := flag.String("endpoint", "/service/rest/v1/script", "The endpoint of interest")
	username := flag.String("u", "admin", "username")
	password := flag.String("p", "password", "password")
	filename := flag.String("f", "dummy.txt", "The file to upload")
	flag.Parse()
	URL := *host + *path

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	payload := &payload{
		url:      URL,
		username: *username,
		password: *password,
		filename: *filename,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	payload.doPost()
}

func (p *payload) doPost() {
	p.infoLog.Printf("POST to '%s' as user '%s' with file '%s'\n", p.url, p.username, p.filename)
	f, err := os.Open(p.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	req, err := http.NewRequest("POST", p.url, f)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(p.username, p.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		p.errorLog.Printf("Something went wrong, the server returned a: '%s'", resp.Status)
		return
	}
	defer resp.Body.Close()
	fmt.Println("*** Success ***")
	fmt.Printf("The server returned a: '%s'\n", resp.Status)
	fmt.Printf("{\n %s \n}\n", resp.Body)
}
