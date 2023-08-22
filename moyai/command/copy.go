package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/structure"
	"github.com/moyai-network/build/moyai/user"
	"github.com/moyai-network/build/moyai/worlds"
	"os"
)

// Copy copies the blocks within the player's selected area.
type Copy struct{}

// CopySave copies the blocks within the player's selected area and saves it to a file.
type CopySave struct {
	Sub  cmd.SubCommand `cmd:"save"`
	Name string         `cmd:"name"`
}

// CopyDelete deletes a structure copy.
type CopyDelete struct {
	Sub  cmd.SubCommand `cmd:"delete"`
	Name structureList  `cmd:"name"`
}

// Run ...
func (Copy) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	h.Copy()
}

// Run ...
func (c CopySave) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	st, ok := h.Copy()
	if !ok {
		return
	}
	path := "assets/structures/" + c.Name

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		o.Errorf("Structure with the name %s already exists.", c.Name)
		return
	}

	err := structure.WriteFile(path, st)
	if err != nil {
		o.Errorf("Error trying to write structure: %s.", err)
		return
	}
}

// Run ...
func (c CopyDelete) Run(s cmd.Source, o *cmd.Output) {
	_ = os.Remove("assets/structures/" + string(c.Name))
}

// Allow ...
func (Copy) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && p.World() != worlds.Manager().DefaultWorld()
}

// Allow ...
func (CopySave) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && p.World() != worlds.Manager().DefaultWorld()
}

// Allow ...
func (CopyDelete) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && p.World() != worlds.Manager().DefaultWorld()
}
