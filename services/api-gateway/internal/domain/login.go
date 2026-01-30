package domain

type LoginErrorResponseType uint8

const (
	LoginErrorResponseTypeUserDeleted LoginErrorResponseType = iota
	LoginErrorResponseTypeBadCredentials
	LoginErrorResponseTypeEmailNotVerified
	LoginErrorResponseTypeBadBypassCode
	LoginErrorResponseTypeLoginFromNewDevice
	LoginErrorResponseTypeUserInBlacklist
)

type LoginErrorResponse struct {
	Type    LoginErrorResponseType `json:"type" validate:"required"`
	Message string                 `json:"message" validate:"required"`
	UserId  *string                `json:"user_id"`
}
