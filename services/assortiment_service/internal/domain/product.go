package domain

import "time"

type ProductStatus uint8

const (
	ProductStatusDraft ProductStatus = iota + 1
	ProductStatusActive
	ProductStatusArchived
	ProductStatusDiscontinued
	ProductStatusRestrictedByRightsHolder
)

type ProductCategory struct {
	ID           string
	Slug         string
	Translations []Translation
}

type ProductBrand struct {
	ID    string
	Title string
	Slug  string
}

type Product struct {
	Translations []Translation
	Categories   []ProductCategory
	Brand        ProductBrand
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ID           string
	BrandID      string
	Status       ProductStatus
}

func (p *Product) IsActive() bool {
	return p.Status == ProductStatusActive
}

func (p *Product) IsArchived() bool {
	return p.Status == ProductStatusArchived
}

func (p *Product) IsDeleted() bool {
	return p.DeletedAt != nil
}
