package commands

import (
	"strconv"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/handlers/user"
	"github.com/moyai-network/build/internal/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Pos struct {
	Pos pos `cmd:"selection"`
}

func (po Pos) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
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
	return ok && !worlds.InDefaultWorld(p)
}

type pos string

func (pos) Type() string {
	return "pos"
}

func (pos) Options(_ cmd.Source) []string {
	return []string{"1", "2"}
}
