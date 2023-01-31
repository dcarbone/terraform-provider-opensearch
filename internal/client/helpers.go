package client

import (
	"context"
	"io"
	"net/http"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func newOpenSearchRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return http.NewRequestWithContext(ctx, method, path, body)
}

func addOpenSearchRequestHeaders(req *http.Request, headers http.Header) {
	l := len(headers)
	if l == 0 {
		return
	}
	if len(req.Header) == 0 {
		req.Header = make(http.Header, l)
	}
	for k, v := range headers {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
}

func buildOpenSearchAPIResponse(res *http.Response) *opensearchapi.Response {
	response := opensearchapi.Response{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Body:       res.Body,
	}
	return &response
}
