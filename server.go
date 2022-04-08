package main

import (
	"log"
	"net/http"
)

var credentialsArr [3]string

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
func isAllowTo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		present := "no"
		for _, x := range credentialsArr {
			if x == id {
				present = "yes"
				break
			}
		}

		w.Write([]byte(present))
		return
	}
	if r.Method == "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
	}
}

func main() {

	credentialsArr = [...]string{"0001", "0002", "0003"}
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/isallowto", isAllowTo)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
