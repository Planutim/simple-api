package jsonerrors

//GenericError wraps error into JSON object
type GenericError struct {
	Message string `json:"message"`
}
