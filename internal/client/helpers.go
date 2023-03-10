package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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

func codeMatch(actual int, targets []int) bool {
	for _, t := range targets {
		if actual == t {
			return true
		}
	}
	return false
}

func codesToString(in []int) []string {
	out := make([]string, len(in))
	for i, n := range in {
		out[i] = strconv.Itoa(n)
	}
	return out
}

func ParseResponse(osResp *opensearchapi.Response, sink interface{}, okCodes ...int) error {
	// immediately queue up body closure handling
	defer HandleResponseCleanup(osResp)

	// test if response code matches a target code
	if codeMatch(osResp.StatusCode, okCodes) {
		// if no sink provided, end here
		if sink == nil {
			return nil
		}
		// otherwise, attempt to unmarshal response into sink
		if err := json.NewDecoder(osResp.Body).Decode(sink); err != nil {
			return err
		}
		// if the provided sink is of type *APIStatusResponse, add warnings from header to it
		if m, ok := sink.(*APIStatusResponse); ok && osResp.HasWarnings() {
			w := osResp.Warnings()
			m.WarningsHeader = make([]string, len(w))
			copy(m.WarningsHeader, w)
		}
		return nil
	}

	// otherwise, attempt to unmarshal response into meta
	meta := new(APIStatusResponse)

	// attempt to decode response
	if err := json.NewDecoder(osResp.Body).Decode(meta); err != nil {
		return err
	}

	// add any warnings from the header
	if osResp.HasWarnings() {
		w := osResp.Warnings()
		meta.WarningsHeader = make([]string, len(w))
		copy(meta.WarningsHeader, w)
	}

	// ensure we append the status code mismatch
	if meta.APIError == nil {
		meta.APIError = new(APIStatusResponseError)
	}

	meta.APIError.RootCause = append(meta.APIError.RootCause, APIStatusResponseErrorRootCause{
		Type:   "Status Code Mismatch",
		Reason: fmt.Sprintf("Actual response code %d does not match expected [%s]", osResp.StatusCode, strings.Join(codesToString(okCodes), " ")),
	})

	// return whatever we got here.
	return meta
}
