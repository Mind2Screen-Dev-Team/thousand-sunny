package xhttp

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type FileInfo struct {
	FileName    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func CopyRequestBody(r *http.Request) []byte {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		// Refill the original request body to preserve it for further usage.
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return bodyBytes
}

func DeepCopyRequest(r *http.Request) *http.Request {
	var (
		// Read the body if it's non-nil
		bodyBytes = CopyRequestBody(r)

		// Create a shallow copy of the request
		rCopy = r.Clone(r.Context())
	)

	// Replace the body of the new request with a new reader wrapping the copied bytes
	if bodyBytes != nil {
		rCopy.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Copy MultipartForm if it has been parsed
	if r.MultipartForm != nil {
		rCopy.MultipartForm = &multipart.Form{
			Value: make(map[string][]string),
			File:  make(map[string][]*multipart.FileHeader),
		}

		// Deep copy form values
		for key, values := range r.MultipartForm.Value {
			rCopy.MultipartForm.Value[key] = append([]string{}, values...)
		}

		// Deep copy file headers
		for key, fileHeaders := range r.MultipartForm.File {
			// Create a new slice to hold the copied file headers
			var copiedFileHeaders []*multipart.FileHeader
			for _, fileHeader := range fileHeaders {
				// Create a new instance of multipart.FileHeader
				newFileHeader := &multipart.FileHeader{
					Filename: fileHeader.Filename,
					Header:   make(textproto.MIMEHeader),
					Size:     fileHeader.Size,
				}

				// Copy all headers (may need to consider the Content-Type or other relevant headers)
				for k, v := range fileHeader.Header {
					newFileHeader.Header[k] = append([]string{}, v...)
				}

				// Append the new file header to the copied slice
				copiedFileHeaders = append(copiedFileHeaders, newFileHeader)
			}

			// Assign the copied file headers to the new request
			rCopy.MultipartForm.File[key] = copiedFileHeaders
		}
	}

	return rCopy
}
