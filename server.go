package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/common"
)

// NewRouter function configures a new router to the API
func NewRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		log.Println(route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(route.HandlerFunc)
	}
	return router
}

func main() {
	// Access $PORT env var
	portVarName := common.EnvVars.Port
	port := os.Getenv(portVarName)
	if port == "" {
		errorMessage := fmt.Sprintf("$%s must be set", portVarName)
		log.Fatal(errorMessage)
	}

	// Setup router
	router := NewRouter(routes)

	// Restrict client interactions
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	// Launch server
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
}
