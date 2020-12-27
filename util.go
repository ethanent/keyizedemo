package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func unmarshalJSON(r *http.Request, d interface{}) error {
	lr := &io.LimitedReader{
		R: r.Body,
		N: 20 * 1000,
	}

	all, err := ioutil.ReadAll(lr)

	if err != nil {
		return err
	}

	return json.Unmarshal(all, d)
}
