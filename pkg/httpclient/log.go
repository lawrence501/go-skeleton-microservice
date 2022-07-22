package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/moul/http2curl"
	"go.uber.org/zap"
)

var (
	jwtRegexp       = regexp.MustCompile(`(?i)(Authorization: bearer |X-Authorization-Header: *)([a-z0-9-_]*)\.([a-z0-9-_]*)\.([a-z0-9-_]*)`)
	basicAuthRegexp = regexp.MustCompile(`(?i)(Authorization: Basic )([A-Za-z0-9+\/=]*)`)
)

type LoggingRoundTrip struct {
	Base http.RoundTripper
}

func (t LoggingRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	logger := zap.L().With(zap.String("component", "httpclient"), zap.String("url", req.URL.String()))
	logger.Debug("req-begin", zap.String("curl", getRequestString(req)))

	resp, err := t.Base.RoundTrip(req)
	if resp != nil {
		logger = logger.With(zap.Int("status", resp.StatusCode))
	}
	logger.Info("req-summary", zap.Int64("duration_ms", time.Since(start).Milliseconds()), zap.Error(err))
	if err != nil {
		logger.Error("Failed to reach endpoint", zap.Error(err))
	}
	logger.Debug("req-end", zap.Any("headers", req.Header), zap.String("responseBody", getResponseBody(resp)))
	return resp, err
}

func getRequestString(req *http.Request) string {
	cmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return fmt.Sprintf("Error creating curl command for %s", req.URL.String())
	}
	return maskCurl(cmd.String())
}

func maskCurl(curl string) string {
	curl = jwtRegexp.ReplaceAllString(curl, "$1$2.$3.[REDACTED]")
	curl = basicAuthRegexp.ReplaceAllString(curl, "${1}[REDACTED]")
	return curl
}

func getResponseBody(res *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "Error reading response body"
	}
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
