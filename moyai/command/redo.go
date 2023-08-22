package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/build/moyai/user"
	"github.com/moyai-network/build/moyai/worlds"
)

// Redo is a command used to place back blocks that were set using Undo.
type Redo struct{}

// Run ...
func (Redo) Run(s cmd.Source, o *cmd.Output) {
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
	return ok && p.World() != worlds.Manager().DefaultWorld()
}
