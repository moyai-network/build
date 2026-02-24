package app

import (
	"net/url"
	"strconv"

	"github.com/df-mc/dragonfly/server"
)

type Config struct {
	server.UserConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string

	MaxConns int32
}

func (c DatabaseConfig) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + strconv.Itoa(c.Port),
		Path:   "/" + c.Name,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

func DefaultConfig() Config {
	c := Config{
		UserConfig: server.DefaultConfig(),
	}
	c.Network.Address = "0.0.0.0:19169"
	c.Database.Port = 5432
	c.Database.MaxConns = 10
	return c
}
