package provider

import "github.com/opensearch-project/opensearch-go"

type Shared struct {
	Client *opensearch.Client
}
