package common

type EnumStruct interface {
	IsEnum() bool
}

// Data vault configuration structure
type DataVaultConfiguration struct {
	Id                string   `json:"id,omitempty"`
	Sequence          uint64   `json:"sequence"`
	Controller        string   `json:"controller"`
	InvokerSingle     string   `json:"invoker,omitempty"`
	InvokerMultiple   []string `json:"invoker,omitempty"`
	DelegatorSingle   string   `json:"delegator,omitempty"`
	DelegatorMultiple []string `json:"delegator,omitempty"`
	ReferenceId       string   `json:"referenceId,omitempty"`
	KeyAgreementKey   struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"keyAgreementKey"`
	Hmac struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"hmac"`
}

// Encrypted document structure
type EncryptedDocument struct {
	Id       string      `json:"id"`
	Sequence uint64      `json:"sequence"`
	Jwe      interface{} `json:"jwe"`
	Indexed  []struct {
		Sequence uint64 `json:"sequence"`
		Hmac     struct {
			Id   string `json:"id"`
			Type string `json:"type"`
		} `json:"hmac"`
		Attributes []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Unique bool   `json:"unique"`
		} `json:"attributes"`
	} `json:"indexed,omitempty"`
}

// Environment variables
type EnvVarOptions struct {
	Host             string
	StorageType      string
	StorageLocalRoot string
}

var (
	EnvVars = EnvVarOptions{
		Host:             "HOST",
		StorageType:      "STORAGE_TYPE",
		StorageLocalRoot: "STORAGE_LOCAL_ROOT",
	}
)
