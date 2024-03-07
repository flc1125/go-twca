package mid

type Config struct {
	Addr       string
	BusinessNo string
	HashKeyNo  string
	HashKey    string
	APIVersion string
}

func (c *Config) init() {
	if c.APIVersion == "" {
		c.APIVersion = "1.0"
	}
}
