package provider

import (
	"context"
	"crypto/tls"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/v3/validation"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opensearch-project/opensearch-go"
)

const (
	configAttrAddresses            = "addresses"
	configAttrUsername             = "username"
	configAttrPassword             = "password"
	configAttrCACert               = "ca_cert"
	configAttrRetryOnStatuses      = "retry_on_statuses"
	configAttrDisableRetry         = "disable_retry"
	configAttrEnableRetryOnTimeout = "enable_retry_on_timeout"
	configAttrMaxRetries           = "max_retries"
	configAttrCompressRequestBody  = "compress_request_body"

	configAttrInsecureSkipTLSVerify = "insecure_skip_tls_verify"
)

type OpenSearchProviderConfig struct {
	Addresses types.List `tfsdk:"addresses"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`

	CACert types.String `tfsdk:"ca_cert"`

	RetryOnStatus        types.List  `tfsdk:"retry_on_status"`
	DisableRetry         types.Bool  `tfsdk:"disable_retry"`
	EnableRetryOnTimeout types.Bool  `tfsdk:"enable_retry_on_timeout"`
	MaxRetries           types.Int64 `tfsdk:"max_retries"`

	CompressRequestBody types.Bool `tfsdk:"compress_request_body"`

	InsecureSkipTLSVerify types.Bool `json:"insecure_skip_tls_verify"`
}

var _ provider.Provider = &OpenSearchProvider{}

type OpenSearchProvider struct {
	version string

	configured bool
}

func (p *OpenSearchProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opensearch"
	resp.Version = p.version
}

func (p *OpenSearchProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			configAttrAddresses: schema.ListAttribute{
				Description: "List of master node addresses in your cluster",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					validation.Required(),
					validation.IsURL(),
				},
			},
			configAttrUsername: schema.StringAttribute{
				Description: "Username for HTTP basic authentication",
				Optional:    true,
			},
			configAttrPassword: schema.StringAttribute{
				Description: "Password for HTTP basic authentication",
				Sensitive:   true,
				Optional:    true,
			},
			configAttrCACert: schema.StringAttribute{
				Description: "PEM Encoded certificate authorities",
				Sensitive:   true,
				Optional:    true,
			},
			configAttrRetryOnStatuses: schema.ListAttribute{
				Description: "List of status codes for retry",
				ElementType: types.Int64Type,
				Optional:    true,
			},
			configAttrDisableRetry: schema.BoolAttribute{
				Description: "Disable all request retries",
				Optional:    true,
			},
			configAttrEnableRetryOnTimeout: schema.BoolAttribute{
				Description: "Enables request retry on timeout",
				Optional:    true,
			},
			configAttrMaxRetries: schema.Int64Attribute{
				Description: "Maximum number of times a given request can be retried",
				Optional:    true,
			},
			configAttrCompressRequestBody: schema.BoolAttribute{
				Description: "Enable request body compression",
				Optional:    true,
			},
			configAttrInsecureSkipTLSVerify: schema.BoolAttribute{
				Description: "Disable TLS verification",
				Optional:    true,
			},
		},
	}
}

func (p *OpenSearchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var (
		conf OpenSearchProviderConfig

		clientConfig opensearch.Config
		client       *opensearch.Client

		// create pooled transport
		transport = cleanhttp.DefaultPooledTransport()
	)

	// attempt to parse provider config
	resp.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// build base opensearch client config
	clientConfig = opensearch.Config{
		Addresses:            conv.StringListToStrings(conf.Addresses),
		Transport:            transport,
		Username:             conv.AttributeValueToString(conf.Username),
		Password:             conv.AttributeValueToString(conf.Password),
		RetryOnStatus:        conv.Int64ListToInts(conf.RetryOnStatus),
		DisableRetry:         conv.BoolValueToBool(conf.DisableRetry),
		EnableRetryOnTimeout: conv.BoolValueToBool(conf.EnableRetryOnTimeout),
		MaxRetries:           conv.Int64ValueToInt(conf.MaxRetries),
		CompressRequestBody:  conv.BoolValueToBool(conf.CompressRequestBody),
	}

	// did they provide ca's?
	if conf.CACert.IsNull() == false && conf.CACert.IsUnknown() == false {
		clientConfig.CACert = []byte(conv.AttributeValueToString(conf.CACert))
	}
	// should
	if conf.InsecureSkipTLSVerify.IsNull() == false && conf.InsecureSkipTLSVerify.IsUnknown() == false && conv.BoolValueToBool(conf.InsecureSkipTLSVerify) == true {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		clientConfig.Transport = transport
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: nope.
	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *OpenSearchProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *OpenSearchProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenSearchProvider{
			version: version,
		}
	}
}
