package errors

// Generic api response structure
type ErrorLogEntry struct {
	Id      string
	Url     string
	Method  string
	Message string
	Status  int
}
