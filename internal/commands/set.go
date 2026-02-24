package commands

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/handlers/user"
	"github.com/moyai-network/build/internal/worlds"
)

// Set is a command used to set a specific block in the area associated with the player.
type Set struct {
	Block blockList `cmd:"block"`
}

// Run ...
func (se Set) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	bl, ok := world.BlockByName("minecraft:"+string(se.Block), nil)
	if !ok {
		o.Errorf("No block with the name %s was found.", se.Block)
		return
	}
	h.Set(bl)
}

// Allow ...
func (Set) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && !worlds.InDefaultWorld(p)
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
