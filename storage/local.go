package storage

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/errors"
)

const (
	EDV_DIR      = "edvs"
	DOC_DIR      = "docs"
	CONFIG_FILE  = "config.json"
	HISTORY_FILE = "history.json"
	INDEX_FILE   = "index.json"
)

// Local storage config structure
type LocalStorageConfig struct {
	EdvHost string // url of edv service
	EdvRoot string // directory of edv data (defaults to current directory)
	// EdvDir  string // defaults to 'edvs'
	// DocDir  string // defaults to 'docs'
}

func InitLocalStorageProvider(edvHost string, edvRoot string) StorageProvider {
	config := LocalStorageConfig{
		EdvHost: edvHost,
		EdvRoot: edvRoot,
	}
	return &config
}

func (provider *LocalStorageConfig) GetEdvDir(edvId string) (string, errors.HttpError) {
	edvDir := fmt.Sprintf("%s/%s/%s", provider.EdvRoot, EDV_DIR, edvId)
	if _, err := os.Stat(edvDir); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find EDV with ID '%s'", edvId)
		status := http.StatusBadRequest
		missingEdvError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return "", missingEdvError
	}
	return edvDir, errors.NilError()
}

func (provider *LocalStorageConfig) GetDocDir(edvId string) (string, errors.HttpError) {
	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return "", err
	}
	docDir := fmt.Sprintf("%s/%s", edvDir, DOC_DIR)
	return docDir, errors.NilError()
}

func (provider *LocalStorageConfig) GetDoc(edvId string, docId string) (string, errors.HttpError) {
	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return "", err
	}
	docDir := fmt.Sprintf("%s/%s", edvDir, DOC_DIR)
	docFileName := fmt.Sprintf("%s/%s.json", docDir, docId)
	if docExists, err := provider.DocExistsClient(edvId, docId); !docExists || err.IsError() {
		return "", err
	}
	return docFileName, errors.NilError()
}

func (provider *LocalStorageConfig) GetSysFile(edvId string, fileType string) (string, errors.HttpError) {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return "", invalidFileTypeError
	}

	var sysFile string
	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return "", err
	}
	switch fileType {
	case SystemFiles.Config:
		sysFile = fmt.Sprintf("%s/%s", edvDir, CONFIG_FILE)
	case SystemFiles.History:
		sysFile = fmt.Sprintf("%s/%s", edvDir, HISTORY_FILE)
	case SystemFiles.Index:
		sysFile = fmt.Sprintf("%s/%s", edvDir, INDEX_FILE)
	}
	return sysFile, errors.NilError()
}

func (provider *LocalStorageConfig) CreateEdv(data []byte) (string, errors.HttpError) {
	var edvConfig common.DataVaultConfiguration
	dataUnmarshalErr := json.Unmarshal(data, &edvConfig)
	if dataUnmarshalErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", dataUnmarshalErr)
		status := http.StatusBadRequest
		bodyParseError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return "", bodyParseError
	}

	var edvId string
	if edvConfig.Id != "" {
		edvId = edvConfig.Id
	} else {
		edvId = uuid.NewString()
	}
	edvConfig.Id = edvId

	docDir := fmt.Sprintf("%s/%s/%s/%s", provider.EdvRoot, EDV_DIR, edvId, DOC_DIR)
	os.MkdirAll(docDir, os.ModePerm)

	configFileName, _ := provider.GetSysFile(edvId, SystemFiles.Config)
	configFile, _ := os.Create(configFileName)
	configFileBytes, _ := json.MarshalIndent(edvConfig, "", "  ")
	configFile.Write(configFileBytes)

	historyFileName, _ := provider.GetSysFile(edvId, SystemFiles.History)
	historyFile, _ := os.Create(historyFileName)
	historyFileString := "[]"
	historyFile.WriteString(historyFileString)

	indexFileName, _ := provider.GetSysFile(edvId, SystemFiles.Index)
	indexFile, _ := os.Create(indexFileName)
	indexFileString := "{\n  \"docIds\": {},\n  \"indexIds\": {}\n}"
	indexFile.WriteString(indexFileString)

	edvLocation := fmt.Sprintf("%s/%s/%s", provider.EdvHost, EDV_DIR, edvId)
	return edvLocation, errors.NilError()
}

func (provider *LocalStorageConfig) CreateDocClient(edvId string, docId string, data []byte) (string, errors.HttpError) {
	var doc common.EncryptedDocument
	dataUnmarshalErr := json.Unmarshal(data, &doc)
	if dataUnmarshalErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", dataUnmarshalErr)
		status := http.StatusBadRequest
		bodyParseError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return "", bodyParseError
	}

	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return "", err
	}
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(doc, "", "  ")
	docFile.Write(docFileBytes)
	// actions.UpdateEdvState(edvId, docId, actions.EncryptedDocumentOperations.Create)
	// actions.UpdateEdvIndexCreate(edvId, doc)

	docLocation := fmt.Sprintf("%s/%s/%s/%s/%s", provider.EdvHost, EDV_DIR, edvId, DOC_DIR, docId)
	return docLocation, errors.NilError()
}

