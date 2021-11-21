package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Area int

const (
	AREA_GRASS = iota
	AREA_CITY
	AREA_ROAD
)

var (
	// The index of this array is also the side it extends to! So 0 == left, 1 == down, 2 == right, 3 == up
	g_sides []Pos = []Pos{Pos{-1, 0}, Pos{0, 1}, Pos{1, 0}, Pos{0, -1}}
)

type Meeple struct {
	sideIndex   int
	playerIndex int
}

type Tile struct {
	id       int
	sides    [4]Area
	cloister bool
	emblem   bool
	// the uint8 is basically a set of indices (bits) 0..3 shift index
	connections []uint8
	meeple      Meeple
}

type Player struct {
	index   int
	score   int
	meeples int
}

type Pos struct {
	x int
	y int
}

type Move struct {
	tile Tile
	pos  Pos
}

func (p Pos) String() string {
	return fmt.Sprintf("Pos(%2d,%2d)", p.x, p.y)
}

func (m Meeple) String() string {
	return fmt.Sprintf("Meeple(side: %v, player: %v)", m.sideIndex, m.playerIndex)
}

func (p Player) String() string {
	return fmt.Sprintf("Player(id: %v, score: %v, meeples: %v)", p.index, p.score, p.meeples)
}

func (area Area) String() string {
	return [...]string{"Grass", "City", "Road"}[area]
}

func (t Tile) String() string {
	cloister := ""
	if t.cloister {
		cloister = " Cloister"
	}
	sides := fmt.Sprintf("[%-5v %-5v %-5v %-5v]", t.sides[0], t.sides[1], t.sides[2], t.sides[3])
	emblem := ""
	if t.emblem {
		emblem = " Emblem"
	}
	conn0 := ""
	if len(t.connections) > 0 {
		conn0 = strconv.FormatInt(int64(t.connections[0]), 2)
	}
	conn1 := ""
	if len(t.connections) > 1 {
		conn1 = strconv.FormatInt(int64(t.connections[1]), 2)
	}

	return fmt.Sprintf("Tile(%v %04v %04v%v%v)", sides, conn0, conn1, cloister, emblem)
}

func add(p, p2 Pos) Pos {
	p.x += p2.x
	p.y += p2.y
	return p
}

func multiplyTile(tiles *[]Tile, t Tile, count int) {
	for i := 0; i < count; i++ {
		*tiles = append(*tiles, t)
	}
}

func newSet(indices []int) (res map[int]bool) {
	for i := range indices {
		res[i] = true
	}
	return
}

func getTiles() (Tile, []Tile) {

	var tiles []Tile
	var id int

	multiplyTile(&tiles, Tile{id, [4]Area{AREA_GRASS, AREA_ROAD, AREA_GRASS, AREA_GRASS}, true, false, nil, Meeple{-1, -1}}, 2)
	id++
	multiplyTile(&tiles, Tile{id, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, nil, Meeple{-1, -1}}, 4)
	id++
	multiplyTile(&tiles, Tile{id, [4]Area{AREA_CITY, AREA_CITY, AREA_CITY, AREA_CITY}, false, true, []uint8{15}, Meeple{-1, -1}}, 1)
	id++
	multiplyTile(&tiles, Tile{id, [4]Area{AREA_ROAD, AREA_GRASS, AREA_ROAD, AREA_CITY}, false, false, []uint8{5}, Meeple{-1, -1}}, 3)
	id++

	startTile := Tile{id, [4]Area{AREA_ROAD, AREA_GRASS, AREA_ROAD, AREA_CITY}, false, false, []uint8{5}, Meeple{-1, -1}}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })

	return startTile, tiles
}

func rotateTile(tile Tile) Tile {

	last := tile.sides[3]
	for i := 3; i >= 1; i-- {
		tile.sides[i] = tile.sides[i-1]
	}
	tile.sides[0] = last

	for i, s := range tile.connections {
		tile.connections[i] = s<<1&0xf | s>>3
	}
	tile.meeple.sideIndex = tile.meeple.sideIndex % 4
	return tile
}

