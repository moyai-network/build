package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/handlers/user"
	"github.com/moyai-network/build/internal/worlds"
)

// Redo is a command used to place back blocks that were set using Undo.
type Redo struct{}

// Run ...
func (Redo) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	h.Redo()
}

// Allow ...
func (Redo) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && !worlds.InDefaultWorld(p)
}
