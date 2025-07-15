package xresp

import "github.com/danielgtaylor/huma/v2"

// GeneralResponse is a generic standardized response format
// that wraps API responses with a code, message, optional data,
// error details, and a trace identifier for observability.
type GeneralResponse[D any, E any] struct {
	// Code is the application-specific status code.
	// Can be an integer, string, or other depending on implementation.
	Code any `json:"code" example:"200" format:"any" doc:"Application-level status code"`

	// Msg is a human-readable message describing the result.
	Msg string `json:"msg" example:"Success" doc:"Human-readable status message"`

	// Data contains the response payload when the request is successful.
	// This field is optional.
	Data D `json:"data" doc:"Response payload"`

	// Err contains error information when the request fails.
	// This field is optional.
	Err E `json:"err" doc:"Detailed error info when applicable"`

	// TraceID is used for tracing logs and request tracking across services.
	// It can be a string or other trace identifier format.
	TraceID string `json:"traceId,omitempty" example:"d1qlgrkqihv7a7754j20" doc:"Tracing identifier (xid) for debugging/logging"`
}

type ErrorModel struct {
	// Title provides a short static summary of the problem. Huma will default this
	// to the HTTP response status code text if not present.
	Title string `json:"title,omitempty" example:"Bad Request" doc:"A short, human-readable summary of the problem type. This value should not change between occurrences of the error."`

	// Detail is an explanation specific to this error occurrence.
	Detail string `json:"detail,omitempty" example:"Property foo is required but is missing." doc:"A human-readable explanation specific to this occurrence of the problem."`

	// Errors provides an optional mechanism of passing additional error details
	// as a list.
	Errors []*huma.ErrorDetail `json:"errors,omitempty" doc:"Optional list of individual error details"`
}

type GeneralResponseError struct {
	// Code is the application-specific status code.
	// Can be an integer, string, or other depending on implementation.
	Code any `json:"code" example:"200" format:"any" doc:"Application-level status code"`

	// Msg is a human-readable message describing the result.
	Msg string `json:"msg" example:"Success" doc:"Human-readable status message"`

	// Data contains the response payload when the request is successful.
	// This field is optional.
	Data any `json:"data" doc:"Response payload"`

	// Err contains error information when the request fails.
	// This field is optional.
	Err *ErrorModel `json:"err" doc:"Detailed error info when applicable"`

	// TraceID is used for tracing logs and request tracking across services.
	// It can be a string or other trace identifier format.
	TraceID string `json:"traceId,omitempty" example:"d1qlgrkqihv7a7754j20" doc:"Tracing identifier (xid) for debugging/logging"`
}

// Error satisfies the `error` interface. It returns the error's detail field.
func (e *GeneralResponseError) Error() string {
	return e.Err.Detail
}

// Add an error to the `Errors` slice. If passed a struct that satisfies the
// `huma.ErrorDetailer` interface, then it is used, otherwise the error
// string is used as the error detail message.
//
//	err := &xresp.GeneralResponseError{ /* ... */ }
//	err.Add(&huma.ErrorDetail{
//		Message: "expected boolean",
//		Location: "body.friends[1].active",
//		Value: 5
//	})
func (e *GeneralResponseError) Add(err error) {
	if converted, ok := err.(huma.ErrorDetailer); ok {
		e.Err.Errors = append(e.Err.Errors, converted.ErrorDetail())
		return
	}

	e.Err.Errors = append(e.Err.Errors, &huma.ErrorDetail{Message: err.Error()})
}

// GetStatus returns the HTTP status that should be returned to the client
// for this error.
func (e *GeneralResponseError) GetStatus() int {
	status, _ := e.Code.(int)
	return status
}

// ContentType provides a filter to adjust response content types. This is
// used to ensure e.g. `application/problem+json` content types defined in
// RFC 9457 Problem Details for HTTP APIs are used in responses to clients.
func (e *GeneralResponseError) ContentType(ct string) string {
	if ct == "application/json" {
		return "application/problem+json"
	}
	if ct == "application/cbor" {
		return "application/problem+cbor"
	}
	return ct
}
