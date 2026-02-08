package domain

import "time"

type Category struct {
	Children     []*Category
	Translations []Translation
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LogoURL      *string
	ParentID     *string
	Slug         *string
	ID           string
	IsActive     bool
}

func (c *Category) GetParentID() string {
	if c.ParentID != nil {
		return *c.ParentID
	}
	return ""
}

func (c *Category) IsDeleted() bool {
	return c.DeletedAt != nil
}
