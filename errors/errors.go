package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func LogError(req *http.Request, message string, status int) {
	errorLog := ErrorLogEntry{
		Id:      uuid.NewString(),
		Url:     req.RequestURI,
		Method:  req.Method,
		Message: message,
		Status:  status,
	}
	log.Println(errorLog)
}

func HandleError(res http.ResponseWriter, req *http.Request, message string, status int) {
	LogError(req, message, status)
	res.WriteHeader(status)
	resBytes, _ := json.Marshal(HttpError{message, status})
	res.Write(resBytes)
}

func NilError() HttpError {
	return HttpError{Message: "", Status: http.StatusOK}
}

func (err HttpError) IsError() bool {
	return err.Status >= 400
}

func (err HttpError) Error() string {
	return fmt.Sprintf("%d error: %s", err.Status, err.Message)
}
