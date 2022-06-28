package errors

// Generic api response structure
type ErrorLogEntry struct {
	Id      string
	Url     string
	Method  string
	Message string
	Status  int
}

// Simple error structure for resource operations
type HttpError struct {
	Message string
	Status  int
}
