package user

import (
	"os"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/structure"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/build/internal/worlds"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Handler struct {
	player.NopHandler

	p  *player.Player
	mu sync.Mutex

	selection [2]cube.Pos

	undo, redo map[world.Block][]cube.Pos
	copy       structure.Structure
}

func NewHandler(p *player.Player) *Handler {
	return &Handler{
		p: p,

		undo: map[world.Block][]cube.Pos{},
		redo: map[world.Block][]cube.Pos{},
	}
}

var lastUse = time.Now()

func (h *Handler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	if time.Now().Before(lastUse.Add(time.Second)) {
		return
	}
	lastUse = time.Now()

	if worlds.InDefaultWorld(h.p) {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place, break or interact with blocks in the default world.</red>"))
		return
	}

	held, _ := h.p.HeldItems()
	if _, ok := held.Value("WAND"); ok {
		ctx.Cancel()

		h.mu.Lock()
		h.selection[1] = pos
		h.mu.Unlock()

		if a, ok := h.Area(); ok {
			h.p.Message(text.Colourf("<green>Selected Blocks: %d</green>", a.Dx()*a.Dy()*a.Dz()))
		}

		h.p.Message(text.Colourf("<green>Area position 2 set to <yellow>%v, %v, %v</yellow>", pos.X(), pos.Y(), pos.Z()))
	}
}

func (h *Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, _ *[]item.Stack, _ *int) {
	if worlds.InDefaultWorld(h.p) {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place, break or interact with blocks in the default world.</red>"))
		return
	}

	held, _ := h.p.HeldItems()
	if _, ok := held.Value("WAND"); ok {
		ctx.Cancel()

		h.mu.Lock()
		h.selection[0] = pos
		h.mu.Unlock()

		if a, ok := h.Area(); ok {
			h.p.Message(text.Colourf("<green>Selected Blocks: %d</green>", a.Dx()*a.Dy()*a.Dz()))
		}

		h.p.Message(text.Colourf("<green>Area position 1 set to <yellow>%v, %v, %v</yellow>", pos.X(), pos.Y(), pos.Z()))
	}
}
func (h *Handler) HandleBlockPlace(ctx *player.Context, _ cube.Pos, _ world.Block) {
	if worlds.InDefaultWorld(h.p) {
		ctx.Cancel()
		h.p.Message(text.Colourf("<red>You may not place, break or interact with blocks in the default world.</red>"))
		return
	}
}

func (h *Handler) SetPos(n int, pos cube.Pos) {
	h.selection[n] = pos

	if a, ok := h.Area(); ok {
		h.p.Message(text.Colourf("<green>Selected Blocks: %d</green>", a.Dx()*a.Dy()*a.Dz()))
	}
}

func (h *Handler) Area() (Area, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.selection[0] == (cube.Pos{}) || h.selection[1] == (cube.Pos{}) {
		return Area{}, false
	}
	first, second := h.selection[0], h.selection[1]
	return NewArea(first.X(), first.Y(), first.Z(), second.X(), second.Y(), second.Z()), true
}

func (h *Handler) Undo() {
	tx := h.p.Tx()
	if tx == nil {
		return
	}

	h.mu.Lock()
	blocks := h.undo

	h.redo = map[world.Block][]cube.Pos{}
	for b, p := range blocks {
		for _, pos := range p {
			bl := tx.Block(pos)
			h.redo[bl] = append(h.redo[bl], pos)
			tx.SetBlock(pos, b, nil)
		}
	}
	h.mu.Unlock()
}

func (h *Handler) Redo() {
	tx := h.p.Tx()
	if tx == nil {
		return
	}

	h.mu.Lock()
	blocks := h.redo

	h.undo = map[world.Block][]cube.Pos{}
	for b, p := range blocks {
		for _, pos := range p {
			bl := tx.Block(pos)
			h.undo[bl] = append(h.undo[bl], pos)
			tx.SetBlock(pos, b, nil)
		}
	}
	h.mu.Unlock()
}

func (h *Handler) Set(b world.Block) {
	a, ok := h.Area()
	if !ok {
		h.p.Message(text.Colourf("<red>You need to have selected the two area boundaries in order to use this.</red>"))
		return
	}
	tx := h.p.Tx()
	if tx == nil {
		return
	}

	var count int

	h.mu.Lock()
	h.undo = map[world.Block][]cube.Pos{}
	a.Range(func(x, y, z int) {
		count++

		pos := cube.Pos{x, y, z}
		bl := tx.Block(pos)
		h.undo[bl] = append(h.undo[bl], pos)

		tx.SetBlock(cube.Pos{x, y, z}, b, nil)
	})
	h.mu.Unlock()

	h.p.Message(text.Colourf("<green>%d blocks were set.</green>", count))
}

func (h *Handler) Replace(old, new world.Block) {
	a, ok := h.Area()
	if !ok {
		h.p.Message(text.Colourf("<red>You need to have selected the two area boundaries in order to use this.</red>"))
		return
	}
	tx := h.p.Tx()
	if tx == nil {
		return
	}

	var count int

	h.mu.Lock()
	h.undo = map[world.Block][]cube.Pos{}
	a.Range(func(x, y, z int) {
		count++

		pos := cube.Pos{x, y, z}
		bl := tx.Block(pos)
		h.undo[bl] = append(h.undo[bl], pos)

		if bl == old {
			tx.SetBlock(cube.Pos{x, y, z}, new, nil)
		}
	})
	h.mu.Unlock()

	h.p.Message(text.Colourf("<green>%d blocks were replaced.</green>", count))
}

func (h *Handler) Paste() {
	if h.copy == (structure.Structure{}) {
		h.p.Message(text.Colourf("<red>You must use copy first, in order to use this.</red>"))
		return
	}
	tx := h.p.Tx()
	if tx == nil {
		return
	}

	h.undo = map[world.Block][]cube.Pos{}

	h.undo = buildStructure(tx, cube.PosFromVec3(h.p.Position()), h.copy)
	h.p.Message(text.Colourf("<green>Successfully pasted your copied structure.</green>"))
}

func (h *Handler) PasteExisting(name string) {
	path := "assets/structures/" + name

	if _, err := os.Stat(path); os.IsNotExist(err) {
		h.p.Message(text.Colourf("<red>Structure with the name %s does not exist.</red>", name))
		return
	}

	s, err := structure.ReadFile(path)
	if err != nil {
		h.p.Message(text.Colourf("<red>Error trying to load structure: %s.</red>", err))
	}

	if s == (structure.Structure{}) {
		h.p.Message(text.Colourf("<red>Structure with the name %s does not exist.</red>", name))
		return
	}

	tx := h.p.Tx()
	if tx == nil {
		return
	}
	h.undo = buildStructure(tx, cube.PosFromVec3(h.p.Position()), s)
	h.p.Message(text.Colourf("<green>Successfully pasted your copied structure.</green>"))
}

func (h *Handler) Copy() (structure.Structure, bool) {
	a, ok := h.Area()
	if !ok {
		h.p.Message(text.Colourf("<red>You need to have selected the two area boundaries in order to use this.</red>"))
		return structure.Structure{}, false
	}
	tx := h.p.Tx()
	if tx == nil {
		return structure.Structure{}, false
	}
	s := structure.New([3]int{a.Dx(), a.Dy(), a.Dz()})
	var count int

	a.Range(func(x, y, z int) {
		count++
		pos := cube.Pos{x, y, z}

		s.Set(x-a.Min[0], y-a.Min[1], z-a.Min[2], tx.Block(pos), nil)
	})

	// We do this to avoid an unsafe error in the df-mc structure library.
	_ = structure.WriteFile("assets/tmp", s)
	h.copy, _ = structure.ReadFile("assets/tmp")

	h.p.Message(text.Colourf("<green>%d blocks were copied.</green>", count))
	return s, true
}
