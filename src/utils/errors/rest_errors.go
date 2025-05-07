package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)


type RestErr struct {
	ErrMessage string        `json:"message"`
	ErrStatus  int           `json:"status"`
	ErrError   string        `json:"error"`
	ErrCauses  []interface{} `json:"causes"`
}


func NewRestError(message string, status int, err string, causes []interface{}) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes
	}
}


func NewRestErrorFromBytes(bytes []byte) (*RestErr, error) {
	var apiErr RestErr
	
	err := json.Unmarshal(bytes, &apiErr)
	
	if (err != nil) {
		return nil, errors.New("Invalid json")
	}
	
	return &apiErr, nil
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request"
	}
}


func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found"
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized"
	}
}


func NewInternalServerError(message string, err error) *RestErr {
	result := &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error"
	}
	
	if (err != nil) {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	
	return result
}
