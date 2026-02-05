package inputs

type CreateCategory struct {
	ParentID     *string
	LogoURL      *string
	Translations []CategoryTranslate
}

type CategoryTranslate struct {
	LanguageCode string
	Title        string
	Description  *string
}

type UpdateCategory struct {
	Title    *string
	LogoURL  *string
	IsActive *bool
}
