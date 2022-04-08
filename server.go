package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Credential struct {
	Code string `json:"Code"`
}

var credentialsArr [3]string

const filePath string = "whitelist.data"

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello there"))

	if r.Method == "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
		return
	}
}

//send device id to get information if it is allow to send data
func credentials(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")

		var allow string = "no"

		file, err := os.Open(filePath)
		if err != nil {
			log.Panic(err)
		}

		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			if text == id {
				allow = "yes"
				break
			}
		}
		w.Write([]byte(allow))
	}

	if r.Method == "PUT" {
		w.Header().Set("Allow", "PUT")
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
	}

	if r.Method == "POST" {
		var credential Credential
		bodyBytes, err := io.ReadAll(r.Body)
		json.Unmarshal([]byte(bodyBytes), &credential)
		fmt.Printf("%s", &credential)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			panic(err)

		}
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			panic(err)
		}

		str := fmt.Sprint(credential.Code, "\n")
		fmt.Printf(str)
		fmt.Printf(str)
		if _, err := file.WriteString(str); err != nil {
			w.WriteHeader(http.StatusConflict)
			panic(err)
		}
		defer file.Close()
		w.WriteHeader(http.StatusCreated)
	}
	return
}

func main() {

	credentialsArr = [...]string{"0001", "0002", "0003"}
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/whitelist", credentials)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
