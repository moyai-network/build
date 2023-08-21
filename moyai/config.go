package moyai

import "github.com/df-mc/dragonfly/server"

type Config struct {
	server.UserConfig
}

func DefaultConfig() Config {
	c := Config{
		UserConfig: server.DefaultConfig(),
	}
	return c
}
