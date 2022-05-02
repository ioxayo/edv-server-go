package actions

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/errors"
)

// Get document by ID
func GetDocumentById(edvId string, docId string) EncryptedDocument {
	var doc EncryptedDocument
	docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	docFileBytes, _ := os.ReadFile(docFileName)
	json.Unmarshal(docFileBytes, &doc)
	return doc
}

// Get documents by ID
func GetDocumentsById(edvId string, docIds []string) []EncryptedDocument {
	docs := make([]EncryptedDocument, 0)
	for _, docId := range docIds {
		doc := GetDocumentById(edvId, docId)
		docs = append(docs, doc)
	}
	return docs
}

// Create document
func CreateDocument(res http.ResponseWriter, req *http.Request) {
	var doc EncryptedDocument
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	bodyUnmarshalErr := json.Unmarshal(body, &doc)

	if bodyReadErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyReadErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	if bodyUnmarshalErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyUnmarshalErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	edvId := mux.Vars(req)["edvId"]
	edvDirName := fmt.Sprintf("./edvs/%s", edvId)
	if _, err := os.Stat(edvDirName); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find EDV with ID '%s'", edvId)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}
	docId := doc.Id
	docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(doc, "", "  ")
	docFile.Write(docFileBytes)
	UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Create)
	UpdateEdvIndexCreate(edvId, doc)

	docLocation := fmt.Sprintf("%s/edvs/%s/docs/%s", req.Host, edvId, docId)
	res.Header().Add("Location", docLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all documents
func GetDocuments(res http.ResponseWriter, req *http.Request) {}

// Get document
func GetDocument(res http.ResponseWriter, req *http.Request) {
	edvId := mux.Vars(req)["edvId"]
	docId := mux.Vars(req)["docId"]
	docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	if _, err := os.Stat(docFileName); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find document with ID '%s' in EDV with ID '%s'", docId, edvId)
		status := http.StatusNotFound
		errors.HandleError(res, req, message, status)
		return
	}
	docFileBytes, _ := os.ReadFile(docFileName)
	res.Write(docFileBytes)
}

// Update document
func UpdateDocument(res http.ResponseWriter, req *http.Request) {
	var doc EncryptedDocument
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	bodyUnmarshalErr := json.Unmarshal(body, &doc)

	if bodyReadErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyReadErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	if bodyUnmarshalErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyUnmarshalErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	edvId := mux.Vars(req)["edvId"]
	docId := mux.Vars(req)["docId"]
	docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	if _, err := os.Stat(docFileName); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find document with ID '%s' in EDV with ID '%s'", docId, edvId)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(doc, "", "  ")
	docFile.Write(docFileBytes)
	UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Update)
}

// Delete document
func DeleteDocument(res http.ResponseWriter, req *http.Request) {
	edvId := mux.Vars(req)["edvId"]
	docId := mux.Vars(req)["docId"]
	docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	if _, err := os.Stat(docFileName); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find document with ID '%s' in EDV with ID '%s'", docId, edvId)
		status := http.StatusNotFound
		errors.HandleError(res, req, message, status)
		return
	}
	os.Remove(docFileName)
	UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Delete)
}
