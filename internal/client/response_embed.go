package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type APIStatusResponseErrorRootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func (e APIStatusResponseErrorRootCause) Error() string {
	return e.String()
}

func (e APIStatusResponseErrorRootCause) String() string {
	if e.Type == "" && e.Reason == "" {
		return ""
	}
	return fmt.Sprintf("type=%q; reason=%q", e.Type, e.Reason)
}

type APIStatusResponseError struct {
	RootCause []APIStatusResponseErrorRootCause `json:"root_cause"`
}

func (e APIStatusResponseError) Error() string {
	// if no errors, return empty
	if len(e.RootCause) == 0 {
		return ""
	}

	// create container for our ultimate error
	var finalErr error

	// append each error to a multierr
	for _, cause := range e.RootCause {
		finalErr = multierror.Append(finalErr, cause)
	}

	// should not be possible, but check just in case
	if finalErr == nil {
		return ""
	}

	// otherwise, return stringified error
	return finalErr.Error()
}

type APIStatusResponse struct {
	Status   string                  `json:"status"`
	Message  string                  `json:"message"`
	APIError *APIStatusResponseError `json:"error"`

	WarningsHeader []string `json:"-"`
}

func (e APIStatusResponse) populated() bool {
	return e.APIError != nil || e.Message != "" || e.Status != ""
}

func (e APIStatusResponse) HasErrors() bool {
	return e.APIError != nil && len(e.APIError.RootCause) > 0
}

func (e APIStatusResponse) Error() string {
	if e.APIError == nil {
		return ""
	}
	return e.APIError.Error()
}

func (e APIStatusResponse) String() string {
	// construct output container
	bits := make([]string, 0)

	// check each field for "non-empty"

	if e.Status != "" {
		bits = append(bits, fmt.Sprintf("status=%q", e.Status))
	}
	if e.Message != "" {
		bits = append(bits, fmt.Sprintf("message=%q", e.Message))
	}
	if len(e.WarningsHeader) > 0 {
		for i, w := range e.WarningsHeader {
			bits = append(bits, fmt.Sprintf("warning_%d=%q", i, w))
		}
	}
	if e.HasErrors() {
		for i, rc := range e.APIError.RootCause {
			bits = append(
				bits,
				fmt.Sprintf("error_%d_type=%q", i, rc.Type),
				fmt.Sprintf("error_%d_reason=%q", i, rc.Reason),
			)
		}
	}

	// if at least 1 field contained a value
	if len(bits) > 0 {
		return strings.Join(bits, "; ")
	}
	return ""
}

func (e APIStatusResponse) AppendDiagnostics(d diag.Diagnostics) {
	// add warnings from header
	for _, w := range e.WarningsHeader {
		d.AddWarning(
			w,
			w,
		)
	}
	// add any / all errors
	if e.HasErrors() {
		for _, source := range e.APIError.RootCause {
			d.AddError(
				source.Type,
				source.Reason,
			)
		}
	}
}

func TryUnmarshalEmbed(b []byte) (APIStatusResponse, map[string]json.RawMessage, error) {
	embed := APIStatusResponse{}

	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return embed, m, err
	}

	if errs, ok := m["error"]; ok {
		if err := json.Unmarshal(errs, embed.APIError); err != nil {
			return embed, m, fmt.Errorf("error unmarshalling error: %w", err)
		}
	}
	if status, ok := m["status"]; ok {
		if err := json.Unmarshal(status, &embed.Status); err != nil {
			return embed, m, fmt.Errorf("error unmarshalling status: %w", err)
		}
	}
	if msg, ok := m["message"]; ok {
		if err := json.Unmarshal(msg, &embed.Message); err != nil {
			return embed, m, fmt.Errorf("error unmarshalling message: %w", err)
		}
	}

	return embed, m, nil
}
