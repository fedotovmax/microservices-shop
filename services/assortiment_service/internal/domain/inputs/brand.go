package inputs

type CreateBrand struct {
	Title       string
	Description *string
	LogoURL     *string
}

type UpdateBrand struct {
	Title       *string
	Description *string
	LogoURL     *string
	IsActive    *bool
}
