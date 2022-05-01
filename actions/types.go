package actions

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
			Unique string `json:"unique"`
		} `json:"attributes"`
	} `json:"indexed,omitempty"`
}

// Encrypted document operations
type EncryptedDocumentOperationOptions struct {
	Created string `json:"created"`
	Updated string `json:"updated"`
	Deleted string `json:"deleted"`
}

var (
	EncryptedDocumentOperations = EncryptedDocumentOperationOptions{
		Created: "created",
		Updated: "updated",
		Deleted: "deleted",
	}
)

// Get EDV history log entry structure
type EdvHistoryLogEntry struct {
	DocumentId string `json:"documentId"`
	Sequence   uint64 `json:"sequence"`
	Operation  string `json:"operation"`
}

// EDV search reuest body structure
type EdvSearchRequest struct {
	Index               string              `json:"index"`
	EqualsAll           map[string]string   `json:"equals"`
	EqualsAny           []map[string]string `json:"equals"`
	ReturnFullDocuments bool                `json:"returnFullDocuments"`
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
