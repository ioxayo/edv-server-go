package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/errors"
	"github.com/ioxayo/edv-server-go/storage"
)

// Update state of EDV
// TODO: may need to add global locking around this function to
// avoid inconsistent state from concurrent client updates
func UpdateEdvState(edvId string, docId string, operation string) errors.HttpError {
	if !common.IsValidEnumMember(EncryptedDocumentOperations, operation) {
		return errors.NilError()
	}

	// Retrieve and parse config
	var config common.DataVaultConfiguration
	// configFileName := fmt.Sprintf("./edvs/%s/config.json", edvId)
	// configFileBytes, _ := os.ReadFile(configFileName)
	configFileBytes, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.Config)
	json.Unmarshal(configFileBytes, &config)

	// Retrieve and parse history
	var history []EdvHistoryLogEntry
	// historyFileName := fmt.Sprintf("./edvs/%s/history.json", edvId)
	// historyFileBytes, _ := os.ReadFile(historyFileName)
	historyFileBytes, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.History)
	json.Unmarshal(historyFileBytes, &history)

	// Update parsed config
	config.Sequence++
	configFileBytes, _ = json.MarshalIndent(config, "", "  ")

	// Update parsed history
	historyEntry := EdvHistoryLogEntry{docId, config.Sequence, operation}
	history = append(history, historyEntry)
	historyFileBytes, _ = json.MarshalIndent(history, "", "  ")

	// Persist updated config and history
	// os.WriteFile(configFileName, configFileBytes, os.ModePerm)
	// os.WriteFile(historyFileName, historyFileBytes, os.ModePerm)

	// Persist updated config
	updateConfigErr := storage.Provider.UpdateDocSystem(edvId, storage.SystemFiles.Config, configFileBytes)
	if updateConfigErr.IsError() {
		return updateConfigErr
	}

	// Persist updated history
	updateHistoryErr := storage.Provider.UpdateDocSystem(edvId, storage.SystemFiles.History, historyFileBytes)
	if updateHistoryErr.IsError() {
		return updateHistoryErr
	}

	return errors.NilError()
}

// Update EDV index for create operation
// TODO: may need to add global locking around this function to
// avoid inconsistent state from concurrent client updates
func UpdateEdvIndexCreate(edvId string, doc common.EncryptedDocument) errors.HttpError {
	// Check if doc has index
	if docIndex := doc.Indexed; docIndex != nil {
		// Fetch index file
		var edvIndex EncryptedIndex
		// edvIndexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
		// edvIndexFileBytesBefore, _ := os.ReadFile(edvIndexFileName)
		edvIndexFileBytesBefore, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.Index)
		json.Unmarshal(edvIndexFileBytesBefore, &edvIndex)

		// Iterate through all doc indexes
		docId := doc.Id
		docIndexes := make([]string, 0)
		for _, index := range docIndex {
			// Update or initialize array for index-ID-keyed map with doc ID
			indexId := index.Hmac.Id
			if docIds, indexExists := edvIndex.DocIds[indexId]; indexExists {
				docIds = append(docIds, docId)
				edvIndex.DocIds[indexId] = docIds
			} else {
				edvIndex.DocIds[indexId] = []string{docId}
			}
			// Build array for doc-ID-keyed map with index IDs
			docIndexes = append(docIndexes, indexId)
		}

		// Bind index array to doc ID
		edvIndex.IndexIds[docId] = docIndexes

		// Update index file
		edvIndexFileBytesAfter, _ := json.MarshalIndent(edvIndex, "", "  ")
		// os.WriteFile(edvIndexFileName, edvIndexFileBytesAfter, os.ModePerm)
		return storage.Provider.UpdateDocSystem(edvId, storage.SystemFiles.Index, edvIndexFileBytesAfter)
	}
	return errors.NilError()
}

