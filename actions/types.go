package actions

import "encoding/json"

// Encrypted index structure
type EncryptedIndex struct {
	// maps index IDs to doc IDs for quick search responses
	DocIds map[string][]string `json:"docIds"`
	// maps doc IDs to index IDs for quick update bookkeeping
	IndexIds map[string][]string `json:"indexIds"`
}

// Encrypted document operations
type EncryptedDocumentOperationOptions struct {
	Create string
	Update string
	Delete string
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
