package errors

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ioxayo/edv-server-go/common"
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
	resBytes, _ := json.Marshal(common.GenericResponse{message, false})
	res.Write(resBytes)
}
