package domain

import "time"

type Translation struct {
	ID           string
	Title        string
	LanguageCode string
	Description  *string
}

type Category struct {
	ID           string
	Slug         string
	ParentID     *string
	LogoURL      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Children     []*Category
	Translations []Translation
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
