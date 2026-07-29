package engine

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"codedock.run/codedock-tunnel/pkg/protocol"
)

type HTTPProxy struct {
	baseDomain   string
	router       *RequestRouter
	maxBodyBytes int64
	recorder     UsageRecorder
}

type UsageMeasurement struct {
	OrganizationID string
	TunnelID       string
	EventType      string
	Bytes          int64
	Connections    int
}

type UsageRecorder interface {
	Record(context.Context, UsageMeasurement) error
}

func NewHTTPProxy(baseDomain string, router *RequestRouter, maxBodyBytes int64) (*HTTPProxy, error) {
	if baseDomain == "" {
		return nil, fmt.Errorf("tunnel base domain is required")
	}
	if router == nil {
		return nil, fmt.Errorf("request router is required")
	}
	return &HTTPProxy{baseDomain: strings.ToLower(strings.TrimSuffix(baseDomain, ".")), router: router, maxBodyBytes: maxBodyBytes}, nil
}

func (p *HTTPProxy) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	route, ok := resolveRouteKey(request.Host, p.baseDomain)
	if !ok {
		http.Error(response, "tunnel not found", http.StatusNotFound)
		return
	}
	body, err := readBody(request.Body, p.maxBodyBytes)
	if err != nil {
		if err == errBodyTooLarge {
			http.Error(response, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(response, "unable to read request body", http.StatusBadRequest)
		return
	}
	forwarded, err := p.router.ForwardHTTP(request.Context(), route, protocol.HTTPRequest{
		Method:  request.Method,
		Path:    request.URL.RequestURI(),
		Headers: requestHeaders(request.Header),
		Body:    body,
	})
	if decodedRequest, decodeErr := base64.StdEncoding.DecodeString(body); decodeErr == nil {
		p.record(request, route, "http_request", int64(len(decodedRequest)))
	}
	if err != nil {
		p.record(request, route, "http_error", 0)
		http.Error(response, "tunnel unavailable", http.StatusBadGateway)
		return
	}
	if forwarded.Error != "" {
		p.record(request, route, "http_error", 0)
		http.Error(response, forwarded.Error, http.StatusBadGateway)
		return
	}
	writeResponseHeaders(response.Header(), forwarded.Headers)
	status := forwarded.StatusCode
	if status < http.StatusContinue || status > 599 {
		status = http.StatusBadGateway
	}
	response.WriteHeader(status)
	if forwarded.Body == "" {
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(forwarded.Body)
	if err != nil {
		p.record(request, route, "http_error", 0)
		return
	}
	p.record(request, route, "http_response", int64(len(decoded)))
	_, _ = response.Write(decoded)
}

func (p *HTTPProxy) SetUsageRecorder(recorder UsageRecorder) {
	p.recorder = recorder
}

func (p *HTTPProxy) record(request *http.Request, tunnelID, eventType string, bytes int64) {
	if p.recorder == nil {
		return
	}
	organizationID, ok := p.router.OrganizationID(tunnelID)
	if !ok {
		return
	}
	_ = p.recorder.Record(request.Context(), UsageMeasurement{OrganizationID: organizationID, TunnelID: tunnelID, EventType: eventType, Bytes: bytes})
}

var errBodyTooLarge = fmt.Errorf("request body too large")

func readBody(body io.ReadCloser, maxBytes int64) (string, error) {
	if body == nil {
		return "", nil
	}
	defer body.Close()
	reader := io.Reader(body)
	if maxBytes > 0 {
		reader = io.LimitReader(body, maxBytes+1)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read request body: %w", err)
	}
	if maxBytes > 0 && int64(len(data)) > maxBytes {
		return "", errBodyTooLarge
	}
	if len(data) == 0 {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func resolveTunnelID(host string, baseDomain string) (string, bool) {
	cleanHost := strings.ToLower(strings.TrimSuffix(host, "."))
	if parsedHost, _, err := net.SplitHostPort(cleanHost); err == nil {
		cleanHost = parsedHost
	}
	baseDomain = strings.ToLower(strings.TrimSuffix(baseDomain, "."))
	prefix, found := strings.CutSuffix(cleanHost, "."+baseDomain)
	if !found || prefix == "" || strings.Contains(prefix, ".") || prefix == "www" {
		return "", false
	}
	return prefix, true
}

func resolveRouteKey(host string, baseDomain string) (string, bool) {
	cleanHost := strings.ToLower(strings.TrimSuffix(host, "."))
	if parsedHost, _, err := net.SplitHostPort(cleanHost); err == nil {
		cleanHost = parsedHost
	}
	baseDomain = strings.ToLower(strings.TrimSuffix(baseDomain, "."))
	prefix, found := strings.CutSuffix(cleanHost, "."+baseDomain)
	if found {
		if prefix == "" || strings.Contains(prefix, ".") || prefix == "www" {
			return "", false
		}
		return prefix, true
	}
	if cleanHost == "" || strings.Contains(cleanHost, " ") || strings.Contains(cleanHost, "..") || cleanHost == "www" {
		return "", false
	}
	return cleanHost, true
}

func requestHeaders(headers http.Header) map[string][]string {
	result := make(map[string][]string)
	for key, values := range headers {
		if isHopByHopHeader(key) {
			continue
		}
		result[key] = append([]string(nil), values...)
	}
	return result
}

func writeResponseHeaders(destination http.Header, source map[string][]string) {
	for key, values := range source {
		if isHopByHopHeader(key) {
			continue
		}
		for _, value := range values {
			destination.Add(key, value)
		}
	}
}

func isHopByHopHeader(key string) bool {
	switch strings.ToLower(key) {
	case "connection", "keep-alive", "proxy-authenticate", "proxy-authorization", "te", "trailer", "transfer-encoding", "upgrade":
		return true
	default:
		return false
	}
}
