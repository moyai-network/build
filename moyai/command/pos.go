package command

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/build/moyai/user"
	"github.com/moyai-network/build/moyai/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strconv"
)

type Pos struct {
	Pos pos `cmd:"selection"`
}

func (po Pos) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	pos := cube.PosFromVec3(p.Position())

	h := p.Handler().(*user.Handler)
	n, _ := strconv.Atoi(string(po.Pos))
	h.SetPos(n-1, pos)

	p.Message(text.Colourf("<green>Area position %s set to <yellow>%v, %v, %v</yellow>", po.Pos, pos.X(), pos.Y(), pos.Z()))
}

// Allow ...
func (Pos) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && p.World() != worlds.Manager().DefaultWorld()
}

type pos string

func (pos) Type() string {
	return "pos"
}

func (pos) Options(_ cmd.Source) []string {
	return []string{"1", "2"}
}
