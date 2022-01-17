package actions

// Create EDV request structure
type CreateEdvRequest struct {
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

// Create document request structure
type CreateDocumentRequest struct {
	Id       string      `json:"id"`
	Sequence uint64      `json:"sequence"`
	Jwe      interface{} `json:"jwe"`
}
