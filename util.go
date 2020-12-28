package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

func genID() (string, error) {
	d := make([]byte, 16)

	if _, err := rand.Read(d); err != nil {
		return "", err
	}

	return hex.EncodeToString(d), nil
}

func getIP(addr string) string {
	ip := strings.Split(addr, ":")[0]
	return ip
}
