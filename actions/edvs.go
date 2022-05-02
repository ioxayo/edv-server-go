package actions

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/errors"
)

// Update state of EDV
// TODO: may need to add global locking around this function to
// avoid inconsistent state from concurrent client updates
func UpdateEdvState(edvId string, docId string, operation string) {
	if !common.IsValidEnumMember(EncryptedDocumentOperations, operation) {
		return
	}
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

	// TODO: Retrieve and parse index

	// Update parsed config and history
	edvConfig.Sequence++
	historyEntry := EdvHistoryLogEntry{docId, edvConfig.Sequence, operation}
	historyEntries = append(historyEntries, historyEntry)
	historyFileBytes, _ = json.MarshalIndent(historyEntries, "", "  ")
	configFileBytes, _ = json.MarshalIndent(edvConfig, "", "  ")
	// TODO: Update parsed index

	// Persist updated config and history
	os.WriteFile(configFileName, configFileBytes, os.ModePerm)
	os.WriteFile(historyFileName, historyFileBytes, os.ModePerm)
}

// Retrieve document IDs associated with an index ID
func IndexToDocuments(edvId string, indexId string) []string {
	var indexEntries map[string][]string
	indexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
	indexFileBytes, _ := os.ReadFile(indexFileName)
	json.Unmarshal(indexFileBytes, &indexEntries)
	return indexEntries[indexId]
}

// Returns all document IDs for which condition is met for all key-value pairs of subfilter of given query operator
func FetchMatchesAll(edvId string, subfilter map[string]string, operator string) []string {
	indexId := subfilter["index"]
	docIds := IndexToDocuments(edvId, indexId)
	docMatches := make([]string, 0)
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return docMatches
	}
	for _, docId := range docIds {
		docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
		if _, err := os.Stat(docFileName); goerrors.Is(err, os.ErrNotExist) {
			continue
		}
		docFileBytes, _ := os.ReadFile(docFileName)
		var encDoc EncryptedDocument
		if err := json.Unmarshal(docFileBytes, &encDoc); err != nil {
			continue
		}
		filterMatches := make(map[string]bool)
		switch operator {
		case EdvSearchOperators.Equals:
			indexes := encDoc.Indexed
			for _, index := range indexes {
				if index.Hmac.Id == indexId {
					attributes := index.Attributes
					for _, attribute := range attributes {
						attributeName := attribute.Name
						attributeValue := attribute.Value
						if subfilterValue, subfilterExists := subfilter[attributeName]; subfilterValue == attributeValue && subfilterExists {
							// Only attributes in the subfilter
							// should affect the result
							filterMatches[attributeName] = true
							break
						}
					}
					// There shouldn't be two indexes with the same ID
					break
				}
			}
		}
		allSubfiltersMatch := true
		for _, matches := range filterMatches {
			if !matches {
				allSubfiltersMatch = false
				break
			}
		}
		if allSubfiltersMatch {
			docMatches = append(docMatches, docId)
		}
	}
	return docMatches
}

// Returns all document IDs for which condition is met for any subfilter of given query operator
func FetchMatchesAny(edvId string, subfilters []map[string]string, operator string) []string {
	docMatches := make([]string, 0)
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return docMatches
	}
	uniqueDocMatches := make(map[string]bool)
	for _, subfilter := range subfilters {
		subfilterMatches := FetchMatchesAll(edvId, subfilter, operator)
		for _, match := range subfilterMatches {
			if !uniqueDocMatches[match] {
				uniqueDocMatches[match] = true
			}
		}
	}
	for docId := range uniqueDocMatches {
		docMatches = append(docMatches, docId)
	}
	return docMatches
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

	indexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
	indexFile, _ := os.Create(indexFileName)
	indexFileString := "{}"
	indexFile.WriteString(indexFileString)

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

// Search EDV with all query
func SearchEdvAll(edvId string, subfilter map[string]string, operator string, searchRequest EdvSearchRequest) []byte {
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return make([]byte, 0)
	}
	if searchRequest.ReturnFullDocuments {
		matches := FetchMatchesAll(edvId, subfilter, operator)
		fullMatches := GetDocumentsById(edvId, matches)
		fullMatchesBytes, _ := json.MarshalIndent(fullMatches, "", "  ")
		return fullMatchesBytes
	}
	matches := FetchMatchesAll(edvId, subfilter, operator)
	matchesBytes, _ := json.MarshalIndent(matches, "", "  ")
	return matchesBytes
}

// Search EDV with any query
func SearchEdvAny(edvId string, subfilters []map[string]string, operator string, searchRequest EdvSearchRequest) []byte {
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return make([]byte, 0)
	}
	if searchRequest.ReturnFullDocuments {
		matches := FetchMatchesAny(edvId, subfilters, operator)
		fullMatches := GetDocumentsById(edvId, matches)
		fullMatchesBytes, _ := json.MarshalIndent(fullMatches, "", "  ")
		return fullMatchesBytes
	}
	matches := FetchMatchesAny(edvId, subfilters, operator)
	matchesBytes, _ := json.MarshalIndent(matches, "", "  ")
	return matchesBytes
}

// Search EDV
func SearchEdv(res http.ResponseWriter, req *http.Request) {
	var edvSearchRequest EdvSearchRequest
	body, bodyReadErr := ioutil.ReadAll(req.Body)
	bodyUnmarshalErr := json.Unmarshal(body, &edvSearchRequest)

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
	if equalsAll := edvSearchRequest.EqualsAll; equalsAll != nil {
		matchesBytes := SearchEdvAll(edvId, equalsAll, EdvSearchOperators.Equals, edvSearchRequest)
		res.Write(matchesBytes)
		return
	}
	if equalsAny := edvSearchRequest.EqualsAny; equalsAny != nil {
		matchesBytes := SearchEdvAny(edvId, equalsAny, EdvSearchOperators.Equals, edvSearchRequest)
		res.Write(matchesBytes)
		return
	}
}
