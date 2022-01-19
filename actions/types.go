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
}

// Get EDV history log entry structure
type EdvHistoryLogEntry struct {
	DocumentId string `json:"documentId"`
	Sequence   uint64 `json:"sequence"`
	Operation  string `json:"operation"`
}
