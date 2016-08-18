package impl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/xtraclabs/appreg/domain"
	"github.com/xtraclabs/oraeventstore"
	"os"
	"strings"
)

var eventStore *oraeventstore.OraEventStore

type Default struct {
}

func init() {
	var configErrors []string

	user := os.Getenv("FEED_DB_USER")
	if user == "" {
		configErrors = append(configErrors, "Configuration missing FEED_DB_USER env variable")
	}

	password := os.Getenv("FEED_DB_PASSWORD")
	if password == "" {
		configErrors = append(configErrors, "Configuration missing FEED_DB_PASSWORD env variable")
	}

	dbhost := os.Getenv("FEED_DB_HOST")
	if dbhost == "" {
		configErrors = append(configErrors, "Configuration missing FEED_DB_HOST env variable")
	}

	dbPort := os.Getenv("FEED_DB_PORT")
	if dbPort == "" {
		configErrors = append(configErrors, "Configuration missing FEED_DB_PORT env variable")
	}

	dbSvc := os.Getenv("FEED_DB_SVC")
	if dbSvc == "" {
		configErrors = append(configErrors, "Configuration missing FEED_DB_SVC env variable")
	}

	if len(configErrors) != 0 {
		log.Fatal(strings.Join(configErrors,"\n"))
	}

	var err error
	eventStore, err = oraeventstore.NewOraEventStore("esusr", "password", "xe.oracle.docker", "localhost", "1521")
	if err != nil {
		log.Fatalf("Error connecting to oracle: %s", err.Error())
	}

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

	err = appReg.Store(eventStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	response := make(map[string]interface{})
	response["apiVersion"] = 1.0
	responseData := make(map[string]interface{})
	responseData["client_id"] = appReg.ClientID
	response["data"] = responseData

	outbytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(outbytes)
}
