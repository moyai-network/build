package main

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/moyai"
	"github.com/moyai-network/build/moyai/command"
	"github.com/moyai-network/build/moyai/user"
	"github.com/moyai-network/build/moyai/worlds"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	c, err := config.Config(log)
	if err != nil {
		log.Fatalln(err)
	}
	c.Name = text.Colourf("<b><red>Build</red></b>")
	c.Allower = moyai.NewAllower(config.Moyai.Whitelist)

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
		log.Fatalln(err)
	}
	registerCommands()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		worlds.Manager().Close()
		if err := srv.Close(); err != nil {
			log.Errorf("close server: %v", err)
		}
	}()

	srv.Listen()
	for srv.Accept(accept) {
		// Do nothing
	}
}

func accept(p *player.Player) {
	p.SetGameMode(world.GameModeCreative)
	p.ShowCoordinates()
	p.Handle(user.NewHandler(p))
}

// registerCommands registers all commands for the build server.
func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("world", "Manage worlds.", []string{"w"}, command.WorldCreate{}, command.WorldDelete{}, command.WorldTeleport{}),
		cmd.New("wand", "Get the magic wand", nil, command.Wand{}),
		cmd.New("set", "Set blocks within your area selection.", nil, command.Set{}),
		cmd.New("undo", "Undo your set / redo usage.", nil, command.Undo{}),
		cmd.New("redo", "Redo your undo usage.", nil, command.Redo{}),
		cmd.New("copy", "Copy a structure.", nil, command.Copy{}, command.CopySave{}, command.CopyDelete{}),
		cmd.New("paste", "Paste your copied structure.", nil, command.Paste{}, command.PasteExisting{}),
	} {
		cmd.Register(c)
	}
}

// readConfig reads the config file and returns a moyai.Config.
// If the config file does not exist, one is generated.
func readConfig() (moyai.Config, error) {
	c := moyai.DefaultConfig()
	g := gophig.NewGophig("./config", "toml", 0777)

	err := g.GetConf(&c)
	if os.IsNotExist(err) {
		err = g.SetConf(c)
	}
	return c, err
}
