package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	d, ok := data[r.RemoteAddr]

	if !ok {
		http.Error(w, "No recordings have been uploaded yet.", http.StatusNotFound)
		return
	}

	res := []*UserInfo{}

	for id, u := range d.users {
		res = append(res, &UserInfo{
			ID:   id,
			Name: u.name,
		})
	}

	rd, err := json.Marshal(res)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rd)
}

func removeUsersHandler(w http.ResponseWriter, r *http.Request) {
	d, ok := data[r.RemoteAddr]

	if !ok {
		http.Error(w, "No recordings have been uploaded yet.", http.StatusNotFound)
		return
	}

	userIDs := []string{}

	if err := unmarshalJSON(r, &userIDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, uid := range userIDs {
		delete(d.users, uid)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
