package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TerraformLogger struct {
	provCTX context.Context

	reqBodyEnabled  bool
	respBodyEnabled bool
}

func NewTerraformLogger(provCTX context.Context, requestBodyEnabled, responseBodyEnabled bool) TerraformLogger {
	tl := TerraformLogger{
		provCTX:         provCTX,
		reqBodyEnabled:  requestBodyEnabled,
		respBodyEnabled: responseBodyEnabled,
	}
	return tl
}

func (tl TerraformLogger) LogRoundTrip(req *http.Request, resp *http.Response, err error, start time.Time, dur time.Duration) error {
	var (
		logFields = map[string]interface{}{
			"start":    start.Format(time.RFC3339),
			"duration": dur,
		}
	)

	defer func() {
		if req != nil && req.Body != nil {
			_ = req.Body.Close()
		}
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if req != nil {
		logFields["url"] = req.URL.String()
		if tl.reqBodyEnabled && req.Body != nil {
			reqBodyBytes, _ := io.ReadAll(req.Body)
			logFields["request_body"] = string(reqBodyBytes)
			logFields["request_body_len"] = len(reqBodyBytes)
		}
	}

	if tl.respBodyEnabled && resp != nil && resp.Body != nil {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		logFields["response_body"] = string(respBodyBytes)
		logFields["response_body_len"] = len(respBodyBytes)
	}

	if err != nil {
		logFields["err"] = err.Error()
		tflog.Error(tl.provCTX, fmt.Sprintf("OpenSearch client error: %v", err), logFields)
	} else {
		tflog.Trace(tl.provCTX, "OpenSearch client query tracer", logFields)
	}

	return nil
}

func (tl TerraformLogger) RequestBodyEnabled() bool {
	return tl.reqBodyEnabled
}

func (tl TerraformLogger) ResponseBodyEnabled() bool {
	return tl.respBodyEnabled
}
