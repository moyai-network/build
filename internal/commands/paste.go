package commands

import (
	"os"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/internal/handlers/user"
	"github.com/moyai-network/build/internal/worlds"
)

// Paste pastes the user's copied structure.
type Paste struct{}

// PasteExisting pastes an existing structure.
type PasteExisting struct {
	Sub       cmd.SubCommand `cmd:"e"`
	Structure structureList
}

// Run ...
func (Paste) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	h.Paste()
}

// Run ...
func (pe PasteExisting) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	h.PasteExisting(string(pe.Structure))
}

// Allow ...
func (Paste) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && !worlds.InDefaultWorld(p)
}

// Allow ...
func (PasteExisting) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && !worlds.InDefaultWorld(p)
}

type (
	structureList string
)

func (structureList) Type() string {
	return "structure_list"
}

func (structureList) Options(cmd.Source) (st []string) {
	dir, err := os.ReadDir("assets/structures/")
	if err != nil {
		return
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		st = append(st, f.Name())
	}
	return
}
