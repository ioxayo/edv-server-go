package storage

import (
	"os"

	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/errors"
)

func GetStorageProvider(edvId string) (StorageProvider, errors.HttpError) {
	var provider StorageProvider
	var providerErr errors.HttpError
	switch os.Getenv(common.EnvVars.StorageType) {
	case StorageProviderTypes.Local:
	default:
		provider, providerErr = GetLocalStorageProvider(edvId)
	}
	return provider, providerErr
}
