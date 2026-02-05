package domain

import "time"

type Brand struct {
	ID          string
	Title       string
	Slug        string
	Description *string
	LogoURL     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	IsActive    bool
}

func (b *Brand) IsDeleted() bool {
	return b.DeletedAt != nil
}
