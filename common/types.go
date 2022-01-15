package common

// Generic api response structure
type GenericResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
