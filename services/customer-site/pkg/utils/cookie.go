package utils

import "net/http"

func CreateCSRFCookie(name, value string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	return cookie
}