// Update EDV index for update operation
// TODO: may need to add global locking around this function to
// avoid inconsistent state from concurrent client updates
func UpdateEdvIndexUpdate(edvId string, doc common.EncryptedDocument) errors.HttpError {
	// Check if doc has index
	if docIndex := doc.Indexed; docIndex != nil {
		// Fetch index file
		var edvIndex EncryptedIndex
		// edvIndexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
		// edvIndexFileBytesBefore, _ := os.ReadFile(edvIndexFileName)
		edvIndexFileBytesBefore, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.Index)
		json.Unmarshal(edvIndexFileBytesBefore, &edvIndex)

		// Iterate through all doc indexes
		docId := doc.Id
		newDocIndexes := make([]string, 0)
		for _, index := range docIndex {
			// Update or initialize array for index-ID-keyed map with doc ID
			indexId := index.Hmac.Id
			if docIds, indexExists := edvIndex.DocIds[indexId]; indexExists {
				// Since this doc already exists, we should only
				// add it to index if it is not already tracked
				if isDocIndexed := common.IsValueInArray(docIds, docId); !isDocIndexed {
					docIds = append(docIds, docId)
					edvIndex.DocIds[indexId] = docIds
				} else {
					// Build array of index IDs that are not yet tracking this doc
					newDocIndexes = append(newDocIndexes, indexId)
				}
			} else {
				edvIndex.DocIds[indexId] = []string{docId}
			}
		}

		// Join existing array for doc-ID-keyed map
		// with newly discovered index IDs
		existingDocIndexes := edvIndex.IndexIds[docId]
		updatedDocIndexes := append(existingDocIndexes, newDocIndexes...)
		edvIndex.IndexIds[docId] = updatedDocIndexes

		// Update index file
		edvIndexFileBytesAfter, _ := json.MarshalIndent(edvIndex, "", "  ")
		// os.WriteFile(edvIndexFileName, edvIndexFileBytesAfter, os.ModePerm)
		return storage.Provider.UpdateDocSystem(edvId, storage.SystemFiles.Index, edvIndexFileBytesAfter)
	}
	return errors.NilError()
}

// Update EDV index for delete operation
// TODO: may need to add global locking around this function to
// avoid inconsistent state from concurrent client updates
func UpdateEdvIndexDelete(edvId string, docId string) errors.HttpError {
	// Fetch index file
	var edvIndex EncryptedIndex
	// edvIndexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
	// edvIndexFileBytesBefore, _ := os.ReadFile(edvIndexFileName)
	edvIndexFileBytesBefore, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.Index)
	json.Unmarshal(edvIndexFileBytesBefore, &edvIndex)

	// Retrieve index array for doc ID
	docIndexIds := edvIndex.IndexIds[docId]

	// Remove doc ID from all indexes
	for _, indexId := range docIndexIds {
		indexDocIds := edvIndex.DocIds[indexId]
		updatedIndexDocIds := common.RemoveValueFromArray(indexDocIds, docId)
		edvIndex.DocIds[indexId] = updatedIndexDocIds
	}

	// Remove indexes associated with doc ID
	delete(edvIndex.IndexIds, docId)

	// Update index file
	edvIndexFileBytesAfter, _ := json.MarshalIndent(edvIndex, "", "  ")
	// os.WriteFile(edvIndexFileName, edvIndexFileBytesAfter, os.ModePerm)
	return storage.Provider.UpdateDocSystem(edvId, storage.SystemFiles.Index, edvIndexFileBytesAfter)
}

// Retrieve document IDs associated with an index ID
func IndexToDocuments(edvId string, indexId string) []string {
	var index EncryptedIndex
	// indexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
	// indexFileBytes, _ := os.ReadFile(indexFileName)
	indexFileBytes, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.Index)
	json.Unmarshal(indexFileBytes, &index)
	return index.DocIds[indexId]
}

