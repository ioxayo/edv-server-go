# Go EDV Server

## Introduction
This is a Go implementation of the server interface defined in the [*Encrypted Data Vaults*](https://identity.foundation/edv-spec) specification, which "describes a privacy-preserving mechanism for storing, indexing, and retrieving encrypted data at a storage provider." This is a useful service for individuals and organizations that wish to manage private data at a remote storage provider without revealing the data to them in the clear. This convenience minimizes the risk of storage providers exploiting sensitive data and provides a portable storage layer that can be migrated between compatible providers.

## Routes
There are several routes defined in the EDV specification. These are the routes that we support in this library:
- `POST /edvs` - create EDV
- `GET /edvs` - get all EDVs
- `GET /edvs/{edvId}` - get single EDV
- `POST /edvs/{edvId}/query` - search for EDV
- `POST /edvs/{edvId}/docs` - create document
- `GET /edvs/{edvId}/docs` - get all documents
- `GET /edvs/{edvId}/docs/{docId}` - get single document
- `POST /edvs/{edvId}/docs/{docId}` - update document
- `DELETE /edvs/{edvId}/docs/{docId}` - delete document
- `GET /edvs/{edvId}/history` - get EDV history

## `StorageProvider`
We acknowledge that EDV providers may wish to configure the API service independent of the storage layer. The `StorageProvider` interface offers this convenience. Please follow these steps if you wish to implement this interface for a new storage provider:
1. Create a new file in the storage package and implement the methods in `StorageProvider` defined below
2. Add a new `case` block for the new storage provider in the `switch` statement in `CreateEdv` (`actions/edvs.go`) that configures the appropriate storage provider based on the storage type
3. Add a new `case` block for the new storage provider in the `switch` statement in `GetStorageProvider` (`storage/utils.go`) that retrieves the appropriate storage provider based on the storage type
4. Define a new environment variable for the new storage type in `StorageProviderTypes` (`storage/types.go`)
5. Define all necessary environment variables for the new storage type in `EnvVars` (`common/types.go`)

Below are the required methods of `StorageProvider` (*Note: At the time of this writing, we have implemented a local storage implementation, which hosts the storage layer in the local filesystem of the same machine as the EDV service*):
- `CreateDocClient(edvId string, docId string, data []byte) (string, errors.HttpError)`: creates an encrypted doc and returns location and error (if any)
- `CreateDocSystem(edvId string, fileType string, data []byte) errors.HttpError`: creates a system doc (e.g., config, index) and returns error (if any)
- `ReadDocClient(edvId string, docId string) ([]byte, errors.HttpError)`: retrieves an encrypted doc and returns data and error (if any)
- `ReadDocSystem(edvId string, fileType string) ([]byte, errors.HttpError)`: retrieves a system doc and returns data and error (if any)
- `UpdateDocClient(edvId string, docId string, data []byte) errors.HttpError`: updates an encrypted doc and returns error (if any)
- `UpdateDocSystem(edvId string, fileType string, data []byte) errors.HttpError`: updates a system doc and returns error (if any)
- `DeleteDocClient(edvId string, docId string) errors.HttpError`: deletes an encrypted doc and returns error (if any)
- `DeleteDocSystem(edvId string, fileType string) errors.HttpError`: deletes a system doc and returns error (if any)
- `DocExistsClient(edvId string, docId string) (bool, errors.HttpError)`: checks if an encrypted doc exists

## Run
We have provided two commands to run the EDV server:
1. Local: `./bin/sys/run.sh -l`
2. Docker: `./bin/sys/run.sh -d`

## Test
We have provided two ways to test the functionality of the EDV server:
1. Run sample scripts in the `bin` folder (note relevant command-line inputs)
2. Run `go test -v` to run all tests located in files ending in `_test.go`

## Contributions
To make a contribution, please do one of the following:
- For code, submit a pull request against this repo via a fork or branch
- For questions, recommendations, and bug reports, create an issue in this repo
