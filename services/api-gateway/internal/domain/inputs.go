package domain

type LoginInput struct {
	Email            string  `json:"email,omitempty" example:"makc-dgek@mail.ru" validate:"required" format:"email"`
	Password         string  `json:"password,omitempty" validate:"required"`
	UserAgent        string  `json:"user_agent,omitempty" validate:"required" example:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"`
	Ip               string  `json:"ip,omitempty" validate:"required" format:"ip4" example:"19.56.186.122"`
	BypassCode       *string `json:"bypass_code,omitempty" validate:"optional" example:"1A2B3C4D5E6F"`
	DeviceTrustToken *string `json:"device_trust_token,omitempty" validate:"optional" example:"skdjfsdfsdifsdfsdf123"`
}
