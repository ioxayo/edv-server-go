package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/ioxayo/edv-server-go/common"
	"github.com/ioxayo/edv-server-go/storage"
)

func BeforeAllTests() {
	currentDir, _ := os.Getwd()
	edvRoot := fmt.Sprintf("%s/%s", currentDir, "tmp")
	os.Setenv(common.EnvVars.StorageType, storage.StorageProviderTypes.Local)
	os.Setenv(common.EnvVars.StorageLocalRoot, edvRoot)
	os.Setenv(common.EnvVars.Port, "5000")
}

func AfterAllTests() {
	os.Setenv(common.EnvVars.StorageType, "")
	os.Setenv(common.EnvVars.StorageLocalRoot, "")
	os.Setenv(common.EnvVars.Port, "")
}

func TestMain(m *testing.M) {
	BeforeAllTests()
	code := m.Run()
	AfterAllTests()
	os.Exit(code)
}
