package impl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-oci8"
	"github.com/xtracdev/oraeventstore"
	"github.com/xtraclabs/appreg/domain"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	eventStore *oraeventstore.OraEventStore
	db         *sql.DB
)

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
		log.Fatal(strings.Join(configErrors, "\n"))
	}

	var err error
	eventStore, err = oraeventstore.NewOraEventStore(user, password, dbSvc, dbhost, dbPort)
	if err != nil {
		log.Fatalf("Error connecting to oracle: %s", err.Error())
	}

	connectStr := fmt.Sprintf("%s/%s@//%s:%s/%s",
		user, password, dbhost, dbPort, dbSvc)

	db, err = sql.Open("oci8", connectStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	//Are we really in an ok state for starters?
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to oracle: %s", err.Error())
	}

}

func buildGetByIdResponse(appReg *domain.ApplicationReg) ([]byte, error) {
	response := make(map[string]interface{})
	response["appVersion"] = "1.0"

	data := make(map[string]interface{})
	data["name"] = appReg.Name
	data["description"] = appReg.Description
	data["client_id"] = appReg.ID

	created := time.Unix(0, appReg.Created).Format(time.RFC3339Nano)
	data["created"] = created

	response["data"] = data

	return json.Marshal(response)
}

func ApplicationsClientIdGet(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["client_id"]
	log.Printf("Read client id '%s'\n", clientID)

	events, err := eventStore.RetrieveEvents(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(events) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	appReg := domain.NewApplicationRegFromHistory(events)

	log.Println("Read app from event store...", appReg)

	out, err := buildGetByIdResponse(appReg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
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

	response := make(map[string]interface{})
	response["apiVersion"] = "1.0"
	var data []interface{}

	rows, err := db.Query(`select client_id, name, created from app_summary order by name`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var clientId, name string
	var create time.Time

	for rows.Next() {
		rows.Scan(&clientId, &name, &create)

		row := make(map[string]interface{})
		row["name"] = name
		row["clientId"] = clientId
		row["created"] = create.Format(time.RFC3339Nano)

		data = append(data, row)
	}

	response["data"] = data

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
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

	log.Println("create app with name and desc", name, desc)
	appReg, _ := domain.NewApplicationReg(name, desc)
	log.Println("storing...", appReg)

	err = appReg.Store(eventStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["apiVersion"] = "1.0"
	responseData := make(map[string]interface{})
	responseData["client_id"] = appReg.ID
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
