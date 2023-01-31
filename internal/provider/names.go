package provider

import "fmt"

const (
	providerName = "opensearch"
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

const (
	resourceSuffixSecurityPluginRole = "security_plugin_role"
	resourceSuffixSecurityPluginUser = "security_plugin_user"
)

const (
	resourceAttrAllowedActions          = "allowed_actions"
	resourceAttrAttributes              = "attributes"
	resourceAttrBackendRoles            = "backend_roles"
	resourceAttrClusterPermissions      = "cluster_permissions"
	resourceAttrDescription             = "description"
	resourceAttrDLS                     = "dls"
	resourceAttrFLS                     = "fls"
	resourceAttrHash                    = "hash"
	resourceAttrHidden                  = "hidden"
	resourceAttrIndexPatterns           = "index_patterns"
	resourceAttrIndexPermissions        = "index_permissions"
	resourceAttrMaskedFields            = "masked_fields"
	resourceAttrOpenDistroSecurityRoles = "open_distro_security_roles"
	resourceAttrPassword                = "password"
	resourceAttrReserved                = "reserved"
	resourceAttrRoleName                = "role_name"
	resourceAttrRoles                   = "roles"
	resourceAttrStatic                  = "static"
	resourceAttrTenantPatterns          = "tenant_patterns"
	resourceAttrTenantPermissions       = "tenant_permissions"
	resourceAttrUsername                = "username"
)

func makeResourceName(providerName, suffix string) string {
	return fmt.Sprintf("%s_%s", providerName, suffix)
}

func makeDatasourceName(providerName, suffix string) string {
	return fmt.Sprintf("%s_%s", providerName, suffix)
}
