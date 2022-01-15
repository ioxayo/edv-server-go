package main

import (
	"net/http"

	"github.com/ioxayo/edv-server-go/actions"
)

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Create router
var routes Routes = Routes{
	Route{
		"CreateEdv",
		"POST",
		"/",
		actions.CreateEdv,
	},
	Route{
		"GetEdvs",
		"GET",
		"/",
		actions.GetEdvs,
	},
	Route{
		"GetEdv",
		"GET",
		"/{edvId}",
		actions.GetEdv,
	},
	Route{
		"SearchEdv",
		"GET",
		"/{edvId}/search",
		actions.SearchEdv,
	},
	Route{
		"CreateDocument",
		"POST",
		"/{edvId}/documents",
		actions.CreateDocument,
	},
	Route{
		"GetDocuments",
		"GET",
		"/{edvId}/documents",
		actions.GetDocuments,
	},
	Route{
		"GetDocument",
		"GET",
		"/{edvId}/document/{docId}",
		actions.GetDocument,
	},
	Route{
		"UpdateDocument",
		"POST",
		"/{edvId}/document/{docId}",
		actions.UpdateDocument,
	},
}