func placementPossible(board map[Pos]Tile, tile Tile, pos Pos) bool {
	if v, ok := board[add(pos, Pos{-1, 0})]; ok && tile.sides[0] != v.sides[2] {
		return false
	}
	if v, ok := board[add(pos, Pos{0, 1})]; ok && tile.sides[1] != v.sides[3] {
		return false
	}
	if v, ok := board[add(pos, Pos{1, 0})]; ok && tile.sides[2] != v.sides[0] {
		return false
	}
	if v, ok := board[add(pos, Pos{0, -1})]; ok && tile.sides[3] != v.sides[1] {
		return false
	}
	return true
}

func generatePossibleMoves(board map[Pos]Tile, tiles []Tile, openPlacements map[Pos]bool, player Player) (moves []Move) {

	for place := range openPlacements {

		alreadyPlaced := make(map[int]bool)
		for _, t := range tiles {
			if _, ok := alreadyPlaced[t.id]; ok {
				continue
			}
			alreadyPlaced[t.id] = true

			for rot := 0; rot < 4; rot++ {
				if rot > 0 {
					t = rotateTile(t)
				}
				if placementPossible(board, t, place) {
					moves = append(moves, Move{t, place})

					if player.meeples > 0 {
						for side := 0; side < 4; side++ {
							if t.sides[side] != AREA_GRASS {
								t.meeple = Meeple{side, player.index}
								moves = append(moves, Move{t, place})
							}
						}
					}
				}
			}
		}
	}

	return
}

func placeTile(board *map[Pos]Tile, openPlacements *map[Pos]bool, players *[]Player, tile Tile, pos Pos) {
	(*board)[pos] = tile
	delete(*openPlacements, pos)
	if tile.meeple.playerIndex != -1 {
		(*players)[tile.meeple.playerIndex].meeples -= 1
	}

	for _, s := range g_sides {
		if _, ok := (*board)[add(pos, s)]; !ok {
			(*openPlacements)[add(pos, s)] = true
		}
	}

}

// returns (pointCount, listOfPositions, cityIsClosed) for the streets/cities. The cityIsClosed can just be ignored for roads!
// sideFrom is the index on the current tile (not the last) that needs to be further investigated.
// foundMeeples = map[playerIndex]countOfFoundMeeples
func _calcNewPoints(board map[Pos]Tile, pos Pos, sideFrom int, areaType Area, searched *map[Pos]bool) (int, []Pos, bool, map[int]int) {

	score := 0
	positions := []Pos{}
	cityIsClosed := false
	foundMeeples := map[int]int{}

	// There is no tile on the board or it doesn't make sense
	if t, ok := board[pos]; !ok || t.sides[sideFrom] == AREA_GRASS || t.sides[sideFrom] != areaType {
		cityIsClosed = false
		return score, positions, cityIsClosed, foundMeeples
	}
	// We already visited this tile
	if _, ok := (*searched)[pos]; ok {
		cityIsClosed = true
		return score, positions, cityIsClosed, foundMeeples
	}

	(*searched)[pos] = true
	score += 1
	positions = append(positions, pos)
	cityIsClosed = true
	tile := board[pos]
	if tile.meeple.playerIndex != -1 {
		foundMeeples[tile.meeple.playerIndex] = 1
	}

	for _, c := range tile.connections {
		// if there is a connection from the incoming side to somewhere
		if c>>sideFrom&1 == 1 {
			for sideIndex := 0; sideIndex < 4; sideIndex++ {
				// Check only, if there actually is a connection.
				if c>>sideIndex&1 == 1 && sideIndex != sideFrom {
					new_score, new_positions, closed, meeples := _calcNewPoints(board, add(pos, g_sides[sideIndex]), (sideIndex+2)%4, areaType, searched)
					score += new_score
					positions = append(positions, new_positions...)
					cityIsClosed = cityIsClosed && closed
					for playerIndex, count := range meeples {
						foundMeeples[playerIndex] += count
					}
				}
			}
		}
	}

	return score, positions, cityIsClosed, foundMeeples
}

