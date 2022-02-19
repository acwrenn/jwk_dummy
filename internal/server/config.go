package server

type Config struct {
	Protocol string `default:"http"`
	Address  string `default:"localhost"`
	Port     int    `default:"3333"`

	ConfigRoute string `default:"/.well-known/openid-configuration"`
	KeyFile     string `default:""`
}
