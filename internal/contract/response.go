package contract

//Response used to represent a http response object
type Response struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}
