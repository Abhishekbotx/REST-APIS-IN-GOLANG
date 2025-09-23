package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"` // by adding these tags in response u ll get status instead of Status
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	//Responsewriter is applying same io.write type so we can use that in NewEncoder coz thats also of NewEncoder type
	w.Header().Set("Content-Type", "appliaction/json")
	w.WriteHeader(status)
	//STRUCT DATA TO JSON IS encode and json data to struct is decode

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {

	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}

}

func ValidationError(errs validator.ValidationErrors) Response { //this validator.ValidationErrors is a slice
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required files", err.Field())) //field prop will give the error
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid ", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs,", "),//slice and sepeerator
	}
}
