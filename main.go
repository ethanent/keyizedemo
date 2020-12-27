package main

import (
	"fmt"
	"github.com/KeyizeBiometry/keyize"
	"net/http"
)

type UserData struct {
	name   string
	avgDyn *keyize.Dynamics
}

type NetworkData struct {
	users map[string]*UserData
}

var data = map[string]*NetworkData{}

func main() {
	s := http.NewServeMux()

	s.Handle("/", http.FileServer(http.Dir("./static")))

	s.HandleFunc("/upload", uploadHandler)
	s.HandleFunc("/search", searchHandler)

	s.HandleFunc("/listUsers", listUsersHandler)
	s.HandleFunc("/removeUsers", removeUsersHandler)

	fmt.Println("Listening...")

	if err := http.ListenAndServe(":80", s); err != nil {
		panic(err)
	}
}
