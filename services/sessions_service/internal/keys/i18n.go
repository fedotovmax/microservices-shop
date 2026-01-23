package keys

const (
	ValidationUUID          = "validation_uuid"
	ValidationIP            = "validation_ip"
	ValidationStrSymbolsMin = "validation_str_symbols_min"

	ValidationFailed = "validation_failed"
)

const (
	CreateSessionInternal  = "create_session_internal"
	RefreshSessionInternal = "refresh_session_internal"
	VerifyAccessInternal   = "verify_access_internal"
)

const (
	UserNotFound           = "user_not_found"
	UserDeleted            = "user_deleted"
	UserAlreadyExists      = "user_already_exists"
	LoginFromNewIPOrDevice = "login_from_new_ip_or_device"
	BadBypassCode          = "bad_bypass_code"
	BadBlacklistCode       = "bad_blacklist_code"
	UserInBlacklist        = "user_in_blacklist"
)

const (
	SessionNotFound       = "session_not_found"
	InvalidTokenOrExpired = "invalid_token_or_expired"
)
