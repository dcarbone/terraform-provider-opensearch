package fields

import "fmt"

const (
	ProviderName = "opensearch"
)

const (
	ConfigAttrAddresses             = "addresses"
	ConfigAttrUsername              = "username"
	ConfigAttrPassword              = "password"
	ConfigAttrCACert                = "ca_cert"
	ConfigAttrRetryOnStatuses       = "retry_on_statuses"
	ConfigAttrDisableRetry          = "disable_retry"
	ConfigAttrEnableRetryOnTimeout  = "enable_retry_on_timeout"
	ConfigAttrMaxRetries            = "max_retries"
	ConfigAttrCompressRequestBody   = "compress_request_body"
	ConfigAttrInsecureSkipTLSVerify = "insecure_skip_tls_verify"
	ConfigAttrUseResponseCheckOnly  = "use_response_check_only"
	ConfigAttrSkipInitProductCheck  = "skip_init_product_check"
	ConfigAttrLogging               = "logging"
	ConfigAttrEnabled               = "enabled"
	ConfigAttrIncludeRequestBody    = "include_request_body"
	ConfigAttrIncludeResponseBody   = "include_response_body"
)

const (
	ResourceSuffixSecurityPluginRole = "security_plugin_role"
	ResourceSuffixSecurityPluginUser = "security_plugin_user"
)

const (
	ResourceAttrAllowedActions          = "allowed_actions"
	ResourceAttrAttributes              = "attributes"
	ResourceAttrBackendRoles            = "backend_roles"
	ResourceAttrClusterPermissions      = "cluster_permissions"
	ResourceAttrDescription             = "description"
	ResourceAttrDLS                     = "dls"
	ResourceAttrFLS                     = "fls"
	ResourceAttrHash                    = "hash"
	ResourceAttrHidden                  = "hidden"
	ResourceAttrIndexPatterns           = "index_patterns"
	ResourceAttrIndexPermissions        = "index_permissions"
	ResourceAttrMaskedFields            = "masked_fields"
	ResourceAttrOpenDistroSecurityRoles = "open_distro_security_roles"
	ResourceAttrPassword                = "password"
	ResourceAttrReserved                = "reserved"
	ResourceAttrRoleName                = "role_name"
	ResourceAttrRoles                   = "roles"
	ResourceAttrStatic                  = "static"
	ResourceAttrTenantPatterns          = "tenant_patterns"
	ResourceAttrTenantPermissions       = "tenant_permissions"
	ResourceAttrUsername                = "username"
)

func MakeResourceName(providerName, suffix string) string {
	return fmt.Sprintf("%s_%s", providerName, suffix)
}

func MakeDatasourceName(providerName, suffix string) string {
	return fmt.Sprintf("%s_%s", providerName, suffix)
}
