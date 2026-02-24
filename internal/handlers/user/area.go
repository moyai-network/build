package user

// Area represents a 3D cuboid selection where Min/Max are inclusive.
type Area struct {
	Min [3]int
	Max [3]int
}

func NewArea(x1, y1, z1, x2, y2, z2 int) Area {
	return Area{
		Min: [3]int{min(x1, x2), min(y1, y2), min(z1, z2)},
		Max: [3]int{max(x1, x2), max(y1, y2), max(z1, z2)},
	}
}

func (a Area) Dx() int {
	return a.Max[0] - a.Min[0] + 1
}

func (a Area) Dy() int {
	return a.Max[1] - a.Min[1] + 1
}

func (a Area) Dz() int {
	return a.Max[2] - a.Min[2] + 1
}

func (a Area) Range(f func(x, y, z int)) {
	for x := a.Min[0]; x <= a.Max[0]; x++ {
		for y := a.Min[1]; y <= a.Max[1]; y++ {
			for z := a.Min[2]; z <= a.Max[2]; z++ {
				f(x, y, z)
			}
		}
	}
}
