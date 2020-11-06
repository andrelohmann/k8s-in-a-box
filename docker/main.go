package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
  "os"
)

type request struct {
	URL        string      `json:"url"`
	Method     string      `json:"method"`
	Hostname   string      `json:"hostname"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
}

func handle(rw http.ResponseWriter, r *http.Request) {
  h_name, h_err := os.Hostname()
  if h_err != nil {
		panic(h_err)
	}
	var err error
	rr := &request{}
	rr.Method = r.Method
	rr.Headers = r.Header
  rr.Hostname = h_name
	rr.URL = r.URL.String()
	rr.Body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rrb, err := json.Marshal(rr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(rrb)
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8000", nil)
}
