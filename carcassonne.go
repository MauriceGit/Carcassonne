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

const (
	SIDE_LEFT   = 0
	SIDE_DOWN   = 1
	SIDE_RIGHT  = 2
	SIDE_UP     = 3
	SIDE_CENTER = 4
)

var (
	// The index of this array is also the side it extends to! So 0 == left, 1 == down, 2 == right, 3 == up
	g_sides []Pos = []Pos{Pos{-1, 0}, Pos{0, 1}, Pos{1, 0}, Pos{0, -1}}
	// Even though those are not positions per se, it's two ints which is exactly what we need here.
	g_connection_indices = []Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}}
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
	// 4x4 grid of 1/0, if a column is connected to a row index! So 0010 means, that this row_index is connected to side 2.
	connections uint16
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

// To keep track of partially visited tiles
// (city ends on tile where a road goes through. City is already visited, road not yet!)
type Visited struct {
	pos  Pos
	side int
}

func (p Pos) String() string {
	return fmt.Sprintf("Pos(%2d,%2d)", p.x, p.y)
}

func (m Meeple) String() string {
	return fmt.Sprintf("Meeple(side: %v, player: %v)", m.sideIndex, m.playerIndex)
}

func (p Player) String() string {
	cStart, cEnd := playerIndexColor(p.index)
	return fmt.Sprintf("%vPlayer%v(id: %v, score: %v, meeples: %v)", cStart, cEnd, p.index, p.score, p.meeples)
}

func (area Area) String() string {
	return [...]string{"Grass", "City", "Road"}[area]
}

func (area Area) StringShort() string {
	return [...]string{"~", "c", "r"}[area]
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
	conn := strconv.FormatInt(int64(t.connections), 2)
	return fmt.Sprintf("Tile(%v %v %016v%v%v)", t.id, sides, conn, cloister, emblem)
}

func playerIndexColor(i int) (string, string) {
	return fmt.Sprintf("\033[0;%vm", i+31), "\033[0m"
}

