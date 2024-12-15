package main

import (
	"fmt"
	"os"
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
	Position
	warehouse *Warehouse
}

type Wall struct{ WarehouseItem }

type Box struct{ WarehouseItem }

type Item interface {
	getIcon() rune
	setPosition(p Position)
	getPosition() Position
	canMove(m *Movement) bool
}

type Warehouse struct {
	width, height int
	itemMap       map[Position]Item
	walls         []*Wall
	boxes         []*Box

	submarine Position
}

func (w *Wall) getIcon() rune {
	return '#'
}

func (b *Box) getIcon() rune {
	return 'O'
}

func (w *Wall) getPosition() Position {
	return w.Position
}

func (b *Box) getPosition() Position {
	return b.Position
}

func (w *Wall) setPosition(p Position) {
	w.Position = p
}

func (b *Box) setPosition(p Position) {
	b.Position = p
}

func (w *Wall) canMove(m *Movement) bool {
	return false
}

func (w *Box) canMove(m *Movement) bool {
	obj, found := w.warehouse.itemMap[w.Position.displace(m)]
	return !found || obj.canMove(m)
}

func (w *Warehouse) addWall(wall *Wall) {
	w.walls = append(w.walls, wall)
	w.itemMap[wall.Position] = wall
}

func (w *Warehouse) addBox(box *Box) {
	w.boxes = append(w.boxes, box)
	w.itemMap[box.Position] = box
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
				sb.WriteRune('@')
			} else if obj, ok := w.itemMap[position]; ok {
				sb.WriteRune(obj.getIcon())
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (w *Warehouse) moveSubmarine(m *Movement) bool {
	newCandidatePosition := w.submarine.displace(m)
	if obj, ok := w.itemMap[newCandidatePosition]; !ok {
		w.submarine = newCandidatePosition
		return true
	} else {
		if obj.canMove(m) {
			stack := BoxStack{}
			for obj != nil {
				stack.Push(obj)
				obj, ok = w.itemMap[obj.getPosition().displace(m)]
			}
			for len(stack) > 0 {
				obj := stack.Pop()
				oldPos := obj.getPosition()
				newPos := oldPos.displace(m)
				obj.setPosition(newPos)
				delete(w.itemMap, oldPos)
				w.itemMap[newPos] = obj
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
		width:   width,
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
						Position: Position{
							x: c,
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
						Position: Position{
							x: c,
							y: r,
						},
						warehouse: &warehouse,
					},
				}
				warehouse.addWall(&wall)
			}
			if o == '@' {
				warehouse.submarine = Position{x: c, y: r}
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
	input, _ := os.ReadFile("input.txt")

	inputSections := strings.Split(string(input), "\n\n")

	warehouse := newWarehouse(inputSections[0])
	movements := getMovements(inputSections[1])

	for _, m := range movements {
		warehouse.moveSubmarine(m)
	}

	gpsSum := 0

	for _, b := range warehouse.boxes {
		gpsSum += b.Position.y*100 + b.Position.x
	}

	fmt.Println(warehouse)
	fmt.Println(gpsSum)
	// for pos, wo := range warehouse.itemMap {
	// 	if isWall(*wo) {
	// 		fmt.Println("Wall in position", pos)
	// 	} else if isBox(*wo) {
	// 		fmt.Println("Box in position", pos)
	// 	}
	// }
}
