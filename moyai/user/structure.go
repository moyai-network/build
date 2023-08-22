package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// buildStructure is basically the same as world.(*World).BuildStructure.
// The only difference is that here, we have a way of logging all the blocks that are being placed.
func buildStructure(w *world.World, pos cube.Pos, s world.Structure) map[world.Block][]cube.Pos {
	undo := map[world.Block][]cube.Pos{}

	dim := s.Dimensions()
	width, height, length := dim[0], dim[1], dim[2]
	maxX, maxY, maxZ := pos[0]+width, pos[1]+height, pos[2]+length

	for chunkX := pos[0] >> 4; chunkX <= maxX>>4; chunkX++ {
		for chunkY := pos[1] >> 4; chunkY <= maxY>>4; chunkY++ {
			for chunkZ := pos[2] >> 4; chunkZ <= maxZ>>4; chunkZ++ {
				// We approach this on a per-chunk basis, so that we can keep only one chunk in memory at a time
				// while not needing to acquire a new chunk lock for every block. This also allows us not to send
				// block updates, but instead send a single chunk update once.
				baseX, baseY, baseZ := chunkX<<4, chunkY<<4, chunkZ<<4

				for localY := 0; localY < 16; localY++ {
					yOffset := baseY + localY
					if yOffset > w.Range()[1] || yOffset >= maxY {
						// We've hit the height limit for blocks.
						break
					} else if yOffset < w.Range()[0] || yOffset < pos[1] {
						// We've got a block below the minimum, but other blocks might still reach above
						// it, so don't break but continue.
						continue
					}
					for localX := 0; localX < 16; localX++ {
						xOffset := baseX + localX
						if xOffset < pos[0] || xOffset >= maxX {
							continue
						}
						for localZ := 0; localZ < 16; localZ++ {
							zOffset := baseZ + localZ
							if zOffset < pos[2] || zOffset >= maxZ {
								continue
							}
							b, _ := s.At(xOffset-pos[0], yOffset-pos[1], zOffset-pos[2], nil)
							if b != nil {
								ps := cube.Pos{xOffset, yOffset, zOffset}
								bl := w.Block(ps)
								undo[bl] = append(undo[bl], ps)

								w.SetBlock(ps, b, nil)
							}
						}
					}
				}
			}
		}
	}
	return undo
}
