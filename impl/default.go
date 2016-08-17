package impl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/xtraclabs/appreg/domain"
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

	name, ok := parsed["name"].(string)
	if !ok {
		http.Error(w, "Unable to extract name parameter from request payload", http.StatusInternalServerError)
	}

	desc, ok := parsed["description"].(string)
	if !ok {
		http.Error(w, "Unable to extract description parameter from request payload", http.StatusInternalServerError)
	}


	appReg,_ := domain.NewApplicationReg(name, desc)
	log.Println(appReg)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Yeah allright"))
}
