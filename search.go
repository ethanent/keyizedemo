package main

import (
	"encoding/json"
	"fmt"
	"github.com/KeyizeBiometry/keyize"
	"net/http"
)

type CompareResult struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	PropMatch  float64 `json:"propMatch"`
	Confidence float64 `json:"confidence"`
	Dist       float64 `json:"dist"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	d := &struct {
		V1Recording string `json:"rec"`
	}{}

	if err := unmarshalJSON(r, d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse rec

	rec, err := keyize.ImportKeyizeV1(d.V1Recording)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dyn := rec.Dynamics()

	// Get network

	nd, ok := data[r.RemoteAddr]

	if ok == false {
		http.Error(w, "No recordings have been uploaded yet.", http.StatusNotFound)
		return
	}

	// Compare rec to network user dynamics

	res := &struct {
		Results []*CompareResult `json:"results"`
	}{
		Results: []*CompareResult{},
	}

	for id, user := range nd.users {
		confidence := dyn.ProportionSharedProperties(user.avgDyn, keyize.Both)

		propMatch := dyn.ProportionMatch(user.avgDyn)

		dist := dyn.AvgScaledPropDiff(user.avgDyn, nil)

		res.Results = append(res.Results, &CompareResult{
			ID:         id,
			Name:       user.name,
			PropMatch:  propMatch,
			Confidence: confidence,
			Dist:       dist,
		})
	}

	// Write resp

	resd, err := json.Marshal(res)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resd)
}