func closingSide(tile Tile, index int) bool {
	if tile.sides[index] == AREA_GRASS {
		return false
	}
	if len(tile.connections) == 0 {
		return true
	}
	for _, connection := range tile.connections {
		if connection>>index&1 == 1 {
			return false
		}
	}
	return true
}

func closingTileSides(board map[Pos]Tile, pos Pos) (closingSides []int) {

	tile := board[pos]

	skipSides := map[int]bool{}

	for sideIndex, _ := range g_sides {
		if _, ok := skipSides[sideIndex]; ok {
			continue
		}
		if closingSide(tile, sideIndex) {
			closingSides = append(closingSides, sideIndex)
		} else {
			// the tile has at least one connection to another side. So check, if all sides have tiles on the board!
			// If not, then this is not a closing tile!
			allSidesAreOnBoard := true
			for _, connection := range tile.connections {
				if connection>>sideIndex&1 == 1 {
					for i := 0; i < 4; i++ {
						if connection>>i&1 == 1 {
							_, ok := board[add(pos, g_sides[i])]
							allSidesAreOnBoard = allSidesAreOnBoard && ok
							skipSides[i] = true
						}
					}
				}
			}
			if allSidesAreOnBoard {
				closingSides = append(closingSides, sideIndex)
			}
		}
	}
	return
}

// Sums up points for finished roads or finished cities for the given move
func calcNewPoints(board *map[Pos]Tile, pos Pos, players *[]Player) {

	// This should only check the recursive points, iff either of the following conditions is true
	// - a road is ended with this tile
	// - a city is closed with this tile
	// - there is a connection between sides and on all ends are already tiles on the board
	// Additionally, check if a cloister is closed within the surrounding tiles.
	// Those optimizations are ignored right now. This might make it a bit faster later on! Right now, we just evaluate
	// all directions if it's not grass!

	tile := (*board)[pos]

	for _, sideIndex := range closingTileSides(*board, pos) {
		searched := map[Pos]bool{}
		score, positions, closed, meeples := _calcNewPoints(*board, add(pos, g_sides[sideIndex]), (sideIndex+2)%4, tile.sides[sideIndex], &searched)

		if !closed {
			continue
		}

		// Who has the most meeples on the board? Is there a clear winner?
		bestPlayer, bestCount, secondBestCount := -1, 0, 0
		for playerIndex, count := range meeples {
			if bestPlayer == -1 || count > bestCount {
				secondBestCount = bestCount
				bestCount = count
				bestPlayer = playerIndex
			}
		}
		if bestCount > secondBestCount {
			(*players)[bestPlayer].score += score
		}
		// Clean up and remove meeples from the board. Add them back to the players inventory!
		for _, p := range positions {
			if t := (*board)[p]; t.meeple.playerIndex != -1 {
				(*players)[t.meeple.playerIndex].meeples += 1
				t.meeple = Meeple{-1, -1}
				(*board)[p] = t
			}
		}
	}

}

func main() {

	board := make(map[Pos]Tile)
	startTile, tiles := getTiles()
	openPlacements := map[Pos]bool{Pos{-1, 0}: true, Pos{1, 0}: true, Pos{0, -1}: true, Pos{0, 1}: true}

	playerCount := 3
	var players []Player
	for i := 0; i < playerCount; i++ {
		players = append(players, Player{i, 0, 6})
	}
	board[Pos{0, 0}] = startTile

	for round := 0; round < 3; round++ {
		for _, player := range players {
			moves := generatePossibleMoves(board, tiles, openPlacements, player)
			move := moves[rand.Intn(len(moves))]
			placeTile(&board, &openPlacements, &players, move.tile, move.pos)
			calcNewPoints(&board, move.pos, &players)
		}
	}

	for p, t := range board {
		fmt.Printf("%v --> %v\n", p, t)
	}
	for _, p := range players {
		fmt.Println(p)
	}

}
