package inputs

import "time"

type UpdateBadge struct {
	EndsAt   *time.Time
	Color    *string
	Code     string
	Priority *uint8
}

type BadgeTranslate struct {
	LanguageCode string
	Title        string
}

type SaveBadge struct {
	EndsAt   *time.Time
	StartsAt *time.Time
	Color    *string
	Code     string
	Priority uint8
}

type UpdateBadgeTranslate struct {
	ID    string
	Title string
}
