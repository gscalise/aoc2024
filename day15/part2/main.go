package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Movement struct {
	dx, dy int
}

var UP Movement = Movement{0, -1}
var DOWN Movement = Movement{0, 1}
var LEFT Movement = Movement{-1, 0}
var RIGHT Movement = Movement{1, 0}

type Position struct {
	x, y int
}

func (p Position) displace(m *Movement) Position {
	return Position{
		x: p.x + m.dx,
		y: p.y + m.dy,
	}
}

type WarehouseItem struct {
	position1 Position
	position2 Position
	warehouse *Warehouse
}

type Wall struct{ WarehouseItem }

type Box struct{ WarehouseItem }

type Item interface {
	getIcon() string
	setPosition(p1 Position, p2 Position)
	getPositions() [2]Position
	canMove(m *Movement) bool
	gatherMovementCandidateBoxes(m *Movement) map[*Box]bool
}

type Warehouse struct {
	width, height int
	itemMap       map[Position]Item
	walls         []*Wall
	boxes         []*Box

	submarine Position
}

func (w *Wall) getIcon() string {
	return "##"
}

func (b *Box) getIcon() string {
	return "[]"
}

func (w *Wall) getPositions() [2]Position {
	return [...]Position{w.position1, w.position2}
}

func (b *Box) getPositions() [2]Position {
	return [...]Position{b.position1, b.position2}
}

func (w *Wall) setPosition(p1 Position, p2 Position) {
	w.position1 = p1
	w.position2 = p2
}

func (b *Box) setPosition(p1 Position, p2 Position) {
	b.position1 = p1
	b.position2 = p2
}

func (w *Wall) canMove(m *Movement) bool {
	return false
}

func (w *Box) canMove(m *Movement) bool {
	if *m == RIGHT {
		obj, found := w.warehouse.itemMap[w.position2.displace(m)]
		return !found || obj.canMove(m)
	} else if *m == LEFT {
		obj, found := w.warehouse.itemMap[w.position1.displace(m)]
		return !found || obj.canMove(m)
	} else {
		obj1, found1 := w.warehouse.itemMap[w.position1.displace(m)]
		obj2, found2 := w.warehouse.itemMap[w.position2.displace(m)]
		if !found1 && !found2 {
			return true
		} else if found1 && found2 {
			if obj1 == obj2 {
				return obj1.canMove(m)
			} else {
				return obj1.canMove(m) && obj2.canMove(m)
			}
		} else if found1 {
			return obj1.canMove(m)
		} else {
			return obj2.canMove(m)
		}
	}
}

func (w *Wall) gatherMovementCandidateBoxes(m *Movement) map[*Box]bool {
	return make(map[*Box]bool)
}

func (w *Box) gatherMovementCandidateBoxes(m *Movement) map[*Box]bool {
	candidates := map[*Box]bool{w: true}
	if *m == RIGHT {
		obj, found := w.warehouse.itemMap[w.position2.displace(m)]
		if found {
			for kb, _ := range obj.gatherMovementCandidateBoxes(m) {
				candidates[kb] = true
			}
		}
	} else if *m == LEFT {
		obj, found := w.warehouse.itemMap[w.position1.displace(m)]
		if found {
			for kb, _ := range obj.gatherMovementCandidateBoxes(m) {
				candidates[kb] = true
			}
		}
	} else {
		obj1, found1 := w.warehouse.itemMap[w.position1.displace(m)]
		obj2, found2 := w.warehouse.itemMap[w.position2.displace(m)]
		if !found1 && !found2 {
		} else if found1 && found2 {
			if obj1 == obj2 {
				for kb, _ := range obj1.gatherMovementCandidateBoxes(m) {
					candidates[kb] = true
				}
			} else {
				for kb, _ := range obj1.gatherMovementCandidateBoxes(m) {
					candidates[kb] = true
				}
				for kb, _ := range obj2.gatherMovementCandidateBoxes(m) {
					candidates[kb] = true
				}
			}
		} else if found1 {
			for kb, _ := range obj1.gatherMovementCandidateBoxes(m) {
				candidates[kb] = true
			}
		} else {
			for kb, _ := range obj2.gatherMovementCandidateBoxes(m) {
				candidates[kb] = true
			}
		}
	}
	return candidates
}

func (w *Warehouse) addWall(wall *Wall) {
	w.walls = append(w.walls, wall)
	w.itemMap[wall.position1] = wall
	w.itemMap[wall.position2] = wall
}

