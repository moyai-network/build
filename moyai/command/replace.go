package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/moyai/user"
)

type Replace struct {
	Old blockList `cmd:"old"`
	New blockList `cmd:"new"`
}

func (r Replace) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	oldBlock, ok := world.BlockByName("minecraft:"+string(r.Old), nil)
	if !ok {
		o.Errorf("No block with the name %s was found.", r.Old)
		return
	}

	newBlock, ok := world.BlockByName("minecraft:"+string(r.New), nil)
	if !ok {
		o.Errorf("No block with the name %s was found.", r.New)
		return
	}

	h.Replace(oldBlock, newBlock)
}
