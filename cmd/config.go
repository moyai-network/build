package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/moyai-network/build/internal/app"
	"github.com/restartfu/gophig"
)

// readConfig reads the config file and returns an app.Config.
// If the config file does not exist, one is generated.
func readConfig() (app.Config, error) {
	c := app.DefaultConfig()
	g := gophig.NewGophig("./config", "toml", 0o777)

	err := g.GetConf(&c)
	if os.IsNotExist(err) {
		err = g.SetConf(c)
	}
	if err != nil {
		return c, err
	}
	if err := applyDatabaseEnvOverrides(&c); err != nil {
		return c, err
	}
	return c, err
}

func applyDatabaseEnvOverrides(cfg *app.Config) error {
	if cfg == nil {
		return nil
	}
	applyEnvString("DB_HOST", &cfg.Database.Host)
	applyEnvString("DB_NAME", &cfg.Database.Name)
	applyEnvString("DB_USER", &cfg.Database.User)
	applyEnvPassword(&cfg.Database.Password)

	if err := applyEnvInt("DB_PORT", &cfg.Database.Port); err != nil {
		return err
	}
	if err := applyEnvInt32("DB_MAX_CONNS", &cfg.Database.MaxConns); err != nil {
		return err
	}
	return nil
}

func applyEnvString(key string, target *string) {
	if target == nil {
		return
	}
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return
	}
	*target = value
}

func applyEnvPassword(target *string) {
	if target == nil {
		return
	}
	if pass, ok := os.LookupEnv("DB_PASS"); ok {
		*target = pass
		return
	}
	if pass, ok := os.LookupEnv("DB_PASSWORD"); ok {
		*target = pass
	}
}

func applyEnvInt(key string, target *int) error {
	if target == nil {
		return nil
	}
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("invalid %s %q: %w", key, raw, err)
	}
	*target = value
	return nil
}

func applyEnvInt32(key string, target *int32) error {
	if target == nil {
		return nil
	}
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("invalid %s %q: %w", key, raw, err)
	}
	*target = int32(value)
	return nil
}
