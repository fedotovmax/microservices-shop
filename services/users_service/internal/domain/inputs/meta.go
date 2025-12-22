package inputs

type MetaParams struct {
	userID string
	locale string
}

func NewMetaParams() *MetaParams {
	return &MetaParams{}
}

func (m *MetaParams) GetLocale() string {
	return m.locale
}

func (m *MetaParams) GetUserID() string {
	return m.userID
}

func (m *MetaParams) SetLocale(l string) {
	m.locale = l
}

func (m *MetaParams) SetUserID(userID string) {
	m.userID = userID
}
