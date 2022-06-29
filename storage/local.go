package storage

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/errors"
)

const (
	EDV_DIR      = "edvs"
	DOC_DIR      = "docs"
	CONFIG_FILE  = "config.json"
	HISTORY_FILE = "history.json"
	INDEX_FILE   = "index.json"
	STORAGE_FILE = "storage.json"
)

// Local storage config structure
type LocalStorageConfig struct {
	Type    string `json:"type"`    // type of storage provider (must be "local")
	EdvRoot string `json:"edvRoot"` // directory of edv data (defaults to current directory)
}

func InitLocalStorageProvider(edvRoot string, edvId string) StorageProvider {
	config := LocalStorageConfig{
		Type:    StorageProviderTypes.Local,
		EdvRoot: edvRoot,
	}
	docDir := config.GetDocsPath(edvId)
	os.MkdirAll(docDir, os.ModePerm)
	return config
}

func GetStorageFilePath(edvId string) (string, errors.HttpError) {
	currentDir, _ := os.Getwd()
	storageFile := fmt.Sprintf("%s/%s/%s/%s", currentDir, EDV_DIR, edvId, STORAGE_FILE)
	return storageFile, errors.NilError()
}

func GetLocalStorageProvider(edvId string) (StorageProvider, errors.HttpError) {
	var provider LocalStorageConfig
	storageFile, _ := GetStorageFilePath(edvId)
	storageFileBytes, _ := os.ReadFile(storageFile)
	providerUnmarshalErr := json.Unmarshal(storageFileBytes, &provider)
	if providerUnmarshalErr != nil {
		message := fmt.Sprintf("Error parsing provider: %v", providerUnmarshalErr)
		status := http.StatusInternalServerError
		providerParseError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return provider, providerParseError
	}
	return provider, errors.NilError()
}

func (provider LocalStorageConfig) GetEdvPath(edvId string) (string, errors.HttpError) {
	edvDir := fmt.Sprintf("%s/%s/%s", provider.EdvRoot, EDV_DIR, edvId)
	if _, err := os.Stat(edvDir); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find EDV with ID '%s'", edvId)
		status := http.StatusBadRequest
		missingEdvError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return edvDir, missingEdvError
	}
	return edvDir, errors.NilError()
}

func (provider LocalStorageConfig) GetDocsPath(edvId string) string {
	edvDir, err := provider.GetEdvPath(edvId)
	docDir := fmt.Sprintf("%s/%s", edvDir, DOC_DIR)
	if err.IsError() {
		return docDir
	}
	return docDir
}

func (provider LocalStorageConfig) GetDocPath(edvId string, docId string) string {
	docDir := provider.GetDocsPath(edvId)
	docFileName := fmt.Sprintf("%s/%s.json", docDir, docId)
	return docFileName
}

func (provider LocalStorageConfig) GetDoc(edvId string, docId string) (string, errors.HttpError) {
	docDir := provider.GetDocsPath(edvId)
	docFileName := fmt.Sprintf("%s/%s.json", docDir, docId)
	if _, err := os.Stat(docFileName); goerrors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Could not find document with ID '%s' in EDV with ID '%s'", docId, edvId)
		status := http.StatusBadRequest
		missingDocError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return docFileName, missingDocError
	}
	return docFileName, errors.NilError()
}

func (provider LocalStorageConfig) GetSysFile(edvId string, fileType string) (string, errors.HttpError) {
	var sysFile string
	edvDir, err := provider.GetEdvPath(edvId)
	switch fileType {
	case SystemFiles.Config:
		sysFile = fmt.Sprintf("%s/%s", edvDir, CONFIG_FILE)
	case SystemFiles.History:
		sysFile = fmt.Sprintf("%s/%s", edvDir, HISTORY_FILE)
	case SystemFiles.Index:
		sysFile = fmt.Sprintf("%s/%s", edvDir, INDEX_FILE)
	case SystemFiles.Storage:
		sysFile = fmt.Sprintf("%s/%s", edvDir, STORAGE_FILE)
	default:
		sysFile = fmt.Sprintf("%s/%s", edvDir, fileType)
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return sysFile, invalidFileTypeError
	}
	if err.IsError() {
		return sysFile, err
	}
	return sysFile, errors.NilError()
}