func (provider *LocalStorageConfig) CreateDocSystem(edvId string, fileType string, data []byte) errors.HttpError {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return invalidFileTypeError
	}

	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return err
	}

	switch fileType {
	case SystemFiles.Config:
		var edvConfig common.DataVaultConfiguration
		json.Unmarshal(data, &edvConfig)
		configFileName := fmt.Sprintf("%s/config.json", edvDir)
		configFile, _ := os.Create(configFileName)
		configFileBytes, _ := json.MarshalIndent(edvConfig, "", "  ")
		configFile.Write(configFileBytes)
	case SystemFiles.History:
		historyFileName := fmt.Sprintf("%s/history.json", edvDir)
		historyFile, _ := os.Create(historyFileName)
		historyFileString := "[]"
		historyFile.WriteString(historyFileString)
	case SystemFiles.Index:
		indexFileName := fmt.Sprintf("%s/index.json", edvDir)
		indexFile, _ := os.Create(indexFileName)
		indexFileString := "{\n  \"docIds\": {},\n  \"indexIds\": {}\n}"
		indexFile.WriteString(indexFileString)
	}
	return errors.NilError()
}

func (provider *LocalStorageConfig) ReadDocClient(edvId string, docId string) ([]byte, errors.HttpError) {
	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return make([]byte, 0), err
	}

	docFileBytes, docFileErr := os.ReadFile(docFileName)

	if docFileErr != nil {
		message := fmt.Sprintf("Error retrieving document: %s", docFileErr.Error())
		status := http.StatusInternalServerError
		docFileError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return make([]byte, 0), docFileError
	}

	return docFileBytes, errors.NilError()
}

func (provider *LocalStorageConfig) ReadDocSystem(edvId string, fileType string) ([]byte, errors.HttpError) {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return make([]byte, 0), invalidFileTypeError
	}

	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return make([]byte, 0), err
	}

	fileName := fmt.Sprintf("%s/%s.json", edvDir, fileType)
	fileBytes, fileReadErr := os.ReadFile(fileName)

	if fileReadErr != nil {
		message := fmt.Sprintf("Error retrieving document: %s", fileReadErr.Error())
		status := http.StatusInternalServerError
		docFileError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return make([]byte, 0), docFileError
	}

	return fileBytes, errors.NilError()
}

func (provider *LocalStorageConfig) UpdateDocClient(edvId string, docId string, data []byte) errors.HttpError {
	var doc common.EncryptedDocument
	dataUnmarshalErr := json.Unmarshal(data, &doc)
	if dataUnmarshalErr != nil {
		message := fmt.Sprintf("Error parsing request body: %v", dataUnmarshalErr)
		status := http.StatusBadRequest
		bodyParseError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return bodyParseError
	}

	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return err
	}
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(doc, "", "  ")
	docFile.Write(docFileBytes)
	// actions.UpdateEdvState(edvId, docId, actions.EncryptedDocumentOperations.Update)
	// actions.UpdateEdvIndexUpdate(edvId, doc)
	return errors.NilError()
}

func (provider *LocalStorageConfig) UpdateDocSystem(edvId string, fileType string, data []byte) errors.HttpError {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return invalidFileTypeError
	}

	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return err
	}

	switch fileType {
	case SystemFiles.Config:
		configFileName := fmt.Sprintf("%s/config.json", edvDir)
		os.WriteFile(configFileName, data, os.ModePerm)
	case SystemFiles.History:
		historyFileName := fmt.Sprintf("%s/history.json", edvDir)
		os.WriteFile(historyFileName, data, os.ModePerm)
	case SystemFiles.Index:
		indexFileName := fmt.Sprintf("%s/index.json", edvDir)
		os.WriteFile(indexFileName, data, os.ModePerm)
	}
	return errors.NilError()
}

func (provider *LocalStorageConfig) DeleteDocClient(edvId string, docId string) errors.HttpError {
	// docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return err
	}

	os.Remove(docFileName)
	// UpdateEdvState(edvId, docId, EncryptedDocumentOperations.Delete)
	// UpdateEdvIndexDelete(edvId, docId)
	return errors.NilError()
}

func (provider *LocalStorageConfig) DeleteDocSystem(edvId string, fileType string) errors.HttpError {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return invalidFileTypeError
	}

	edvDir, err := provider.GetEdvDir(edvId)
	if err.IsError() {
		return err
	}

	switch fileType {
	case SystemFiles.Config:
		configFileName := fmt.Sprintf("%s/config.json", edvDir)
		os.Remove(configFileName)
	case SystemFiles.History:
		historyFileName := fmt.Sprintf("%s/history.json", edvDir)
		os.Remove(historyFileName)
	case SystemFiles.Index:
		indexFileName := fmt.Sprintf("%s/index.json", edvDir)
		os.Remove(indexFileName)
	}
	return errors.NilError()
}

func (provider *LocalStorageConfig) DocExistsClient(edvId string, docId string) (bool, errors.HttpError) {
	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return false, err
	}
	// docFileName := fmt.Sprintf("./edvs/%s/docs/%s.json", edvId, docId)
	if _, err := os.Stat(docFileName); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find document with ID '%s' in EDV with ID '%s'", docId, edvId)
		status := http.StatusBadRequest
		missingDocError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return false, missingDocError
	}
	return true, errors.NilError()
}
