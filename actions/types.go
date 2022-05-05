package actions

import "encoding/json"

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

// Encrypted index structure
type EncryptedIndex struct {
	// maps index IDs to doc IDs for quick search responses
	DocIds map[string][]string `json:"docIds"`
	// maps doc IDs to index IDs for quick update bookkeeping
	IndexIds map[string][]string `json:"indexIds"`
}

// Encrypted document operations
type EncryptedDocumentOperationOptions struct {
	Create string `json:"created"`
	Update string `json:"updated"`
	Delete string `json:"deleted"`
}

var (
	EncryptedDocumentOperations = EncryptedDocumentOperationOptions{
		Create: "created",
		Update: "updated",
		Delete: "deleted",
	}
)

func (enumStruct EncryptedDocumentOperationOptions) IsEnum() bool {
	return true
}

// Get EDV history log entry structure
type EdvHistoryLogEntry struct {
	DocumentId string `json:"documentId"`
	Sequence   uint64 `json:"sequence"`
	Operation  string `json:"operation"`
}

// EDV search reuest body structure
type EdvSearchRequest struct {
	Index               string          `json:"index"`
	Equals              json.RawMessage `json:"equals"` // map[string]string | []map[string]string
	ReturnFullDocuments bool            `json:"returnFullDocuments"`
}

// Search operators supported by this implementation of EDV
type EdvSearchOperatorOptions struct {
	Equals string `json:"equals"`
}

var (
	EdvSearchOperators = EdvSearchOperatorOptions{
		Equals: "equals",
	}
)

func (enumStruct EdvSearchOperatorOptions) IsEnum() bool {
	return true
}
