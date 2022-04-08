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

const filePath string = "whitelist.data"

//home endpoint
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
//add new device to whitelist
func credentials(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")

		var allow string = "no"

		//read file which contains all allowed devices ids
		file, err := os.Open(filePath)
		if err != nil {
			log.Panic(err)
		}

		defer file.Close()
		//search for matching id in file with given from request
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			if text == id {
				allow = "yes"
				break
			}
		}

		if allow == "yes" {
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
		w.Write([]byte(allow))
	}

	if r.Method == "PUT" {
		w.Header().Set("Allow", "PUT")
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
	}

	if r.Method == "POST" {
		//MAPING DATA FROM REQ.BODY TO STRUCT
		var credential Credential
		bodyBytes, err := io.ReadAll(r.Body)
		json.Unmarshal([]byte(bodyBytes), &credential)
		fmt.Printf("%s", &credential)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			panic(err)

		}
		//SAVE TO DATA FILE
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			panic(err)
		}
		//CONCAT 2 STRINGS
		str := fmt.Sprint(credential.Code, "\n")
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

	//API ENDPOINTS:
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/whitelist", credentials)

	//RUNMING SERVER
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
