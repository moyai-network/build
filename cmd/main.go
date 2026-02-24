package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/app"
	"github.com/moyai-network/build/internal/handlers/user"
	"github.com/moyai-network/build/internal/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	config, err := readConfig()
	if err != nil {
		panic(err)
	}
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	c, err := config.Config(log)
	if err != nil {
		panic(err)
	}

	store, db, err := loadStorage(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("close db", "error", err)
		}
	}()

	whitelistEnabled, err := store.WhitelistStore().Enabled(context.Background())
	if err != nil {
		panic(err)
	}

	c.Name = text.Colourf("<b><red>Build</red></b>")
	c.Allower = app.NewAllower(store.WhitelistStore(), whitelistEnabled, log)

	srv := c.New()

	w := srv.World()
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeCreative)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()
	w.StopRaining()

	err = worlds.NewManager(w, "assets/worlds", log)
	if err != nil {
		panic(err)
	}
	registerCommands()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		worlds.Manager().Close()
		if err := srv.Close(); err != nil {
			log.Error("close server", "error", err)
		}
	}()

	srv.Listen()
	for p := range srv.Accept() {
		accept(p)
	}
}

func accept(p *player.Player) {
	p.SetSpeed(5)
	p.SetGameMode(world.GameModeCreative)
	p.ShowCoordinates()
	p.Handle(user.NewHandler(p))
}
