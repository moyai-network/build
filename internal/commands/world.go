package commands

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// WorldCreate creates a new world with the given name.
type WorldCreate struct {
	Sub  cmd.SubCommand `cmd:"create"`
	Name string         `cmd:"name"`
}

// WorldDelete deletes the world with the given name.
type WorldDelete struct {
	Sub  cmd.SubCommand `cmd:"delete"`
	Name worldList      `cmd:"name"`
}

// WorldTeleport is a command used to teleport to a given world.
type WorldTeleport struct {
	Sub  cmd.SubCommand `cmd:"tp"`
	Name worldList      `cmd:"name"`
}

// Run ...
func (w WorldCreate) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	if _, ok := worlds.Manager().World(w.Name); ok {
		o.Errorf("A world with the name %s already exists.", w.Name)
		return
	}

	wr, err := worlds.Manager().CreateWorld(w.Name)
	if err != nil {
		o.Error(err)
		return
	}
	<-wr.Exec(func(tx *world.Tx) {
		tx.SetBlock(tx.World().Spawn().Sub(cube.Pos{0, 1, 0}), block.Grass{}, nil)
	})
	if !worlds.MovePlayer(p, wr) {
		o.Error("failed to move player to the new world")
		return
	}

	o.Print(text.Colourf("<green>Successfully created world %s.</green>", w.Name))
}

// Run ...
func (w WorldDelete) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	err := worlds.Manager().DeleteWorld(string(w.Name))
	if err != nil {
		o.Error(err)
		return
	}
	o.Print(text.Colourf("<green>Successfully deleted world %s.</green>", w.Name))
}

// Run ...
func (w WorldTeleport) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	wr, ok := worlds.Manager().World(string(w.Name))
	if !ok {
		o.Errorf("No world with the name %s was found", w.Name)
		return
	}
	if !worlds.MovePlayer(p, wr) {
		o.Error("failed to teleport player to the selected world")
		return
	}

	o.Print(text.Colourf("<green>You have been teleported to the world %s.</green>", w.Name))
}

// Allow ...
func (WorldCreate) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

// Allow ...
func (WorldDelete) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

// Allow ...
func (WorldTeleport) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

type (
	// worldList represents the world list enum type for commands.
	worldList string
)

// Type ...
func (worldList) Type() string {
	return "world_list"
}

// Options ...
func (worldList) Options(_ cmd.Source) (wl []string) {
	for _, w := range worlds.Manager().Worlds() {
		wl = append(wl, w.Name())
	}
	return
}
