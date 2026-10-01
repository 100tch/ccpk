package packer

type Config struct {
	Verbose bool
}

func DefaultConfig() Config {
	return Config{}
}

func (c Config) Validate() error {
	return nil
}
