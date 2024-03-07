package mid

type Config struct {
	Addr       string
	BusinessNo string
	HashKeyNo  string
	HashKey    string
	ApiVersion string
}

func (c *Config) init() {
	if c.ApiVersion == "" {
		c.ApiVersion = "1.0"
	}
}
