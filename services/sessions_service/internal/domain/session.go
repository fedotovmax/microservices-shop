package domain

import (
	"net"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NewAccessToken struct {
	AccessToken   string    `json:"access_token"`
	AccessExpTime time.Time `json:"access_exp_time"`
}

type User struct {
	UID       string
	Email     string
	DeletedAt *time.Time
}

type BlackList struct {
	Code          string
	CodeExpiresAt time.Time
}

type Bypass struct {
	Code            string
	BypassExpiresAt time.Time
}

type SessionsUser struct {
	Info      User
	BlackList *BlackList
	Bypass    *Bypass
	//TODO:
	twoFactor bool
}

func (u *SessionsUser) Clone() *SessionsUser {
	var bl *BlackList

	if u.BlackList != nil {
		bl = &BlackList{
			Code:          u.BlackList.Code,
			CodeExpiresAt: u.BlackList.CodeExpiresAt,
		}
	}
	var bp *Bypass

	if u.Bypass != nil {
		bp = &Bypass{
			Code:            u.Bypass.Code,
			BypassExpiresAt: u.Bypass.BypassExpiresAt,
		}
	}
	return &SessionsUser{
		Info: User{
			UID:   u.Info.UID,
			Email: u.Info.Email,
		},
		BlackList: bl,
		Bypass:    bp,
	}
}

func (u *SessionsUser) HasTwoFactor() bool {
	return u.twoFactor
}

func (u *SessionsUser) IsDeleted() bool {
	return u.Info.DeletedAt != nil
}

func (u *SessionsUser) IsInBlackList() bool {
	return u.BlackList != nil
}

func (u *SessionsUser) HasBypass() bool {
	return u.Bypass != nil
}

func (bl *BlackList) IsCodeExpired() bool {
	return time.Now().After(bl.CodeExpiresAt)
}

func (bl *BlackList) ComapreCodes(code string) bool {
	return bl.Code == code
}

func (bp *Bypass) IsCodeExpired() bool {
	return time.Now().After(bp.BypassExpiresAt)
}

func (bp *Bypass) ComapreCodes(code string) bool {
	return bp.Code == code
}

type Session struct {
	ID             string
	User           *SessionsUser
	RefreshHash    string
	IP             net.IP
	Browser        string
	BrowserVersion string
	OS             string
	Device         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RevokedAt      *time.Time
	ExpiresAt      time.Time
}

func (s *Session) Clone() Session {

	return Session{
		ID:             s.ID,
		User:           s.User.Clone(),
		RefreshHash:    s.RefreshHash,
		IP:             s.IP,
		Browser:        s.Browser,
		BrowserVersion: s.BrowserVersion,
		OS:             s.OS,
		Device:         s.Device,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		RevokedAt:      s.RevokedAt,
		ExpiresAt:      s.ExpiresAt,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

type PreparedTrustTokenAction int8

const (
	TrustTokenNone PreparedTrustTokenAction = iota
	TrustTokenCreated
	TrustTokenUpdated
)

type PreparedTrustToken struct {
	UID                     string
	DeviceTrustTokenValue   string
	DeviceTrustTokenHash    string
	DeviceTrustTokenExpTime time.Time
	Action                  PreparedTrustTokenAction
}

type SessionResponseTrustToken struct {
	DeviceTrustTokenExpTime time.Time
	DeviceTrustTokenValue   string
}

type SessionResponse struct {
	AccessToken    string
	RefreshToken   string
	AccessExpTime  time.Time
	RefreshExpTime time.Time
	TrustToken     *SessionResponseTrustToken
}

func (r *SessionResponse) ToProto() *sessionspb.CreateSessionResponse {

	var trustToken *sessionspb.CreatedTrustToken

	if r.TrustToken != nil {
		trustToken = &sessionspb.CreatedTrustToken{
			TrustTokenValue:   r.TrustToken.DeviceTrustTokenValue,
			TrustTokenExpTime: timestamppb.New(r.TrustToken.DeviceTrustTokenExpTime),
		}
	}

	return &sessionspb.CreateSessionResponse{
		AccessToken:    r.AccessToken,
		RefreshToken:   r.RefreshToken,
		AccessExpTime:  timestamppb.New(r.AccessExpTime),
		RefreshExpTime: timestamppb.New(r.RefreshExpTime),
		TrustToken:     trustToken,
	}

}

type DeviceTrustToken struct {
	TokenHash string
	UID       string
	LastUsed  time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func (s *DeviceTrustToken) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *DeviceTrustToken) IsRevoked() bool {
	return s.RevokedAt != nil
}
