package keys

const (
	Development = "development"
	Production  = "production"
)

const AppEnv = "APP_ENV"

var SupportedEnv = map[string]struct{}{
	Development: {},
	Production:  {},
}
