package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type APIResponseMetaErrorRootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func (e APIResponseMetaErrorRootCause) Error() string {
	if e.Type == "" && e.Reason == "" {
		return ""
	}
	return fmt.Sprintf("type=%q; reason=%q", e.Type, e.Reason)
}

type APIResponseMetaError struct {
	RootCause []APIResponseMetaErrorRootCause `json:"root_cause"`
}

func (e APIResponseMetaError) Error() string {
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

type APIResponseMeta struct {
	Status   string                `json:"status"`
	Message  string                `json:"message"`
	TheError *APIResponseMetaError `json:"error"`
}

func (e APIResponseMeta) HasErrors() bool {
	return e.TheError != nil && len(e.TheError.RootCause) > 0
}

func (e APIResponseMeta) Error() string {
	if e.TheError == nil {
		return ""
	}
	return e.TheError.Error()
}

func (e APIResponseMeta) populated() bool {
	return e.TheError != nil || e.Message != "" || e.Status != ""
}

func (e APIResponseMeta) String() string {
	bits := make([]string, 0)
	if e.Status != "" {
		bits = append(bits, fmt.Sprintf("status=%q", e.Status))
	}
	if e.Message != "" {
		bits = append(bits, fmt.Sprintf("message=%q", e.Message))
	}
	if e.HasErrors() {
		bits = append(bits, fmt.Sprintf("errors=%q", e.Error()))
	}
	if len(bits) > 0 {
		return strings.Join(bits, "; ")
	}
	return ""
}

func (e APIResponseMeta) AppendErrorsToDiagnostic(d diag.Diagnostics) {
	if !e.HasErrors() {
		return
	}
	for _, source := range e.TheError.RootCause {
		d.AddError(
			source.Reason,
			source.Error(),
		)
	}
}

func TryUnmarshalEmbed(b []byte) (APIResponseMeta, map[string]json.RawMessage, error) {
	embed := APIResponseMeta{}

	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return embed, m, err
	}

	if errs, ok := m["error"]; ok {
		if err := json.Unmarshal(errs, embed.TheError); err != nil {
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
