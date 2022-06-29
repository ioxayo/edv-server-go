# Go EDV Server

## StorageProvider
We respect that different EDV providers may wish to configure the API service independent of the storage layer. The `StorageProvider` interface offers this convenience. Below are the required methods of `StorageProvider`. (*Note: At the time of this writing, we have implemented a local storage implementation, which hosts the storage layer in the local filesystem of the same machine as the EDV service*):
- `CreateDocClient(edvId string, docId string, data []byte) (string, errors.HttpError)`: creates an encrypted doc and returns location and error (if any)
- `CreateDocSystem(edvId string, fileType string, data []byte) errors.HttpError`: creates a system doc (e.g., config, index) and returns error (if any)
- `ReadDocClient(edvId string, docId string) ([]byte, errors.HttpError)`: retrieves an encrypted doc and returns data and error (if any)
- `ReadDocSystem(edvId string, fileType string) ([]byte, errors.HttpError)`: retrieves a system doc and returns data and error (if any)
- `UpdateDocClient(edvId string, docId string, data []byte) errors.HttpError`: updates an encrypted doc and returns error (if any)
- `UpdateDocSystem(edvId string, fileType string, data []byte) errors.HttpError`: updates a system doc and returns error (if any)
- `DeleteDocClient(edvId string, docId string) errors.HttpError`: deletes an encrypted doc and returns error (if any)
- `DeleteDocSystem(edvId string, fileType string) errors.HttpError`: deletes a system doc and returns error (if any)
- `DocExistsClient(edvId string, docId string) (bool, errors.HttpError)`: checks if an encrypted doc exists