func (m Meeple) colorCodeStartEnd() (string, string) {
	return playerIndexColor(m.playerIndex)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Creates an optimized uint16 that represents the connections, from
// a list of actual connection indices
func connectionsToUint16(conns []Pos) (out uint16) {
	for _, c := range conns {
		out |= 1 << (c.x*4 + c.y)
		out |= 1 << (c.y*4 + c.x)
	}
	return
}

func (t Tile) hasConnectionAtSide(i int) bool {
	if t.hasNoConnections() {
		return false
	}
	if t.allSidesAreConnected() {
		return true
	}
	for _, c := range g_connection_indices {
		if (c.x == i || c.y == i) && (t.connections>>(c.x*4))&(1<<c.y) != 0 {
			return true
		}
	}
	return false
}

func (t Tile) allSidesAreConnected() bool {
	return t.connections == 0x7B2E
}

func (t Tile) hasNoConnections() bool {
	return t.connections == 0x0
}

func drawColor(m Meeple, drawColor bool, s string) {
	cStart, cEnd := m.colorCodeStartEnd()

	if drawColor && m.playerIndex != -1 {
		fmt.Printf("%v%v", cStart, s)
	} else {
		fmt.Printf("%v", s)
	}

	fmt.Printf("%v", cEnd)
}

func drawField(board map[Pos]Tile) {

	minX, maxX := 0, 0
	minY, maxY := 0, 0
	for k := range board {
		minX = min(minX, k.x)
		minY = min(minY, k.y)
		maxX = max(maxX, k.x)
		maxY = max(maxY, k.y)
	}

	for y := minY; y <= maxY; y++ {
		for row := 0; row < 3; row++ {
			for x := minX; x <= maxX; x++ {
				if t, ok := board[Pos{x, y}]; ok {
					filler := "~"

					switch row {
					case 0:
						if t.emblem {
							fmt.Printf("#")
						} else {
							fmt.Printf("%v", filler)
						}
						fmt.Printf("%v", filler)
						drawColor(t.meeple, t.meeple.sideIndex == SIDE_UP, fmt.Sprintf("%v", t.sides[SIDE_UP].StringShort()))
						fmt.Printf("%v%v", filler, filler)
					case 1:
						drawColor(t.meeple, t.meeple.sideIndex == SIDE_LEFT, fmt.Sprintf("%v", t.sides[SIDE_LEFT].StringShort()))

						if t.hasConnectionAtSide(SIDE_LEFT) {
							fmt.Printf("%v", t.sides[SIDE_LEFT].StringShort())
						} else {
							fmt.Printf("%v", filler)
						}
						if t.cloister {
							drawColor(t.meeple, t.meeple.sideIndex == SIDE_CENTER, "Ħ")
						} else {
							if !t.hasNoConnections() {
								fmt.Printf("+")
							} else {
								fmt.Printf("%v", filler)
							}
						}
						if t.hasConnectionAtSide(SIDE_RIGHT) {
							fmt.Printf("%v", t.sides[SIDE_RIGHT].StringShort())
						} else {
							fmt.Printf("%v", filler)
						}
						drawColor(t.meeple, t.meeple.sideIndex == SIDE_RIGHT, t.sides[SIDE_RIGHT].StringShort())
					case 2:
						fmt.Printf("%v%v", filler, filler)
						drawColor(t.meeple, t.meeple.sideIndex == SIDE_DOWN, fmt.Sprintf("%v", t.sides[SIDE_DOWN].StringShort()))
						fmt.Printf("%v%v", filler, filler)
					}

				} else {
					fmt.Printf("     ")
				}
			}
			fmt.Println("")
		}
	}
	fmt.Println("")

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

func getTiles() (Tile, []Tile) {

	var tiles []Tile
	var id int

	multiplyTile(&tiles, Tile{id, [4]Area{AREA_GRASS, AREA_ROAD, AREA_GRASS, AREA_GRASS}, true, false, 0, Meeple{-1, -1}}, 2)
	id++
	multiplyTile(&tiles, Tile{id, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, true, false, 0, Meeple{-1, -1}}, 4)
	id++
	conn := connectionsToUint16([]Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}})
	multiplyTile(&tiles, Tile{id, [4]Area{AREA_CITY, AREA_CITY, AREA_CITY, AREA_CITY}, false, true, conn, Meeple{-1, -1}}, 1)
	id++
	conn = connectionsToUint16([]Pos{Pos{0, 2}})
	multiplyTile(&tiles, Tile{id, [4]Area{AREA_ROAD, AREA_GRASS, AREA_ROAD, AREA_CITY}, false, false, conn, Meeple{-1, -1}}, 3)
	id++

	conn = connectionsToUint16([]Pos{Pos{0, 2}})
	startTile := Tile{id, [4]Area{AREA_ROAD, AREA_GRASS, AREA_ROAD, AREA_CITY}, false, false, conn, Meeple{-1, -1}}

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

	// 0x7B2E == 0b0111101100101110 --> All sides are connected!
	if tile.connections != 0 && tile.connections != 0x7B2E {
		a := tile.connections & 0xf
		b := (tile.connections >> 4) & 0xf
		c := (tile.connections >> 8) & 0xf
		d := (tile.connections >> 12) & 0xf

		a = (a>>1 | a<<3) & 0xf
		b = (b>>1 | b<<3) & 0xf
		c = (c>>1 | c<<3) & 0xf
		d = (d>>1 | d<<3) & 0xf

		tile.connections = (a << 12) | (d << 8) | (c << 4) | b
	}

	tile.meeple.sideIndex = (tile.meeple.sideIndex + 1) % 4
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

		// At some point - implement a statistic (remaining tile_type * tile_count / all_tile_count or something)
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
						if t.cloister {
							t.meeple = Meeple{SIDE_CENTER, player.index}
							moves = append(moves, Move{t, place})
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

/*
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

	for _, c := range g_connection_indices {

		if (tile.connections>>(c.x*4))&(1<<c.y) != 0 {
			new_score, new_positions, closed, meeples := _calcNewPoints(board, add(pos, g_sides[c.x]), (c.x+2)%4, areaType, searched)
			score += new_score
			positions = append(positions, new_positions...)
			cityIsClosed = cityIsClosed && closed
			for playerIndex, count := range meeples {
				foundMeeples[playerIndex] += count
			}

		}
	}

	return score, positions, cityIsClosed, foundMeeples
}

func closingSide(tile Tile, index int) bool {
	if tile.sides[index] == AREA_GRASS {
		return false
	}
	if tile.connections == 0 {
		return true
	}

	for _, c := range g_connection_indices {
		if (c.x == index || c.y == index) && (tile.connections>>(c.x*4))&(1<<c.y) != 0 {
			return false
		}
	}
	return true
}

func closingTileSides(board map[Pos]Tile, pos Pos) (closingSides []int) {

	tile := board[pos]

	for sideIndex, _ := range g_sides {
		if tile.sides[sideIndex] == AREA_GRASS {
			continue
		}
		if closingSide(tile, sideIndex) {
			closingSides = append(closingSides, sideIndex)
		} else {
			// the tile has at least one connection to another side. So check, if all sides have tiles on the board!
			// If not, then this is not a closing tile!
			allSidesAreOnBoard := true

			for _, c := range g_connection_indices {
				if (c.x == sideIndex || c.y == sideIndex) && (tile.connections>>(c.x*4))&(1<<c.y) != 0 {
					_, ok1 := board[add(pos, g_sides[c.x])]
					_, ok2 := board[add(pos, g_sides[c.y])]
					allSidesAreOnBoard = allSidesAreOnBoard && ok1 && ok2
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
		if bestPlayer != -1 && bestCount > secondBestCount {
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
*/

// Returns:
// Score, positions_with_meeple_on_them, is_closed
func calcRecursivePoints(board map[Pos]Tile, pos Pos, side int, searched *map[Pos]bool, meeples *[]int) (int, []Pos, bool) {

	// We already visited this tile
	if _, ok := (*searched)[pos]; ok {
		return 0, nil, true
	}

	tile, ok := board[pos]
	// If tile is not even on the board. Do we need this check?
	if !ok {
		return 0, nil, false
	}

	if tile.sides[side] == AREA_GRASS {
		return 0, nil, false
	}

	(*searched)[pos] = true

	score, positions, closed := calcRecursivePoints(board, add(pos, g_sides[side]), (side+2)%4, searched, meeples)

	// This tile also counts
	score += 1
	// An extra point, if we are building a city which has an emblem on it!
	if tile.emblem && tile.sides[side] == AREA_CITY {
		score += 1
	}

	if tile.meeple.playerIndex != -1 && tile.meeple.sideIndex == side {
		positions = append(positions, pos)
		(*meeples)[tile.meeple.playerIndex] += 1
	}

	for _, c := range g_connection_indices {
		if (c.x == side || c.y == side) && (tile.connections>>(c.x*4))&(1<<c.y) != 0 {

			otherSide := c.x
			if c.x == side {
				otherSide = c.y
			}

			if tile.meeple.playerIndex != -1 && tile.meeple.sideIndex == otherSide {
				positions = append(positions, pos)
				(*meeples)[tile.meeple.playerIndex] += 1
			}

			_score, _positions, _closed := calcRecursivePoints(board, add(pos, g_sides[otherSide]), (otherSide+2)%4, searched, meeples)
			score += _score
			positions = append(positions, _positions...)
			closed = closed && _closed
		}
	}

	if !closed {
		positions = nil
	}

	return score, positions, closed

}

// meeples: array with index == player_index and value == meeple_count
func getBestPlayerIndex(meeples []int) int {

	bestPlayer, bestCount := -1, 0
	for playerIndex, count := range meeples {
		if count > bestCount {
			bestCount = count
			bestPlayer = playerIndex
		}
	}
	return bestPlayer

	/*

		// Who has the most meeples on the board? Is there a clear winner?
		bestPlayer, bestCount, secondBestCount := -1, 0, 0
		for playerIndex, count := range meeples {
			if bestPlayer == -1 || count > bestCount {
				secondBestCount = bestCount
				bestCount = count
				bestPlayer = playerIndex
			}
		}
		if bestPlayer != -1 && bestCount > secondBestCount {
			return bestPlayer
		}
		return -1
	*/
}

// TODO: Only remove a meeple, if it was actually on a now-closed structure???
// This should automatically be checked when the positions are searched!
// positions should only be tiles with a meeple on it, that needs to be removed!!!
func cleanupUsedMeeplesFromBoard(board *map[Pos]Tile, players *[]Player, positions []Pos) {
	//fmt.Println(positions)
	// Clean up and remove meeples from the board. Add them back to the players inventory!
	for _, p := range positions {
		t := (*board)[p]
		(*players)[t.meeple.playerIndex].meeples += 1
		t.meeple = Meeple{-1, -1}
		(*board)[p] = t
	}
}

// UpdateFinalPoints adds up the final points that are achieved by placing the tile at pos.
// It only counts finished cities and closed roads1
func updateFinalPoints(board *map[Pos]Tile, pos Pos, players *[]Player) {

	tile := (*board)[pos]

	for side := 0; side < 4; side++ {
		if tile.hasConnectionAtSide(side) {

			// We always skip one side of a connection, so we only search each possible way once!
			ok := false
			for _, c := range g_connection_indices {
				if c.x == side {
					ok = true
				}
			}
			if !ok {
				continue
			}
		}

		searched := map[Pos]bool{}
		meeples := make([]int, len(*players), len(*players))
		score, positions, closed := calcRecursivePoints(*board, pos, side, &searched, &meeples)

		if bestPlayer := getBestPlayerIndex(meeples); closed && bestPlayer != -1 {

			// Closed cities count twice!
			if t, ok := (*board)[pos]; ok && t.sides[side] == AREA_CITY {
				score *= 2
			}

			for playerIndex, count := range meeples {
				if count == meeples[bestPlayer] {
					(*players)[playerIndex].score += score
				}
			}
		}

		cleanupUsedMeeplesFromBoard(board, players, positions)
	}
}

// Assembles a set of all positions with meeples. That can be used as a starting point for point evaluation later.
func getMeeplePositions(board map[Pos]Tile) map[Pos]bool {
	meeplePositions := map[Pos]bool{}
	for p, t := range board {
		if t.meeple.playerIndex != -1 {
			meeplePositions[p] = true
		}
	}
	return meeplePositions
}

// Returns any possibly random key from the map. Order doesn't matter here!
func getKey(positions map[Pos]bool) Pos {
	for k := range positions {
		return k
	}
	return Pos{0, 0}
}

// Calculates the immediate points, that are not yet finalized. So unfinished roads,
// unfinished cities or unfinished cloisters
func updateImmediatePoints(board map[Pos]Tile, playerScores *[]int) {

	// Initial set!
	meeplePositions := getMeeplePositions(board)

	for len(meeplePositions) > 0 {
		pos := getKey(meeplePositions)
		side := board[pos].meeple.sideIndex

		delete(meeplePositions, pos)

		searched := map[Pos]bool{}
		meeples := make([]int, len(*playerScores), len(*playerScores))
		score, positions, closed := calcRecursivePoints(board, pos, side, &searched, &meeples)

		// Closed structures should be handled by the updateFinalPoints() function. Not here, as it must handle
		// meeple removal as well!
		if !closed {
			bestPlayer := getBestPlayerIndex(meeples)
			for playerIndex, count := range meeples {
				if count == meeples[bestPlayer] {
					(*playerScores)[playerIndex] += score
				}
			}
		}

		// remove found meeple-positions from the initial set of meeplePositions! So we don't seach the same structure multiple times
		for _, p := range positions {
			delete(meeplePositions, p)
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

	for rounds := 0; rounds < 1; rounds++ {
		i := 0
		for i < len(tiles) {
			for _, player := range players {
				if i >= len(tiles) {
					break
				}
				tile := tiles[i]
				i += 1

				moves := generatePossibleMoves(board, []Tile{tile}, openPlacements, player)
				if len(moves) > 0 {
					move := moves[rand.Intn(len(moves))]
					placeTile(&board, &openPlacements, &players, move.tile, move.pos)
					updateFinalPoints(&board, move.pos, &players)
				}
			}
		}
	}

	drawField(board)

	for _, p := range players {
		fmt.Println(p)
	}

}
