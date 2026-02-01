package domain

import "time"

type EmailVerifyLink struct {
	Link          string
	UserID        string
	LinkExpiresAt time.Time
}

func (l *EmailVerifyLink) IsExpired() bool {
	return time.Now().After(l.LinkExpiresAt)
}
