package keys

const FallbackLocale = "ru"

const (
	ValidationPassword   = "validation_password"
	ValidationEmail      = "validation_email"
	ValidationUUID       = "validation_uuid"
	ValidationNumRange   = "validation_num_range"
	ValidationNumMin     = "validation_num_min"
	ValidationNumMax     = "validation_num_max"
	ValidationNumBetween = "validation_num_between"

	ValidationStrFilePath     = "validation_str_filepath"
	ValidationStrSymbolsRange = "validation_str_symbols_range"
	ValidationStrSymbolsMax   = "validation_str_symbols_max"
	ValidationStrSymbolsMin   = "validation_str_symbols_min"

	ValidationGender     = "validation_gender"
	ValidationFullName   = "validation_fullname"
	ValidationDateFormat = "validation_date_format"

	ValidationFailed = "validation_failed"
)

const (
	UserNotFound      = "user_not_found"
	UserAlreadyExists = "user_already_exists"

	GetUserInternal           = "get_user_internal"
	CreateUserInternal        = "create_user_internal"
	UpdateUserProfileInternal = "update_user_profile_internal"
	UserSessionActionInternal = "user_session_action_internal"
)

const (
	UserDeleted          = "user_deleted"
	UserBadCredentials   = "bad_credentials"
	UserEmailNotVerified = "email_not_verified"
)