func (w *Warehouse) addBox(box *Box) {
	w.boxes = append(w.boxes, box)
	w.itemMap[box.position1] = box
	w.itemMap[box.position2] = box
}

type BoxStack []Item

func (s *BoxStack) Push(v Item) {
	*s = append(*s, v)
}

func (s *BoxStack) Pop() Item {
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]

	return ret
}

func (w *Warehouse) String() string {
	var sb strings.Builder

	for r := range w.height {
		for c := range w.width {
			position := Position{x: c, y: r}
			if w.submarine == position {
				sb.WriteString("@")
			} else if obj, ok := w.itemMap[position]; ok {
				if position == obj.getPositions()[0] {
					sb.WriteString(obj.getIcon())
				}
			} else {
				sb.WriteString(".")
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (w *Warehouse) moveSubmarine(m *Movement) bool {
	newCandidatePosition := w.submarine.displace(m)
	var obj Item = nil
	var ok = false

	if obj, ok = w.itemMap[newCandidatePosition]; !ok {
		w.submarine = newCandidatePosition
		return true
	} else {
		if obj.canMove(m) {

			stack := BoxStack{}
			if m == &LEFT || m == &RIGHT {
				for ok {
					stack.Push(obj)
					if m == &LEFT {
						obj, ok = w.itemMap[obj.getPositions()[0].displace(m)]
					} else {
						obj, ok = w.itemMap[obj.getPositions()[1].displace(m)]
					}
				}
			} else {
				candidateBoxes := []*Box{}
				for k := range obj.gatherMovementCandidateBoxes(m) {
					candidateBoxes = append(candidateBoxes, k)
				}
				var yComp func(a, b *Box) int = nil
				if m == &UP {
					yComp = func(i1, i2 *Box) int {
						return i2.getPositions()[0].y - i1.getPositions()[0].y
					}
				} else {
					yComp = func(i1, i2 *Box) int {
						return i1.getPositions()[0].y - i2.getPositions()[0].y
					}
				}
				slices.SortFunc(candidateBoxes, yComp)
				for _, b := range candidateBoxes {
					stack.Push(b)
				}
			}
			for len(stack) > 0 {
				obj := stack.Pop()
				oldPos := obj.getPositions()
				newPos := [2]Position{oldPos[0].displace(m), oldPos[1].displace(m)}
				obj.setPosition(newPos[0], newPos[1])
				delete(w.itemMap, oldPos[0])
				delete(w.itemMap, oldPos[1])
				w.itemMap[newPos[0]] = obj
				w.itemMap[newPos[1]] = obj
			}
			w.submarine = newCandidatePosition
			return true
		}
		return false
	}
}

func newWarehouse(input string) *Warehouse {
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	warehouse := Warehouse{
		width:   width * 2,
		height:  height,
		itemMap: map[Position]Item{},
		walls:   []*Wall{},
		boxes:   []*Box{},
	}

	for r, l := range lines {
		for c, o := range l {
			if o == 'O' {
				b := Box{
					WarehouseItem: WarehouseItem{
						position1: Position{
							x: 2 * c,
							y: r,
						},
						position2: Position{
							x: 2*c + 1,
							y: r,
						},
						warehouse: &warehouse,
					},
				}
				warehouse.addBox(&b)
			}
			if o == '#' {
				wall := Wall{
					WarehouseItem: WarehouseItem{
						position1: Position{
							x: 2 * c,
							y: r,
						},
						position2: Position{
							x: 2*c + 1,
							y: r,
						},
						warehouse: &warehouse,
					},
				}
				warehouse.addWall(&wall)
			}
			if o == '@' {
				warehouse.submarine = Position{x: 2 * c, y: r}
			}
		}
	}
	return &warehouse
}

func getMovements(input string) []*Movement {
	movements := []*Movement{}
	for _, l := range input {
		var movement *Movement = nil
		switch l {
		case '^':
			movement = &UP
		case '>':
			movement = &RIGHT
		case '<':
			movement = &LEFT
		case 'v':
			movement = &DOWN
		}
		if movement != nil {
			movements = append(movements, movement)
		}
	}
	return movements
}

func main() {
	input, _ := os.ReadFile("../input.txt")

	inputSections := strings.Split(string(input), "\n\n")

	warehouse := newWarehouse(inputSections[0])
	movements := getMovements(inputSections[1])

	for _, m := range movements {
		warehouse.moveSubmarine(m)
	}

	gpsSum := 0

	for _, b := range warehouse.boxes {
		gpsSum += b.position1.y*100 + b.position1.x
	}

	fmt.Println(gpsSum)

}
