package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/v3/validation"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
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

	InsecureSkipTLSVerify types.Bool `tfsdk:"insecure_skip_tls_verify"`

	UseResponseCheckOnly types.Bool `tfsdk:"use_response_check_only"`

	SkipInitProductCheck types.Bool `tfsdk:"skip_init_product_check"`
}

var _ provider.Provider = &OpenSearchProvider{}

type OpenSearchProvider struct {
	version string
}

func (p *OpenSearchProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opensearch"
	resp.Version = p.version
}

func (p *OpenSearchProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "OpenSearch Provider",
		Attributes: map[string]schema.Attribute{
			fields.ConfigAttrAddresses: schema.ListAttribute{
				Description: "List of addresses to connect to",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					validation.IsURL(),
				},
			},
			fields.ConfigAttrUsername: schema.StringAttribute{
				Description: "Username for HTTP basic authentication",
				Optional:    true,
			},
			fields.ConfigAttrPassword: schema.StringAttribute{
				Description: "Password for HTTP basic authentication",
				Sensitive:   true,
				Optional:    true,
			},
			fields.ConfigAttrCACert: schema.StringAttribute{
				Description: "PEM Encoded certificate authorities",
				Sensitive:   true,
				Optional:    true,
			},
			fields.ConfigAttrRetryOnStatuses: schema.ListAttribute{
				Description: "List of status codes for retry",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			fields.ConfigAttrDisableRetry: schema.BoolAttribute{
				Description: "Disable all request retries",
				Optional:    true,
			},
			fields.ConfigAttrEnableRetryOnTimeout: schema.BoolAttribute{
				Description: "Enables request retry on timeout",
				Optional:    true,
			},
			fields.ConfigAttrMaxRetries: schema.Int64Attribute{
				Description: "Maximum number of times a given request can be retried",
				Optional:    true,
			},
			fields.ConfigAttrCompressRequestBody: schema.BoolAttribute{
				Description: "Enable request body compression",
				Optional:    true,
			},
			fields.ConfigAttrInsecureSkipTLSVerify: schema.BoolAttribute{
				Description: "Disable TLS verification",
				Optional:    true,
			},
			fields.ConfigAttrUseResponseCheckOnly: schema.BoolAttribute{
				Description: "Disable executing product check on every request",
				Optional:    true,
			},
			fields.ConfigAttrSkipInitProductCheck: schema.BoolAttribute{
				Description: "Skip product check API call on configure",
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

		shared Shared

		// create pooled transport
		transport = cleanhttp.DefaultPooledTransport()
	)

	// attempt to parse provider config
	resp.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// configure transport
	if conf.InsecureSkipTLSVerify.IsNull() == false && conf.InsecureSkipTLSVerify.IsUnknown() == false && conv.BoolValueToBool(conf.InsecureSkipTLSVerify) == true {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// build base opensearch client config
	clientConfig = opensearch.Config{
		Addresses:            conv.StringListToStrings(conf.Addresses),
		Transport:            transport,
		Username:             conf.Username.ValueString(),
		Password:             conf.Password.ValueString(),
		RetryOnStatus:        conv.Int64ListToInts(conf.RetryOnStatus),
		DisableRetry:         conf.DisableRetry.ValueBool(),
		EnableRetryOnTimeout: conf.EnableRetryOnTimeout.ValueBool(),
		MaxRetries:           conv.Int64ValueToInt(conf.MaxRetries),
		CompressRequestBody:  conf.CompressRequestBody.ValueBool(),
		UseResponseCheckOnly: conf.UseResponseCheckOnly.ValueBool(),
	}

	// did they provide ca's?
	if !conf.CACert.IsNull() && !conf.CACert.IsUnknown() {
		clientConfig.CACert = []byte(conf.CACert.ValueString())
	}

	// if there were any errors during config, bail out now
	if resp.Diagnostics.HasError() {
		return
	}

	// attempt to perform connectivity and fitment test
	if !conf.SkipInitProductCheck.ValueBool() {
		infoReq := opensearchapi.InfoRequest{}
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if _, err := infoReq.Do(ctx, client); err != nil {
			resp.Diagnostics.AddError(
				"Error performing init compatibility check",
				fmt.Sprintf("Error occurred during init compatibility check: %v", err),
			)
			return
		}
	}

	// create shared object for use in resource and datasource types
	shared = Shared{
		Client: client,
	}

	// set shared
	resp.ResourceData = &shared
	resp.DataSourceData = &shared
}

func (p *OpenSearchProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPluginSecurityRoleResource,
		//NewPluginSecurityUserResource,
	}
}

func (p *OpenSearchProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewExampleDataSource,
	}
}

func NewOpenSearchProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenSearchProvider{
			version: version,
		}
	}
}
