package client

import (
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

func (e apiResponseEmbed) AddDiagnosticErrors(d diag.Diagnostics) {
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
