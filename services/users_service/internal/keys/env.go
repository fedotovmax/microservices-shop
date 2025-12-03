package keys

const (
	Local       = "local"
	Development = "development"
	Production  = "production"
)

const AppEnv = "APP_ENV"

var SupportedEnv = map[string]struct{}{
	Local:       {},
	Development: {},
	Production:  {},
}
