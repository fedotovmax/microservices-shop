package state

import (
	"net/http"

	"github.com/fedotovmax/microservices-shop/customer-site/internal/keys"
)

type ClientState struct {
	AccessToken  string
	RefreshToken string
	UserAgent    string
	IP           string
	DeviceToken  string
	Dirty        bool // true, если токены обновились
}

func GetClientState(r *http.Request) *ClientState {

	state := &ClientState{}

	if c, err := r.Cookie(keys.CookieAccessToken); err == nil {
		state.AccessToken = c.Value
	}
	if c, err := r.Cookie(keys.CookieRefreshToken); err == nil {
		state.RefreshToken = c.Value
	}

	state.UserAgent = r.Header.Get(keys.UserAgentHeader)

	ip := r.Header.Get(keys.XForwadedForHeader)
	if ip == "" {
		ip = r.RemoteAddr
	}
	state.IP = ip

	//TODO: other vars add

	return state
}
