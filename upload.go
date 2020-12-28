package main

import (
	"fmt"
	"github.com/KeyizeBiometry/keyize"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r.RemoteAddr)

	d := &struct {
		Name         string   `json:"name"`
		V1Recordings []string `json:"recs"`
	}{}

	if err := unmarshalJSON(r, d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure at least 3 recordings provided & name len > 0

	if len(d.V1Recordings) < 3 || len(d.Name) == 0 {
		http.Error(w, "Invalid request. Should have len(recs) >= 3, len(name) > 0", http.StatusBadRequest)
		return
	}

	// Average dynamics

	dynamics := []*keyize.Dynamics{}

	for _, curV1Rec := range d.V1Recordings {
		// Parse rec
		rec, err := keyize.ImportKeyizeV1(curV1Rec)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Append dynamics to dynamics

		dynamics = append(dynamics, rec.Dynamics())
	}

	// Average dynamics from recordings

	avgDynamics := keyize.AvgDynamics(dynamics)

	// Insert data

	if _, ok := data[ip]; ok == false {
		data[ip] = &NetworkData{
			users: map[string]*UserData{},
		}
	}

	id, err := genID()

	if err != nil {
		fmt.Println(err)
		return
	}

	nd := data[ip]

	nd.users[id] = &UserData{
		name:   d.Name,
		avgDyn: avgDynamics,
	}

	fmt.Println("Added user for", ip, ":", d.Name)

	// Reply

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
