package usecases

import (
	"time"
)

type EmailConfig struct {
	EmailVerifyLinkExpiresDuration time.Duration
}
