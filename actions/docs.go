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

// Create document
func CreateDocument(res http.ResponseWriter, req *http.Request) {
	var cdReq CreateDocumentRequest
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	bodyUnmarshalErr := json.Unmarshal(body, &cdReq)

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
	docId := cdReq.Id
	docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(cdReq, "", "  ")
	docFile.Write(docFileBytes)
	docLocation := fmt.Sprintf("%s/edvs/%s/docs/%s", req.Host, edvId, docId)
	res.Header().Add("Location", docLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all documents
func GetDocuments(res http.ResponseWriter, req *http.Request) {}

// Get single document
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
	docFile, _ := os.ReadFile(docFileName)
	res.WriteHeader(http.StatusOK)
	res.Write(docFile)
}

// Update document
func UpdateDocument(res http.ResponseWriter, req *http.Request) {}
