package block

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

// transparent is a struct that may be embedded to make a block transparent to light. Light will be able to
// pass through this block freely.
type transparent struct{}

// LightDiffusionLevel ...
func (transparent) LightDiffusionLevel() uint8 {
	return 0
}

// sourceWaterDisplacer may be embedded to allow displacing water source blocks.
type sourceWaterDisplacer struct{}

// CanDisplace returns true if the world.Liquid passed is of the type Water, not falling and has a depth of 8.
func (s sourceWaterDisplacer) CanDisplace(b world.Liquid) bool {
	w, ok := b.(block.Water)
	return ok && !w.Falling && w.Depth == 8
}

// replaceableWith checks if the block at the position passed is replaceable with the block passed.
func replaceableWith(w *world.World, pos cube.Pos, with world.Block) bool {
	if pos.OutOfBounds(w.Range()) {
		return false
	}
	b := w.Block(pos)
	if replaceable, ok := b.(block.Replaceable); ok {
		return replaceable.ReplaceableBy(with) && b != with
	}
	return false
}

// firstReplaceable finds the first replaceable block position eligible to have a block placed on it after
// clicking on the position and face passed.
// If none can be found, the bool returned is false.
func firstReplaceable(w *world.World, pos cube.Pos, face cube.Face, with world.Block) (cube.Pos, cube.Face, bool) {
	if replaceableWith(w, pos, with) {
		// A replaceableWith block was clicked, so we can replace it. This will then be assumed to be placed on
		// the top face. (Torches, for example, will get attached to the floor when clicking tall grass.)
		return pos, cube.FaceUp, true
	}
	side := pos.Side(face)
	if replaceableWith(w, side, with) {
		return side, face, true
	}
	return pos, face, false
}

// place places the block passed at the position passed. If the user implements the block.Placer interface, it
// will use its PlaceBlock method. If not, the block is placed without interaction from the user.
func place(w *world.World, pos cube.Pos, b world.Block, user item.User, ctx *item.UseContext) {
	if placer, ok := user.(block.Placer); ok {
		placer.PlaceBlock(pos, b, ctx)
		return
	}
	w.SetBlock(pos, b, nil)
	w.PlaySound(pos.Vec3(), sound.BlockPlace{Block: b})
}

// placed checks if an item was placed with the use context passed.
func placed(ctx *item.UseContext) bool {
	return ctx.CountSub > 0
}

// newBreakInfo creates a BreakInfo struct with the properties passed. The XPDrops field is 0 by default. The blast
// resistance is set to the block's hardness*5 by default.
func newBreakInfo(hardness float64, harvestable func(item.Tool) bool, effective func(item.Tool) bool, drops func(item.Tool, []item.Enchantment) []item.Stack) block.BreakInfo {
	return block.BreakInfo{
		Hardness:        hardness,
		BlastResistance: hardness * 5,
		Harvestable:     harvestable,
		Effective:       effective,
		Drops:           drops,
	}
}

// pickaxeEffective is a convenience function for blocks that are effectively mined with a pickaxe.
var pickaxeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypePickaxe
}

// axeEffective is a convenience function for blocks that are effectively mined with an axe.
var axeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeAxe
}

// shearsEffective is a convenience function for blocks that are effectively mined with shears.
var shearsEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeShears
}

// shovelEffective is a convenience function for blocks that are effectively mined with a shovel.
var shovelEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeShovel
}

// hoeEffective is a convenience function for blocks that are effectively mined with a hoe.
var hoeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeHoe
}

// nothingEffective is a convenience function for blocks that cannot be mined efficiently with any tool.
var nothingEffective = func(item.Tool) bool {
	return false
}

// alwaysHarvestable is a convenience function for blocks that are harvestable using any item.
var alwaysHarvestable = func(t item.Tool) bool {
	return true
}

// neverHarvestable is a convenience function for blocks that are not harvestable by any item.
var neverHarvestable = func(t item.Tool) bool {
	return false
}

// pickaxeHarvestable is a convenience function for blocks that are harvestable using any kind of pickaxe.
var pickaxeHarvestable = pickaxeEffective

// oneOf returns a drops function that returns one of each of the item types passed.
func oneOf(i ...world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(item.Tool, []item.Enchantment) []item.Stack {
		var s []item.Stack
		for _, it := range i {
			s = append(s, item.NewStack(it, 1))
		}
		return s
	}
}

// fallDistanceEntity is an entity that has a fall distance.
type fallDistanceEntity interface {
	// ResetFallDistance resets the entities fall distance.
	ResetFallDistance()
	// FallDistance returns the entities fall distance.
	FallDistance() float64
}

// boolByte returns 1 if the bool passed is true, or 0 if it is false.
func boolByte(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// horizontalDirection returns the horizontal direction of the given direction. This is a legacy type still used in
// various blocks.
func horizontalDirection(d cube.Direction) cube.Direction {
	switch d {
	case cube.South:
		return cube.North
	case cube.West:
		return cube.South
	case cube.North:
		return cube.West
	case cube.East:
		return cube.East
	}
	panic("invalid direction")
}
