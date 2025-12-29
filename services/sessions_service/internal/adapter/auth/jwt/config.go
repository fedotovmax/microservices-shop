package jwtadapter

import "time"

type Config struct {
	AccessTokenExpDuration time.Duration
	AccessTokenSecret      string
}
