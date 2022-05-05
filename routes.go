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
		"/edvs",
		actions.CreateEdv,
	},
	Route{
		"GetEdvs",
		"GET",
		"/edvs",
		actions.GetEdvs,
	},
	Route{
		"GetEdv",
		"GET",
		"/edvs/{edvId}",
		actions.GetEdv,
	},
	Route{
		"SearchEdv",
		"POST",
		"/edvs/{edvId}/query",
		actions.SearchEdv,
	},
	Route{
		"CreateDocument",
		"POST",
		"/edvs/{edvId}/docs",
		actions.CreateDocument,
	},
	Route{
		"GetDocuments",
		"GET",
		"/edvs/{edvId}/docs",
		actions.GetDocuments,
	},
	Route{
		"GetDocument",
		"GET",
		"/edvs/{edvId}/docs/{docId}",
		actions.GetDocument,
	},
	Route{
		"UpdateDocument",
		"POST",
		"/edvs/{edvId}/docs/{docId}",
		actions.UpdateDocument,
	},
	Route{
		"DeleteDocument",
		"DELETE",
		"/edvs/{edvId}/docs/{docId}",
		actions.DeleteDocument,
	},
	Route{
		"GetEdvHistory",
		"GET",
		"/edvs/{edvId}/history",
		actions.GetEdvHistory,
	},
}
