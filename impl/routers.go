package impl

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/xtrac/devcenter/",
		Index,
	},

	Route{
		"ApplicationsClientIdGet",
		"GET",
		"/xtrac/devcenter/applications/{client_id}",
		ApplicationsClientIdGet,
	},

	Route{
		"ApplicationsClientIdPut",
		"PUT",
		"/xtrac/devcenter/applications/{client_id}",
		ApplicationsClientIdPut,
	},

	Route{
		"ApplicationsClientIdSecretPost",
		"POST",
		"/xtrac/devcenter/applications/{client_id}/secret",
		ApplicationsClientIdSecretPost,
	},

	Route{
		"ApplicationsGet",
		"GET",
		"/xtrac/devcenter/applications",
		ApplicationsGet,
	},

	Route{
		"ApplicationsPost",
		"POST",
		"/xtrac/devcenter/applications",
		ApplicationsPost,
	},

}