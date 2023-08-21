package command

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/build/moyai/worlds"
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
func (w WorldCreate) Run(s cmd.Source, o *cmd.Output) {
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
	wr.SetBlock(wr.Spawn().Sub(cube.Pos{0, 1, 0}), block.Grass{}, nil)
	wr.AddEntity(p)
	p.Teleport(wr.Spawn().Vec3Middle())

	o.Print(text.Colourf("<green>Successfully created world %s.</green>", w.Name))
}

// Run ...
func (w WorldDelete) Run(s cmd.Source, o *cmd.Output) {
	err := worlds.Manager().DeleteWorld(string(w.Name))
	if err != nil {
		o.Error(err)
		return
	}
	o.Print(text.Colourf("<green>Successfully deleted world %s.</green>", w.Name))
}

// Run ...
func (w WorldTeleport) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	wr, ok := worlds.Manager().World(string(w.Name))
	if !ok {
		o.Errorf("No world with the name %s was found", w.Name)
		return
	}
	wr.AddEntity(p)
	p.Teleport(wr.Spawn().Vec3Middle())

	o.Print(text.Colourf("<green>You have been teleported to the world %s.</green>", w.Name))
}

// Allow ...
func (w WorldCreate) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

// Allow ...
func (w WorldDelete) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

// Allow ...
func (w WorldTeleport) Allow(s cmd.Source) bool {
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
