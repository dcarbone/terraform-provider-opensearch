package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/v3/validation"
	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type OpenSearchProviderConfigClientDebugLogger struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type OpenSearchProviderConfigRequestTraceLogger struct {
	Enabled             types.Bool `tfsdk:"enabled"`
	IncludeRequestBody  types.Bool `tfsdk:"include_request_body"`
	IncludeResponseBody types.Bool `tfsdk:"include_response_body"`
}

type OpenSearchProviderConfig struct {
	Addresses types.List `tfsdk:"addresses"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`

	CACert types.String `tfsdk:"ca_cert"`

	RetryOnStatus        types.List  `tfsdk:"retry_on_status"`
	DisableRetry         types.Bool  `tfsdk:"disable_retry"`
	EnableRetryOnTimeout types.Bool  `tfsdk:"enable_retry_on_timeout"`
	MaxRetries           types.Int64 `tfsdk:"max_retries"`

	CompressRequestBody   types.Bool `tfsdk:"compress_request_body"`
	InsecureSkipTLSVerify types.Bool `tfsdk:"insecure_skip_tls_verify"`
	EnableOnRequestCheck  types.Bool `tfsdk:"enable_on_request_check"`
	SkipInitProductCheck  types.Bool `tfsdk:"skip_init_product_check"`

	ClientDebugLogger  types.Object `tfsdk:"client_debug_logger"`
	RequestTraceLogger types.Object `tfsdk:"request_trace_logger"`
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
			fields.ConfigAttrRetryOnStatus: schema.ListAttribute{
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
			fields.ConfigAttrEnableOnRequestCheck: schema.BoolAttribute{
				Description: "By default, the opensearch-go client executes a \"compatibility check\" on every" +
					" single request made.  This has been disabled by default in this provider.  If you wish to" +
					" re-enable this, for whatever reason, set this to true.",
				Optional: true,
			},
			fields.ConfigAttrSkipInitProductCheck: schema.BoolAttribute{
				Description: "Skip product check API call on configure",
				Optional:    true,
			},
			fields.ConfigAttrClientDebugLogger: schema.ObjectAttribute{
				Description: "OpenSearch client debug logging configuration.  This writes debug-level logging" +
					" directly to stdout.  Do not enable outside of a local development environment.",
				Optional: true,
				AttributeTypes: map[string]attr.Type{
					fields.ConfigAttrEnabled: types.BoolType,
				},
			},
			fields.ConfigAttrRequestTraceLogger: schema.ObjectAttribute{
				Description: "OpenSearch client request tracing logger configuration.  This writes TRACE level" +
					" Terraform logs of every HTTP action by the opensearch-go client.  Can produce very chatty" +
					" logs.",
				Optional: true,
				AttributeTypes: map[string]attr.Type{
					fields.ConfigAttrEnabled:             types.BoolType,
					fields.ConfigAttrIncludeRequestBody:  types.BoolType,
					fields.ConfigAttrIncludeResponseBody: types.BoolType,
				},
			},
		},
	}
}

func (p *OpenSearchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var (
		conf         OpenSearchProviderConfig
		traceLogConf OpenSearchProviderConfigRequestTraceLogger
		dbgLogConf   OpenSearchProviderConfigClientDebugLogger
		osConfig     opensearch.Config
		osClient     *opensearch.Client
		shared       Shared
		err          error

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
	osConfig = opensearch.Config{
		Addresses:            conv.StringListToStrings(conf.Addresses),
		Transport:            transport,
		Username:             conf.Username.ValueString(),
		Password:             conf.Password.ValueString(),
		RetryOnStatus:        conv.Int64ListToInts(conf.RetryOnStatus),
		DisableRetry:         conf.DisableRetry.ValueBool(),
		EnableRetryOnTimeout: conf.EnableRetryOnTimeout.ValueBool(),
		MaxRetries:           conv.Int64ValueToInt(conf.MaxRetries),
		CompressRequestBody:  conf.CompressRequestBody.ValueBool(),
		UseResponseCheckOnly: !conf.EnableOnRequestCheck.ValueBool(),
	}

	// did they provide ca's?
	if !conf.CACert.IsNull() && !conf.CACert.IsUnknown() {
		osConfig.CACert = []byte(conf.CACert.ValueString())
	}

	// attempt to unmarshal request trace logging config
	resp.Diagnostics.Append(conf.RequestTraceLogger.As(ctx, &traceLogConf, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	// attempt to unmarshal client debug logging config
	resp.Diagnostics.Append(conf.ClientDebugLogger.As(ctx, &dbgLogConf, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)

	// check for error(s)
	if resp.Diagnostics.HasError() {
		return
	}

	// if request trace logging enabled, configure it
	if traceLogConf.Enabled.ValueBool() {
		osConfig.Logger = client.NewTerraformLogger(ctx, traceLogConf.IncludeRequestBody.ValueBool(), traceLogConf.IncludeResponseBody.ValueBool())
	}
	// if client debug logging enabled, configure it
	if dbgLogConf.Enabled.ValueBool() {
		osConfig.EnableDebugLogger = true
	}

	if osClient, err = opensearch.NewClient(osConfig); err != nil {
		resp.Diagnostics.AddError(
			"Error constructing OpenSearch client",
			fmt.Sprintf("Error occurred constructing OpenSearch client: %v", err.Error()),
		)
		return
	}

	// attempt to perform connectivity and fitment test
	if !conf.SkipInitProductCheck.ValueBool() {
		infoReq := opensearchapi.InfoRequest{}
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if _, err := infoReq.Do(ctx, osClient); err != nil {
			resp.Diagnostics.AddError(
				"Error performing init compatibility check",
				fmt.Sprintf("Error occurred during init compatibility check: %v", err),
			)
			return
		}
	}

	// create shared object for use in resource and datasource types
	shared = Shared{
		Client: osClient,
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
