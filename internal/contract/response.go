package contract

import "net/http"

//Response general way to mannage the API reponses
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//NewOKResponse returns a Response object. Automatically set the status to http.StatusOK (200 code)
func NewOKResponse(message string) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: message,
	}
}

// ErrorResponse common error response for all of the API`s endpoints
type ErrorResponse struct {
	Response
}

//NewErrorResponse ErrorResponse constructor
func NewErrorResponse(status int, message string) *ErrorResponse {
	return &ErrorResponse{
		Response{
			Status:  status,
			Message: message,
		},
	}
}
