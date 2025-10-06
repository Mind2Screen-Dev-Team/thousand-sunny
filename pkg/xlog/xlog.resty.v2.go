package xlog

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type RestyV2Logger struct {
	Log Logger
}

func NewRestyV2Logger(log Logger) *RestyV2Logger {
	return &RestyV2Logger{log}
}

func (l *RestyV2Logger) Errorf(format string, v ...any) {
	var (
		ctx     = context.Background()
		fields  = make([]any, 0)
		logText = fmt.Sprintf(format, v...)
	)

	fields = append(fields,
		"restyRawLog", logText,
	)

	l.Log.Error(ctx, "resty api log", fields...)
}

func (l *RestyV2Logger) Warnf(format string, v ...any) {
	var (
		ctx       = context.Background()
		fields    = make([]any, 0)
		logText   = fmt.Sprintf(format, v...)
		parsed, _ = ParseRestyLog(logText)
	)

	fields = append(fields,
		"restyRawLog", logText,
		"restyParsedLog", parsed,
	)

	l.Log.Warn(ctx, "resty api log", fields...)
}

func (l *RestyV2Logger) Debugf(format string, v ...any) {
	var (
		ctx       = context.Background()
		fields    = make([]any, 0)
		logText   = fmt.Sprintf(format, v...)
		parsed, _ = ParseRestyLog(logText)
	)

	fields = append(fields,
		"restyRawLog", logText,
		"restyParsedLog", parsed,
	)

	l.Log.Debug(ctx, "resty api log", fields...)
}

type (
	HTTPLog struct {
		Curl     string       `json:"curl"`
		Request  HTTPRequest  `json:"request"`
		Response HTTPResponse `json:"response"`
	}

	HTTPRequest struct {
		Method  string              `json:"method"`
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	HTTPResponse struct {
		Status     string              `json:"status"`
		Proto      string              `json:"proto"`
		ReceivedAt string              `json:"receivedAt"`
		Duration   string              `json:"duration"`
		Headers    map[string][]string `json:"headers"`
		Body       string              `json:"body"`
	}
)

func ParseRestyLog(log string) (*HTTPLog, error) {
	re := regexp.MustCompile(`(?ms)(?:~~~ REQUEST\(CURL\) ~~~\s+(?P<request_curl>.*?)\s+)?~~~ REQUEST ~~~\s+(?P<request>.*?)^-+.*?~~~ RESPONSE ~~~\s+(?P<response>.*?)^=+`)
	match := re.FindStringSubmatch(log)
	if match == nil {
		return nil, fmt.Errorf("log format not matched")
	}

	// Extract groups into a map[string]string
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	var (
		curl, _  = result["request_curl"]
		request  = parseRequest(result["request"])
		response = parseResponse(result["response"])
	)

	return &HTTPLog{Curl: strings.TrimSpace(curl), Request: request, Response: response}, nil
}

func parseRequest(req string) HTTPRequest {
	var (
		lines                = strings.Split(req, "\n")
		headers              = map[string][]string{}
		readBody, readHeader = false, false
		method, url          = "", ""
		body                 = ""
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "GET") || strings.HasPrefix(line, "POST") || strings.HasPrefix(line, "PUT") || strings.HasPrefix(line, "DELETE") || strings.HasPrefix(line, "PATCH"):
			if parts := strings.Fields(line); len(parts) >= 2 {
				method = parts[0]
				url = parts[1]
			}
		case line == "HEADERS:":
			readHeader = true
		case strings.HasPrefix(line, "BODY"):
			readHeader = false
			readBody = true
		case readHeader:
			k, v := parseHeader(line)
			headers[k] = append(headers[k], v)
		case readBody:
			body += line + "\n"
		}
	}

	body = strings.TrimSpace(body)
	if body == "***** NO CONTENT *****" {
		body = ""
	}

	return HTTPRequest{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}
}

func parseResponse(resp string) HTTPResponse {
	var (
		lines                               = strings.Split(resp, "\n")
		headers                             = map[string][]string{}
		readBody, readHeader                = false, false
		status, proto, receivedAt, duration = "", "", "", ""
		body                                = ""
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "STATUS"):
			status = extractValue(line)
		case strings.HasPrefix(line, "PROTO"):
			proto = extractValue(line)
		case strings.HasPrefix(line, "RECEIVED AT"):
			receivedAt = extractValue(line)
		case strings.HasPrefix(line, "TIME DURATION"):
			duration = extractValue(line)
		case strings.HasPrefix(line, "HEADERS"):
			readHeader = true
		case strings.HasPrefix(line, "BODY"):
			readHeader = false
			readBody = true
		case readHeader:
			k, v := parseHeader(line)
			headers[k] = append(headers[k], v)
		case readBody:
			body += line + "\n"
		}
	}

	return HTTPResponse{
		Status:     status,
		Proto:      proto,
		ReceivedAt: receivedAt,
		Duration:   duration,
		Headers:    headers,
		Body:       strings.TrimSpace(body),
	}
}

func extractValue(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func parseHeader(line string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
	if len(parts) == 2 {
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		return key, value
	}
	return "", ""
}
