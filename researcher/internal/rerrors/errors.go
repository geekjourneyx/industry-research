package rerrors

const (
	CodeMissingAPIKey          = "missing_api_key"
	CodeInvalidArgument        = "invalid_argument"
	CodeProviderHTTPError      = "provider_http_error"
	CodeProviderAuthError      = "provider_auth_error"
	CodeProviderQuotaExhausted = "provider_quota_exhausted"
	CodeProviderRateLimited    = "provider_rate_limited"
	CodeProviderTimeout        = "provider_timeout"
	CodeProviderUnavailable    = "provider_unavailable"
	CodeProviderParseError     = "provider_parse_error"
	CodeNoRetrievalTriggered   = "no_retrieval_triggered"
	CodePartialFailure         = "partial_failure"
)

const (
	ExitSuccess              = 0
	ExitInvalidArguments     = 1
	ExitMissingCredentials   = 2
	ExitProviderFailed       = 3
	ExitProviderRateLimited  = 4
	ExitTimeout              = 5
	ExitPartialMultiProvider = 6
)
