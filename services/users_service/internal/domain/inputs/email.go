package inputs

type EmailVerifyNotificationInput struct {
	title       string
	description string
	locale      string
}

func NewEmailVerifyNotificationInput() *EmailVerifyNotificationInput {
	return &EmailVerifyNotificationInput{}
}

func (i *EmailVerifyNotificationInput) GetTitle() string {
	return i.title
}

func (i *EmailVerifyNotificationInput) GetDescription() string {
	return i.description
}

func (i *EmailVerifyNotificationInput) GetLocale() string {
	return i.locale
}

func (i *EmailVerifyNotificationInput) SetTitle(title string) {
	i.title = title
}

func (i *EmailVerifyNotificationInput) SetDescription(description string) {
	i.description = description
}

func (i *EmailVerifyNotificationInput) SetLocale(locale string) {
	i.locale = locale
}
