package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/errors"
)

// Update state of EDV
// TODO: may need to add global locking around this function to
// avoid inconsistent state from concurrent client updates
func UpdateEdvState(edvId string, docId string, operation string) {
	// Retrieve and parse config
	var edvConfig DataVaultConfiguration
	configFileName := fmt.Sprintf("./edvs/%s/config.json", edvId)
	configFileBytes, _ := os.ReadFile(configFileName)
	json.Unmarshal(configFileBytes, &edvConfig)

	// Retrieve and parse history
	var historyEntries []EdvHistoryLogEntry
	historyFileName := fmt.Sprintf("./edvs/%s/history.json", edvId)
	historyFileBytes, _ := os.ReadFile(historyFileName)
	json.Unmarshal(historyFileBytes, &historyEntries)

	// Update parsed config and history
	edvConfig.Sequence++
	historyEntry := EdvHistoryLogEntry{docId, edvConfig.Sequence, operation}
	historyEntries = append(historyEntries, historyEntry)
	historyFileBytes, _ = json.MarshalIndent(historyEntries, "", "  ")
	configFileBytes, _ = json.MarshalIndent(edvConfig, "", "  ")

	// Persist updated config and history
	os.WriteFile(configFileName, configFileBytes, os.ModePerm)
	os.WriteFile(historyFileName, historyFileBytes, os.ModePerm)
}

// Create EDV
func CreateEdv(res http.ResponseWriter, req *http.Request) {
	var edvConfig DataVaultConfiguration
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	bodyUnmarshalErr := json.Unmarshal(body, &edvConfig)

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
	if edvConfig.Id != "" {
		edvId = edvConfig.Id
	} else {
		edvId = uuid.NewString()
	}
	edvConfig.Id = edvId

	edvDirName := filepath.Join(".", "edvs", edvId)
	docDirName := filepath.Join(edvDirName, "docs")
	os.MkdirAll(docDirName, os.ModePerm)
	configFileName := fmt.Sprintf("./edvs/%s/config.json", edvId)
	configFile, _ := os.Create(configFileName)
	configFileBytes, _ := json.MarshalIndent(edvConfig, "", "  ")
	configFile.Write(configFileBytes)
	historyFileName := fmt.Sprintf("./edvs/%s/history.json", edvId)
	historyFile, _ := os.Create(historyFileName)
	historyFileString := "[]"
	historyFile.WriteString(historyFileString)
	edvLocation := fmt.Sprintf("%s/edvs/%s", req.Host, edvId)
	res.Header().Add("Location", edvLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all EDVs
func GetEdvs(res http.ResponseWriter, req *http.Request) {}

// Get EDV
func GetEdv(res http.ResponseWriter, req *http.Request) {}

// Get history of EDV
func GetEdvHistory(res http.ResponseWriter, req *http.Request) {
	var historyEntries []EdvHistoryLogEntry
	var historyEntriesFiltered []EdvHistoryLogEntry

	afterSequenceString := req.URL.Query().Get("afterSequence")
	beforeSequenceString := req.URL.Query().Get("beforeSequence")
	var afterSequence uint64
	var beforeSequence uint64

	edvId := mux.Vars(req)["edvId"]
	historyFileName := fmt.Sprintf("./edvs/%s/history.json", edvId)
	historyFileBytes, _ := os.ReadFile(historyFileName)
	json.Unmarshal(historyFileBytes, &historyEntries)

	if afterSequenceString != "" {
		afterSequence, _ = strconv.ParseUint(afterSequenceString, 10, 64)
	} else {
		afterSequence = 0
	}
	if beforeSequenceString != "" {
		beforeSequence, _ = strconv.ParseUint(beforeSequenceString, 10, 64)
	} else {
		beforeSequence = math.MaxUint64
	}

	for _, entry := range historyEntries {
		if entry.Sequence > afterSequence && entry.Sequence < beforeSequence {
			historyEntriesFiltered = append(historyEntriesFiltered, entry)
		}
	}
	historyFileBytesFiltered, _ := json.MarshalIndent(historyEntriesFiltered, "", "  ")
	res.Write(historyFileBytesFiltered)
}

// Search EDV
func SearchEdv(res http.ResponseWriter, req *http.Request) {}
