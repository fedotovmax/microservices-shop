package inputs

type CreateCategory struct {
	ParentID     *string
	LogoURL      *string
	Translations []AddCategoryTranslate
}

type AddCategoryTranslate struct {
	LanguageCode string
	Title        string
	Description  *string
}

type UpdateCategoryTranslate struct {
	ID          string
	Title       string
	Description *string
}

type UpdateCategory struct {
	LogoURL  *string
	IsActive *bool
}
