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
}

func NewLoginErrorResponse(t LoginErrorResponseType, m string) LoginErrorResponse {
	return LoginErrorResponse{
		Type:    t,
		Message: m,
	}
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IP           string `json:"ip"`
}