func (provider LocalStorageConfig) CreateDocClient(edvId string, docId string, data []byte) (string, errors.HttpError) {
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

	docFileName := provider.GetDocPath(edvId, docId)
	docFile, _ := os.Create(docFileName)
	docFileBytes, _ := json.MarshalIndent(doc, "", "  ")
	docFile.Write(docFileBytes)

	docLocation := provider.GetDocPath(edvId, docId)
	return docLocation, errors.NilError()
}

func (provider LocalStorageConfig) CreateDocSystem(edvId string, fileType string, data []byte) errors.HttpError {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return invalidFileTypeError
	}
	switch fileType {
	case SystemFiles.Config:
		var edvConfig common.DataVaultConfiguration
		json.Unmarshal(data, &edvConfig)
		configFileName, _ := provider.GetSysFile(edvId, SystemFiles.Config)
		configFile, _ := os.Create(configFileName)
		configFileBytes, _ := json.MarshalIndent(edvConfig, "", "  ")
		configFile.Write(configFileBytes)
	default:
		sysFileName, _ := provider.GetSysFile(edvId, fileType)
		sysFile, _ := os.Create(sysFileName)
		sysFile.Write(data)
	}
	return errors.NilError()
}

func (provider LocalStorageConfig) ReadDocClient(edvId string, docId string) ([]byte, errors.HttpError) {
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

func (provider LocalStorageConfig) ReadDocSystem(edvId string, fileType string) ([]byte, errors.HttpError) {
	if !common.IsValidEnumMember(SystemFiles, fileType) {
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return make([]byte, 0), invalidFileTypeError
	}

	edvDir, err := provider.GetEdvPath(edvId)
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

func (provider LocalStorageConfig) UpdateDocClient(edvId string, docId string, data []byte) errors.HttpError {
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
	return errors.NilError()
}

func (provider LocalStorageConfig) UpdateDocSystem(edvId string, fileType string, data []byte) errors.HttpError {
	switch fileType {
	case SystemFiles.Config:
		configFileName, _ := provider.GetSysFile(edvId, SystemFiles.Config)
		os.WriteFile(configFileName, data, os.ModePerm)
	case SystemFiles.History:
		historyFileName, _ := provider.GetSysFile(edvId, SystemFiles.History)
		os.WriteFile(historyFileName, data, os.ModePerm)
	case SystemFiles.Index:
		indexFileName, _ := provider.GetSysFile(edvId, SystemFiles.Index)
		os.WriteFile(indexFileName, data, os.ModePerm)
	default:
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return invalidFileTypeError
	}
	return errors.NilError()
}

func (provider LocalStorageConfig) DeleteDocClient(edvId string, docId string) errors.HttpError {
	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return err
	}

	os.Remove(docFileName)
	return errors.NilError()
}

func (provider LocalStorageConfig) DeleteDocSystem(edvId string, fileType string) errors.HttpError {
	switch fileType {
	case SystemFiles.Config:
		configFileName, _ := provider.GetSysFile(edvId, SystemFiles.Config)
		os.Remove(configFileName)
	case SystemFiles.History:
		historyFileName, _ := provider.GetSysFile(edvId, SystemFiles.History)
		os.Remove(historyFileName)
	case SystemFiles.Index:
		indexFileName, _ := provider.GetSysFile(edvId, SystemFiles.Index)
		os.Remove(indexFileName)
	default:
		message := fmt.Sprintf("Invalid file type '%s'", fileType)
		status := http.StatusBadRequest
		invalidFileTypeError := errors.HttpError{
			Message: message,
			Status:  status,
		}
		return invalidFileTypeError
	}
	return errors.NilError()
}

func (provider LocalStorageConfig) DocExistsClient(edvId string, docId string) (bool, errors.HttpError) {
	docFileName, err := provider.GetDoc(edvId, docId)
	if err.IsError() {
		return false, err
	}
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
