package storage

import "github.com/ioxayo/edv-server-go/errors"

// Storage provider structure
type StorageProvider interface {
	CreateEdv(data []byte) (string, errors.HttpError)
	CreateDocClient(edvId string, docId string, data []byte) (string, errors.HttpError)
	CreateDocSystem(edvId string, fileType string, data []byte) errors.HttpError
	ReadDocClient(edvId string, docId string) ([]byte, errors.HttpError)
	ReadDocSystem(edvId string, fileType string) ([]byte, errors.HttpError)
	UpdateDocClient(edvId string, docId string, data []byte) errors.HttpError
	UpdateDocSystem(edvId string, fileType string, data []byte) errors.HttpError
	DeleteDocClient(edvId string, docId string) errors.HttpError
	DeleteDocSystem(edvId string, fileType string) errors.HttpError
	DocExistsClient(edvId string, docId string) (bool, errors.HttpError)
}

// System file types
type SystemFileOptions struct {
	Config  string
	History string
	Index   string
}

var (
	SystemFiles = SystemFileOptions{
		Config:  "config",
		History: "history",
		Index:   "index",
	}
)

func (enumStruct SystemFileOptions) IsEnum() bool {
	return true
}
