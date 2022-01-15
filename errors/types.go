package errors

// Generic api response structure
type ErrorLog struct {
	Id      string
	Url     string
	Method  string
	Message string
	Status  int
}
