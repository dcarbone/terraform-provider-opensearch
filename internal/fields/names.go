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
	ConfigAttrRetryOnStatus         = "retry_on_status"
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
	ResourceTypeSecurityPluginRole = "security_plugin_role"
	ResourceTypeSecurityPluginUser = "security_plugin_user"
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

func TypeName(providerName, typeName string) string {
	const f = "%s_%s"
	return fmt.Sprintf(f, providerName, typeName)
}

func ResourceTypeFQN(providerName, typeName, resourceName string) string {
	const f = "%s.%s"
	return fmt.Sprintf(f, TypeName(providerName, typeName), resourceName)
}

func DatasourceTypeFQN(providerName, typeName, datasourceName string) string {
	const f = "data.%s.%s"
	return fmt.Sprintf(f, TypeName(providerName, typeName), datasourceName)
}
