package moyai

import "github.com/df-mc/dragonfly/server"

type Config struct {
	server.UserConfig

	Moyai struct {
		Whitelist []string
	}
}

func DefaultConfig() Config {
	c := Config{
		UserConfig: server.DefaultConfig(),
	}
	return c
}
