package utils

import (
	"errors"
	"strings"

	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
)

var emptyHeaderErr = errors.New("empty authorization header")
var badHeaderFormatErr = errors.New("bad authorization header format")

func ValidateAuthHeader(header string) (string, error) {
	if header == "" {
		return "", emptyHeaderErr
	}

	authHeaderParts := strings.Split(header, " ")

	if len(authHeaderParts) != 2 {
		return "", badHeaderFormatErr
	}

	if authHeaderParts[0] != keys.BearerAuthorizationToken {
		return "", badHeaderFormatErr
	}

	return authHeaderParts[1], nil
}
