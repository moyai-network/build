package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Wand is a command used in order to get the magic world-edit wand.
type Wand struct{}

func (Wand) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	_, _ = p.Inventory().AddItem(item.NewStack(item.Axe{Tier: item.ToolTierWood}, 1).WithValue("WAND", true))
	o.Print(text.Colourf("<green>You have been given the magic wand.</green>"))
}

// Allow ...
func (Wand) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
