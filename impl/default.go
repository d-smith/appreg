package impl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Default struct {
}

func ApplicationsClientIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ApplicationsClientIdPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ApplicationsClientIdSecretPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ApplicationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ApplicationsPost(w http.ResponseWriter, r *http.Request) {
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	log.Println(string(rb))

	var parsed map[string]interface{}

	err = json.Unmarshal(rb, &parsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("name", parsed["name"])
	log.Println("description", parsed["description"])

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Yeah allright"))
}
