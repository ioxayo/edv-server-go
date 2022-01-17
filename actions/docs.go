package actions

import (
	"encoding/json"
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
	docId := cdReq.Id
	docFileName := fmt.Sprintf("./edvs/%s/%s.json", edvId, docId)
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(cdReq, "", "  ")
	docFile.Write(docFileBytes)
	docLocation := fmt.Sprintf("%s/edvs/%s/%s", req.Host, edvId, docId)
	res.Header().Add("Location", docLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all documents
func GetDocuments(res http.ResponseWriter, req *http.Request) {}

// Get document
func GetDocument(res http.ResponseWriter, req *http.Request) {}

// Update document
func UpdateDocument(res http.ResponseWriter, req *http.Request) {}
