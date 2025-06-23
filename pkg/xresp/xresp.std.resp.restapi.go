package xresp

import (
	"net/http"
	"sync"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// A Wrapper HTTP Rest API Response Builder Initiator
func NewRestResponse[D, E any](e echo.Context) RestResponseSTD[D, E] {
	traceId, _ := e.Get(xlog.XLOG_TRACE_ID_KEY).(xid.ID)
	return &restResponseSTD[D, E]{
		ResponseSTD: ResponseSTD[D, E]{
			echoCtx:        e,
			responseWriter: e.Response().Writer,
			request:        e.Request(),
			TraceID:        traceId,
		},
	}
}

// A Wrapper Standard Response For HTTP REST API
//
// HTTP Rest API Response Builder
type RestResponseSTD[D any, E any] interface {
	Msg(msg string) RestResponseSTD[D, E]
	Code(code any) RestResponseSTD[D, E]
	Data(data D) RestResponseSTD[D, E]
	Error(err E) RestResponseSTD[D, E]

	// Setter HTTP Response Status Code
	//
	// HTTP status codes as registered with IANA.
	// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
	StatusCode(code int) RestResponseSTD[D, E]

	// Add HTTP Response Header, this is equal function with:
	//	responseWriter.Header().Add(key, value)
	AddHeader(key string, value string) RestResponseSTD[D, E]

	// Delete HTTP Response Header by Key, this is equal function with:
	//	responseWriter.Header().Del(key)
	DelHeader(key string) RestResponseSTD[D, E]

	// A Method that Usefull for Run Interceptor Only,
	//
	// When you not need a response api that only need run interceptor.
	Done() error

	// A JSON Response Encoder for HTTP Response Writer, this is also auto set header and status code
	//	responseWriter.Header().Add("Accept", "application/json")
	//	responseWriter.Header().Add("Content-Type", "application/json")
	//
	//	// Default httpStatusCode is 200
	//	responseWriter.WriteHeader(httpStatusCode)
	JSON() error
	RawJSON() error

	// A General Purpose JSON Response Encoder to text
	JSONText() (string, error)
}

type restResponseSTD[D any, E any] struct {
	onceFn sync.Once
	ResponseSTD[D, E]
}

// # GETTER

func (r *restResponseSTD[D, E]) GetMsg() string {
	return r.ResponseSTD.Msg
}

func (r *restResponseSTD[D, E]) GetCode() any {
	return r.ResponseSTD.Code
}

func (r *restResponseSTD[D, E]) GetData() D {
	return r.ResponseSTD.Data
}

func (r *restResponseSTD[D, E]) GetError() E {
	return r.ResponseSTD.Err
}

func (r *restResponseSTD[D, E]) GetStatusCode() int {
	return r.ResponseSTD.statusCode
}

func (r *restResponseSTD[D, E]) GetResponseHeader() http.Header {
	return r.ResponseSTD.responseHeader
}

// # SETTER

func (r *restResponseSTD[D, E]) Msg(msg string) RestResponseSTD[D, E] {
	r.ResponseSTD.SetMsg(msg)
	return r
}

func (r *restResponseSTD[D, E]) Code(code any) RestResponseSTD[D, E] {
	r.ResponseSTD.SetCode(code)
	return r
}

func (r *restResponseSTD[D, E]) Data(data D) RestResponseSTD[D, E] {
	r.ResponseSTD.SetData(data)
	return r
}

func (r *restResponseSTD[D, E]) Error(err E) RestResponseSTD[D, E] {
	r.ResponseSTD.SetError(err)
	return r
}

// # HTTP SETTER

func (r *restResponseSTD[D, E]) StatusCode(code int) RestResponseSTD[D, E] {
	r.ResponseSTD.SetStatusCode(code)
	return r
}

func (r *restResponseSTD[D, E]) AddHeader(key string, value string) RestResponseSTD[D, E] {
	r.ResponseSTD.SetHeader(key, value)
	return r
}

func (r *restResponseSTD[D, E]) DelHeader(key string) RestResponseSTD[D, E] {
	r.ResponseSTD.DelHeader(key)
	return r
}

// # BUILDER

func (r *restResponseSTD[D, E]) Done() error {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK // OK as Default
	}

	return nil
}

func (r *restResponseSTD[D, E]) JSON() error {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK // OK as Default
	}

	r.onceFn.Do(r.ResponseSTD.RestJSON)

	return nil
}

func (r *restResponseSTD[D, E]) RawJSON() error {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK // OK as Default
	}

	r.onceFn.Do(r.ResponseSTD.RestRawJSON)

	return nil
}

func (r *restResponseSTD[D, E]) JSONText() (string, error) {
	return r.ResponseSTD.JSONText()
}
