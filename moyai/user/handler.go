package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/build/moyai/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Handler struct {
	player.NopHandler

	p *player.Player
}

func NewHandler(p *player.Player) *Handler {
	return &Handler{
		p: p,
	}
}

func (h Handler) HandleBlockBreak(ctx *event.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	if h.p.World() == worlds.Manager().DefaultWorld() {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place or break blocks in the default world.</red>"))
		return
	}
}
func (h Handler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	if h.p.World() == worlds.Manager().DefaultWorld() {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place or break blocks in the default world.</red>"))
		return
	}
}
