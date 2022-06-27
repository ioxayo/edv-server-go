package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/storage"
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
	// Access $EDV_PORT env var
	port := os.Getenv("EDV_PORT")
	if port == "" {
		log.Fatal("$EDV_PORT must be set")
	}

	// Setup storage provider
	// TODO: accept host and root from stdin
	// TODO: setup switch statement for different storage provider types
	edvHost := "http://localhost:5000"
	currentDir, _ := os.Getwd()
	// storage.InitLocalStorageProvider(storage.Provider, edvHost, currentDir)
	storage.Provider = storage.InitLocalStorageProvider(edvHost, currentDir)

	// Setup router
	router := NewRouter(routes)

	// Restrict client interactions
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	// Launch server
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
}
