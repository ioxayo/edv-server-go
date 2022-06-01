package storage

// Storage provider structure
type StorageProvider interface {
	CreateEdv(edvId string, data []byte) (string, error)
	CreateDocClient(edvId string, docId string, data []byte) (string, error)
	CreateDocSystem(edvId string, fileType string, data []byte) error
	ReadDocClient(edvId string, docId string) ([]byte, error)
	ReadDocSystem(edvId string, fileType string) ([]byte, error)
	UpdateDocClient(edvId string, docId string, data []byte) error
	UpdateDocSystem(edvId string, fileType string, data []byte) error
	DeleteDocClient(edvId string, docId string) error
	DeleteDocSystem(edvId string, fileType string) error
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
