package command

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/build/moyai/user"
	"github.com/moyai-network/build/moyai/worlds"
	"os"
)

// Paste pastes the user's copied structure.
type Paste struct{}

// PasteExisting pastes an existing structure.
type PasteExisting struct {
	Sub       cmd.SubCommand `cmd:"e"`
	Structure structureList
}

// Run ...
func (Paste) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	h.Paste()
}

// Run ...
func (pe PasteExisting) Run(s cmd.Source, o *cmd.Output) {
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
	return ok && p.World() != worlds.Manager().DefaultWorld()
}

// Allow ...
func (PasteExisting) Allow(s cmd.Source) bool {
	p, ok := s.(*player.Player)
	return ok && p.World() != worlds.Manager().DefaultWorld()
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
		fmt.Println(f.Name())
		st = append(st, f.Name())
	}
	fmt.Println(st)
	return
}
