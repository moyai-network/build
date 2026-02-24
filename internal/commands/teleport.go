package commands

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/build/internal/worlds"
)

// TeleportToPos is a command that teleports the user to a position.
type TeleportToPos struct {
	Position mgl64.Vec3 `cmd:"destination"`
}

// TeleportToTarget is a command that teleports the user to another player.
type TeleportToTarget struct {
	Targets []cmd.Target `cmd:"destination"`
}

// TeleportTargetsToTarget is a command that teleports player(s) to another player.
type TeleportTargetsToTarget struct {
	Targets  []cmd.Target `cmd:"target"`
	Position []cmd.Target `cmd:"destination"`
}

// TeleportTargetsToPos is a command that teleports player(s) to a position.
type TeleportTargetsToPos struct {
	Targets  []cmd.Target `cmd:"target"`
	Position mgl64.Vec3   `cmd:"destination"`
}

// Run ...
func (t TeleportToPos) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p := s.(*player.Player)
	p.Teleport(t.Position)
}

// Run ...
func (tp TeleportToTarget) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	t, ok := tp.Targets[0].(*player.Player)
	if !ok {
		return
	}

	dstWorld := worlds.PlayerWorld(t)
	if dstWorld == nil {
		return
	}
	_ = movePlayerTo(p, dstWorld, t.Position())
}

// teleportTargets teleports a list of targets to a specified position and world. If the world is nil, it will only
// teleport to the position. If the position is empty, it will only teleport to the world of the player. It returns the
// players affected in a readable string.
func teleportTargets(targets []cmd.Target, destination mgl64.Vec3, t *player.Player) string {
	dst := worlds.PlayerWorld(t)
	if dst == nil {
		return ""
	}

	for _, target := range targets {
		if p, ok := target.(*player.Player); ok {
			_ = movePlayerTo(p, dst, destination)
		}
	}
	if l := len(targets); l > 1 {
		return fmt.Sprintf("%d players", l)
	}
	return targets[0].(cmd.NamedTarget).Name()
}

func movePlayerTo(p *player.Player, dst *world.World, pos mgl64.Vec3) bool {
	if p == nil || p.Tx() == nil || dst == nil {
		return false
	}
	if p.Tx().World() == dst {
		p.Teleport(pos)
		return true
	}

	handle := p.Tx().RemoveEntity(p)
	if handle == nil {
		handle = p.H()
	}

	<-dst.Exec(func(tx *world.Tx) {
		e := tx.AddEntity(handle)
		moved, ok := e.(*player.Player)
		if !ok {
			return
		}
		moved.Teleport(pos)
	})
	return true
}
