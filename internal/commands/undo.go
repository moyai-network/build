package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/handlers/user"
	"github.com/moyai-network/build/internal/worlds"
)

// Undo is a command used to place back blocks that were set using Set or Redo.
type Undo struct{}

// Run ...
func (Undo) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	h.Undo()
}

// Allow ...
func (Undo) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && !worlds.InDefaultWorld(p)
}