// Returns all document IDs for which condition is met for all key-value pairs of subfilter of given query operator
func FetchMatchesAll(edvId string, indexId string, subfilter map[string]string, operator string) []string {
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return make([]string, 0)
	}

	docIds := IndexToDocuments(edvId, indexId)
	docMatches := make([]string, 0)

	for _, docId := range docIds {
		// docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
		if docExists, _ := storage.Provider.DocExistsClient(edvId, docId); !docExists {
			continue
		}

		// docFileBytes, _ := os.ReadFile(docFileName)
		docFileBytes, _ := storage.Provider.ReadDocClient(edvId, docId)
		var doc common.EncryptedDocument
		if err := json.Unmarshal(docFileBytes, &doc); err != nil {
			continue
		}

		// Track subfilter matches in map initialized with keys in subfilter
		filterMatches := make(map[string]bool)
		for key, _ := range subfilter {
			filterMatches[key] = false
		}

		switch operator {
		case EdvSearchOperators.Equals:
			indexes := doc.Indexed
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

		// Check if all subfilter pairs match
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
func FetchMatchesAny(edvId string, indexId string, subfilters []map[string]string, operator string) []string {
	docMatches := make([]string, 0)
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return docMatches
	}
	uniqueDocMatches := make(map[string]bool)
	for _, subfilter := range subfilters {
		subfilterMatches := FetchMatchesAll(edvId, indexId, subfilter, operator)
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
	// var edvConfig common.DataVaultConfiguration
	body, bodyReadErr := ioutil.ReadAll(req.Body)

	if bodyReadErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", bodyReadErr)
		status := http.StatusBadRequest
		errors.HandleError(res, req, message, status)
		return
	}

	// var edvId string
	// if edvConfig.Id != "" {
	// 	edvId = edvConfig.Id
	// } else {
	// 	edvId = uuid.NewString()
	// }
	// edvConfig.Id = edvId

	// edvDirName := filepath.Join(".", "edvs", edvId)
	// docDirName := filepath.Join(edvDirName, "docs")
	// os.MkdirAll(docDirName, os.ModePerm)

	// configFileName := fmt.Sprintf("./edvs/%s/config.json", edvId)
	// configFile, _ := os.Create(configFileName)
	// configFileBytes, _ := json.MarshalIndent(edvConfig, "", "  ")
	// configFile.Write(configFileBytes)

	// historyFileName := fmt.Sprintf("./edvs/%s/history.json", edvId)
	// historyFile, _ := os.Create(historyFileName)
	// historyFileString := "[]"
	// historyFile.WriteString(historyFileString)

	// indexFileName := fmt.Sprintf("./edvs/%s/index.json", edvId)
	// indexFile, _ := os.Create(indexFileName)
	// indexFileString := "{\n  \"docIds\": {},\n  \"indexIds\": {}\n}"
	// indexFile.WriteString(indexFileString)

	// edvLocation := fmt.Sprintf("%s/edvs/%s", req.Host, edvId)

	edvLocation, createEdvErr := storage.Provider.CreateEdv(body)
	if createEdvErr.IsError() {
		message := createEdvErr.Message
		status := createEdvErr.Status
		errors.HandleError(res, req, message, status)
		return
	}
	res.Header().Add("Location", edvLocation)
	res.WriteHeader(http.StatusCreated)
}

// Get all EDVs
func GetEdvs(res http.ResponseWriter, req *http.Request) {}

// Get EDV
func GetEdv(res http.ResponseWriter, req *http.Request) {}

// Get history of EDV
func GetEdvHistory(res http.ResponseWriter, req *http.Request) {
	var history []EdvHistoryLogEntry
	var historyFiltered []EdvHistoryLogEntry

	afterSequenceString := req.URL.Query().Get("afterSequence")
	beforeSequenceString := req.URL.Query().Get("beforeSequence")
	var afterSequence uint64
	var beforeSequence uint64

	edvId := mux.Vars(req)["edvId"]
	// historyFileName := fmt.Sprintf("./edvs/%s/history.json", edvId)
	// historyFileBytes, _ := os.ReadFile(historyFileName)
	historyFileBytes, _ := storage.Provider.ReadDocSystem(edvId, storage.SystemFiles.History)
	json.Unmarshal(historyFileBytes, &history)

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

	for _, entry := range history {
		if entry.Sequence > afterSequence && entry.Sequence < beforeSequence {
			historyFiltered = append(historyFiltered, entry)
		}
	}
	historyFileBytesFiltered, _ := json.MarshalIndent(historyFiltered, "", "  ")
	res.Write(historyFileBytesFiltered)
}

// Search EDV with all query
func SearchEdvAll(edvId string, indexId string, subfilter map[string]string, operator string, searchRequest EdvSearchRequest) []byte {
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return make([]byte, 0)
	}
	if searchRequest.ReturnFullDocuments {
		matches := FetchMatchesAll(edvId, indexId, subfilter, operator)
		fullMatches := GetDocumentsById(edvId, matches)
		fullMatchesBytes, _ := json.MarshalIndent(fullMatches, "", "  ")
		return fullMatchesBytes
	}
	matches := FetchMatchesAll(edvId, indexId, subfilter, operator)
	matchesBytes, _ := json.MarshalIndent(matches, "", "  ")
	return matchesBytes
}

// Search EDV with any query
func SearchEdvAny(edvId string, indexId string, subfilters []map[string]string, operator string, searchRequest EdvSearchRequest) []byte {
	if !common.IsValidEnumMember(EdvSearchOperators, operator) {
		return make([]byte, 0)
	}
	if searchRequest.ReturnFullDocuments {
		matches := FetchMatchesAny(edvId, indexId, subfilters, operator)
		fullMatches := GetDocumentsById(edvId, matches)
		fullMatchesBytes, _ := json.MarshalIndent(fullMatches, "", "  ")
		return fullMatchesBytes
	}
	matches := FetchMatchesAny(edvId, indexId, subfilters, operator)
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
	indexId := edvSearchRequest.Index
	equals := edvSearchRequest.Equals
	var equalsAll map[string]string
	var equalsAny []map[string]string

	if equalsAllUnmarshalErr := json.Unmarshal(equals, &equalsAll); equalsAllUnmarshalErr == nil {
		matchesBytes := SearchEdvAll(edvId, indexId, equalsAll, EdvSearchOperators.Equals, edvSearchRequest)
		res.Write(matchesBytes)
		return
	}

	if equalsAnyUnmarshalErr := json.Unmarshal(equals, &equalsAny); equalsAnyUnmarshalErr == nil {
		matchesBytes := SearchEdvAny(edvId, indexId, equalsAny, EdvSearchOperators.Equals, edvSearchRequest)
		res.Write(matchesBytes)
		return
	}
}
