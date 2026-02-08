package domain

import "time"

type Badge struct {
	Translations []Translation
	EndsAt       *time.Time
	StartsAt     *time.Time
	Color        *string
	Code         string
	Priority     uint8
}
