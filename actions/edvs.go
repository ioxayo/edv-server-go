package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ioxayo/edv-server-go/errors"
)

// Create EDV
func CreateEdv(res http.ResponseWriter, req *http.Request) {
	var ceReq CreateEdvRequest
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	bodyUnmarshalErr := json.Unmarshal(body, &ceReq)

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

	var edvId string
	if ceReq.Id != "" {
		edvId = ceReq.Id
	} else {
		edvId = uuid.NewString()
	}
	ceReq.Id = edvId

	var configFileName string
	var configFile *os.File
	edvDir := filepath.Join(".", "edvs", edvId)
	os.MkdirAll(edvDir, os.ModePerm)
	configFileName = fmt.Sprintf("./edvs/%s/config.json", edvId)
	configFile, _ = os.Create(configFileName)
	configFileBytes, _ := json.MarshalIndent(ceReq, "", "  ")
	configFile.Write(configFileBytes)
	edvLocation := fmt.Sprintf("%s/edvs/%s", req.Host, edvId)
	res.Header().Add("Location", edvLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all EDVs
func GetEdvs(res http.ResponseWriter, req *http.Request) {}

// Get EDV
func GetEdv(res http.ResponseWriter, req *http.Request) {}

// Search EDV
func SearchEdv(res http.ResponseWriter, req *http.Request) {}
