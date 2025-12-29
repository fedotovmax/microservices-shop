package domain

import (
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NewAccessToken struct {
	AccessToken   string    `json:"access_token"`
	AccessExpTime time.Time `json:"access_exp_time"`
}

type User struct {
	UID   string
	Email string
}

type BlackList struct {
	Code          string
	CodeExpiresAt time.Time
}

type SessionsUser struct {
	Info      User
	BlackList *BlackList
}

func (u *SessionsUser) IsInBlackList() bool {
	return u.BlackList != nil
}

type Session struct {
	ID             string
	User           SessionsUser
	RefreshHash    string
	IP             string
	Browser        string
	BrowserVersion string
	OS             string
	Device         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RevokedAt      *time.Time
	ExpiresAt      time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

type SessionResponse struct {
	AccessToken    string
	RefreshToken   string
	AccessExpTime  time.Time
	RefreshExpTime time.Time
}

func (r *SessionResponse) ToProto() *sessionspb.CreateSessionResponse {

	return &sessionspb.CreateSessionResponse{
		AccessToken:    r.AccessToken,
		RefreshToken:   r.RefreshToken,
		AccessExpTime:  timestamppb.New(r.AccessExpTime),
		RefreshExpTime: timestamppb.New(r.RefreshExpTime),
	}

}
