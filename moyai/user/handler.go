package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/we/geo"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/build/moyai/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

type Handler struct {
	player.NopHandler

	p *player.Player

	selection [2]cube.Pos
}

func NewHandler(p *player.Player) *Handler {
	return &Handler{
		p: p,
	}
}

var lastUse = time.Now()

func (h *Handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	if time.Now().Before(lastUse.Add(time.Second)) {
		return
	}
	lastUse = time.Now()

	if h.p.World() == worlds.Manager().DefaultWorld() {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place, break or interact with blocks in the default world.</red>"))
		return
	}

	held, _ := h.p.HeldItems()
	if _, ok := held.Value("WAND"); ok {
		ctx.Cancel()
		h.selection[1] = pos
		h.p.Message(text.Colourf("<green>Area position 2 set to <yellow>%v, %v, %v</yellow>", pos.X(), pos.Y(), pos.Z()))
	}
}

func (h *Handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, _ *[]item.Stack, _ *int) {
	if h.p.World() == worlds.Manager().DefaultWorld() {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place, break or interact with blocks in the default world.</red>"))
		return
	}

	held, _ := h.p.HeldItems()
	if _, ok := held.Value("WAND"); ok {
		ctx.Cancel()
		h.selection[0] = pos
		h.p.Message(text.Colourf("<green>Area position 1 set to <yellow>%v, %v, %v</yellow>", pos.X(), pos.Y(), pos.Z()))
	}
}
func (h *Handler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	if h.p.World() == worlds.Manager().DefaultWorld() {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place, break or interact with blocks in the default world.</red>"))
		return
	}
}

func (h *Handler) Area() (geo.Area, bool) {
	if h.selection[0] == (cube.Pos{}) || h.selection[1] == (cube.Pos{}) {
		return geo.Area{}, false
	}
	first, second := h.selection[0], h.selection[1]
	return geo.NewArea(first.X(), first.Y(), first.Z(), second.X(), second.Y(), second.Z()), true
}
