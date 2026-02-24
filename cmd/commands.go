package main

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/build/internal/commands"
)

// registerCommands registers all commands for the build server.
func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("tp", "teleportation commands", nil, commands.TeleportToPos{}, commands.TeleportToTarget{}),
		cmd.New("world", "Manage worlds.", []string{"w"}, commands.WorldCreate{}, commands.WorldDelete{}, commands.WorldTeleport{}),
		cmd.New("wand", "Get the magic wand", nil, commands.Wand{}),
		cmd.New("set", "Set blocks within your area selection.", nil, commands.Set{}),
		cmd.New("replace", "Replace blocks within your area selection.", nil, commands.Replace{}),
		cmd.New("undo", "Undo your set / redo usage.", nil, commands.Undo{}),
		cmd.New("redo", "Redo your undo usage.", nil, commands.Redo{}),
		cmd.New("copy", "Copy a structure.", nil, commands.Copy{}, commands.CopySave{}, commands.CopyDelete{}),
		cmd.New("paste", "Paste your copied structure.", nil, commands.Paste{}, commands.PasteExisting{}),
		cmd.New("pos", "Update a selection using your player position.", nil, commands.Pos{}),
	} {
		cmd.Register(c)
	}
}
