package command

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/moyai/user"
	"github.com/moyai-network/build/moyai/worlds"
	"strings"
)

// Set is a command used to set a specific block in the area associated with the player.
type Set struct {
	Block blockList `cmd:"block"`
}

// Run ...
func (se Set) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}
	a, ok := h.Area()
	if !ok {
		o.Error("You need to have selected the two area boundaries in order to use this.")
		return
	}

	bl, ok := world.BlockByName("minecraft:"+string(se.Block), nil)
	if !ok {
		o.Errorf("No block with the name %s was found.", se.Block)
		return
	}
	a.Range(func(x, y, z int) {
		p.World().SetBlock(cube.Pos{x, y, z}, bl, nil)
	})
}

// Allow ...
func (Set) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && p.World() != worlds.Manager().DefaultWorld()
}

type (
	blockList string
)

// Type ...
func (blockList) Type() string {
	return "block_list"
}

// Options ...
func (blockList) Options(_ cmd.Source) (bl []string) {
	i := 0
	for {
		i++
		b, ok := world.BlockByRuntimeID(uint32(i))
		if !ok {
			return
		}

		enc, m := b.EncodeBlock()
		if len(m) > 0 {
			continue
		}
		bl = append(bl, strings.Split(enc, "minecraft:")[1])
	}
}
