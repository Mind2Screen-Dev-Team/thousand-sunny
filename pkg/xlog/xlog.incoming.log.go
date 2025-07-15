package xlog

import "time"

type (
	IncomingLogData struct {
		TimeStart time.Time `json:"timeStart"`
		TimeEnd   time.Time `json:"timeEnd"`

		ReqFormBody IncomingLogFormData `json:"reqFormBody"`
		ReqHeader   map[string][]string `json:"reqHeader"`
		ReqTraceID  string              `json:"reqTraceId"`
		ReqProtocol string              `json:"reqProtocol"`
		ReqIP       string              `json:"reqIP"`
		ReqIPs      []string            `json:"reqIPs"`
		ReqUA       []byte              `json:"reqUA"`
		ReqURI      []byte              `json:"reqURI"`
		ReqMethod   []byte              `json:"reqMethod"`
		ReqBody     []byte              `json:"reqBody"`
		ReqSize     int64               `json:"reqSize"`

		ResHeader map[string][]string `json:"resHeader"`
		ResStatus int                 `json:"resStatus"`
		ResSize   int64               `json:"resSize"`
		ResBody   []byte              `json:"resBody"`

		PanicMsg   string `json:"panicMsg"`
		PanicStack []byte `json:"panicStack"`

		IsPanic            bool `json:"isPanic"`
		IsHideRes          bool `json:"isHideRes"`
		IsMultipart        bool `josn:"isMultipart"`
		IsMultipartEncoded bool `josn:"isMultipartEncoded"`
		IsCompressed       bool `json:"isCompressed"`
	}

	IncomingLogFormData struct {
		Values map[string][]string `json:"values"`
		Files  map[string]FileInfo `json:"files"`
	}

	FileInfo struct {
		FileName    string `json:"filename"`
		ContentType string `json:"contentType"`
		Size        int64  `json:"size"`
	}
)
