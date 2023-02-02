package client

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type apiResponseEmbedErrorRootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func (e apiResponseEmbedErrorRootCause) Error() string {
	if e.Type == "" && e.Reason == "" {
		return ""
	}
	return fmt.Sprintf("type=%q; reason=%q", e.Type, e.Reason)
}

type apiResponseEmbedError struct {
	RootCause []apiResponseEmbedErrorRootCause `json:"root_cause"`
}

func (e apiResponseEmbedError) Error() string {
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

type apiResponseEmbed struct {
	Status   string                 `json:"status"`
	Message  string                 `json:"message"`
	TheError *apiResponseEmbedError `json:"error"`
}

func (e apiResponseEmbed) HasErrors() bool {
	return e.TheError != nil && len(e.TheError.RootCause) > 0
}

func (e apiResponseEmbed) Error() string {
	if e.TheError == nil {
		return ""
	}
	return e.TheError.Error()
}

func (e apiResponseEmbed) populated() bool {
	return e.TheError != nil || e.Message != "" || e.Status != ""
}

func (e apiResponseEmbed) AppendErrorsToDiagnostic(d diag.Diagnostics) {
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

func tryUnmarshalEmbed(b []byte) (apiResponseEmbed, map[string]json.RawMessage, error) {
	embed := apiResponseEmbed{}

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
