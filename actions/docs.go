package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/errors"
	"github.com/ioxayo/edv-server-go/storage"
)

// Get document by ID
func GetDocumentById(edvId string, docId string) common.EncryptedDocument {
	var doc common.EncryptedDocument
	provider, _ := storage.GetStorageProvider(edvId)
	docFileBytes, _ := provider.ReadDocClient(edvId, docId)
	json.Unmarshal(docFileBytes, &doc)
	return doc
}

// Get documents by ID
func GetDocumentsById(edvId string, docIds []string) []common.EncryptedDocument {
	docs := make([]common.EncryptedDocument, 0)
	for _, docId := range docIds {
		doc := GetDocumentById(edvId, docId)
		docs = append(docs, doc)
	}
	return docs
}

// Create document
func CreateDocument(res http.ResponseWriter, req *http.Request) {
	var doc common.EncryptedDocument
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &doc)

	if bodyReadErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyReadErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	edvId := mux.Vars(req)["edvId"]
	docId := doc.Id
	provider, _ := storage.GetStorageProvider(edvId)

	docLocation, createDocErr := provider.CreateDocClient(edvId, docId, body)
	if createDocErr.IsError() {
		message := createDocErr.Message
		status := createDocErr.Status
		errors.HandleError(res, req, message, status)
		return
	}

	UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Create)
	UpdateEdvIndexCreate(edvId, doc)

	res.Header().Add("Location", docLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all documents
func GetDocuments(res http.ResponseWriter, req *http.Request) {}

// Get document
func GetDocument(res http.ResponseWriter, req *http.Request) {
	edvId := mux.Vars(req)["edvId"]
	docId := mux.Vars(req)["docId"]
	provider, _ := storage.GetStorageProvider(edvId)
	docFileBytes, getDocErr := provider.ReadDocClient(edvId, docId)
	if getDocErr.IsError() {
		message := getDocErr.Message
		status := getDocErr.Status
		errors.HandleError(res, req, message, status)
		return
	}

	res.Write(docFileBytes)
}

// Update document
func UpdateDocument(res http.ResponseWriter, req *http.Request) {
	var doc common.EncryptedDocument
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &doc)

	if bodyReadErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyReadErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	edvId := mux.Vars(req)["edvId"]
	docId := mux.Vars(req)["docId"]
	provider, _ := storage.GetStorageProvider(edvId)

	updateDocErr := provider.UpdateDocClient(edvId, docId, body)
	if updateDocErr.IsError() {
		message := updateDocErr.Message
		status := updateDocErr.Status
		errors.HandleError(res, req, message, status)
		return
	}

	UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Update)
	UpdateEdvIndexUpdate(edvId, doc)
}

// Delete document
func DeleteDocument(res http.ResponseWriter, req *http.Request) {
	edvId := mux.Vars(req)["edvId"]
	docId := mux.Vars(req)["docId"]
	provider, _ := storage.GetStorageProvider(edvId)

	deleteDocErr := provider.DeleteDocClient(edvId, docId)
	if deleteDocErr.IsError() {
		message := deleteDocErr.Message
		status := deleteDocErr.Status
		errors.HandleError(res, req, message, status)
		return
	}

	UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Delete)
	UpdateEdvIndexDelete(edvId, docId)
}
