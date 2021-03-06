package lmdrouter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// MarshalResponse generated an events.APIGatewayProxyResponse object that can
// be directly returned via the lambda's handler function. It receives an HTTP
// status code for the response, a map of HTTP headers (can be empty or nil),
// and a value (probably a struct) representing the response body. This value
// will be marshaled to JSON (currently without base 64 encoding).
func MarshalResponse(status int, headers map[string]string, data interface{}) (
	events.APIGatewayProxyResponse,
	error,
) {
	b, err := json.Marshal(data)
	if err != nil {
		status = http.StatusInternalServerError
		b = []byte(`{"code":500,"message":"the server has encountered an unexpected error"}`)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json; charset=UTF-8"

	return events.APIGatewayProxyResponse{
		StatusCode:      status,
		IsBase64Encoded: false,
		Headers:         headers,
		Body:            string(b),
	}, nil
}

// HandleError generates an events.APIGatewayProxyResponse from an error value.
// If the error is an HTTPError, the response's status code will be taken from
// the error. Otherwise, the error is assumed to be 500 Internal Server Error.
// Regardless, all errors will generate a JSON response in the format
// `{ "code": 500, "error": "something failed" }`
// This format cannot currently be changed.
func HandleError(err error) (events.APIGatewayProxyResponse, error) {
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		return MarshalResponse(httpErr.Code, nil, httpErr)
	}

	return MarshalResponse(
		http.StatusInternalServerError,
		nil,
		HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	)
}
