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
	g_sides    []Pos = []Pos{Pos{-1, 0}, Pos{0, 1}, Pos{1, 0}, Pos{0, -1}}
	g_allSides []Pos = []Pos{Pos{-1, 0}, Pos{0, 1}, Pos{1, 0}, Pos{0, -1}, Pos{-1, -1}, Pos{1, -1}, Pos{1, -1}, Pos{1, 1}}
	// Even though those are not positions per se, it's two ints which is exactly what we need here.
	g_connectionIndices = []Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}}
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

type GameState struct {
	board          map[Pos]Tile
	tiles          []Tile
	players        []Player
	openPlacements map[Pos]bool
	lastMoves      []ReverseMove
}

type ReverseMeeplePlacement struct {
	playerIndex int
	pos         Pos
	side        int
}

type ReversePlayerPoints struct {
	playerIndex int
	points      int
}

// This struct is used to track changes to the whole gamestate when a player makes a move.
// That way we can reverse moves without having to copy the whole gamestate for branching
type ReverseMove struct {
	// A Meeple should be taken from the Player and placed (back) on the board
	playerToBoardMeeple []ReverseMeeplePlacement
	// A Meeple is taken from the board and put back into the player inventory
	boardToPlayerMeeple ReverseMeeplePlacement
	// Placed tile on the board that needs to be removed again. This one needs to
	// be added to openPlacements as well!
	// At the same time, remove all positions around this tile, that are not on the
	// board from the openPlacements set!
	// The tile also needs to be placed back onto the tiles-list of the game-state!
	removeTileFromBoard Pos
	// Final points that were awarded to a player
	awardedPoints []ReversePlayerPoints
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
	for _, c := range g_connectionIndices {
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

	if tile.meeple.sideIndex != SIDE_CENTER {
		tile.meeple.sideIndex = (tile.meeple.sideIndex + 1) % 4
	}
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

func placeTile(game *GameState, tile Tile, pos Pos, revMove *ReverseMove) {
	game.board[pos] = tile
	delete(game.openPlacements, pos)
	if tile.meeple.playerIndex != -1 {
		game.players[tile.meeple.playerIndex].meeples -= 1
		revMove.boardToPlayerMeeple = ReverseMeeplePlacement{tile.meeple.playerIndex, pos, tile.meeple.sideIndex}

	}

	revMove.removeTileFromBoard = pos

	for _, s := range g_sides {
		if _, ok := game.board[add(pos, s)]; !ok {
			game.openPlacements[add(pos, s)] = true
		}
	}

}

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

	for _, c := range g_connectionIndices {
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
// All players with the most meeples on the structure get full points
func getBestPlayerIndex(meeples []int) int {

	bestPlayer, bestCount := -1, 0
	for playerIndex, count := range meeples {
		if count > bestCount {
			bestCount = count
			bestPlayer = playerIndex
		}
	}
	return bestPlayer
}

// positions should only be tiles with a meeple on it, that needs to be removed!!!
func cleanupUsedMeeplesFromBoard(board *map[Pos]Tile, players *[]Player, positions []Pos) {
	// Clean up and remove meeples from the board. Add them back to the players inventory!
	for _, p := range positions {
		t := (*board)[p]
		(*players)[t.meeple.playerIndex].meeples += 1
		t.meeple = Meeple{-1, -1}
		(*board)[p] = t
	}
}

func countSurroundingTiles(board map[Pos]Tile, pos Pos) (count int) {
	for _, d := range g_allSides {
		if _, ok := board[add(pos, d)]; ok {
			count += 1
		}
	}
	return
}

// UpdateFinalPoints adds up the final points that are achieved by placing the tile at pos.
// It only counts finished cities and closed roads1
func updateFinalPoints(board *map[Pos]Tile, pos Pos, players *[]Player, revMove *ReverseMove) {

	tile := (*board)[pos]

	// Did we close all tiles around a cloister?
	for _, d := range g_allSides {
		tmpPos := add(pos, d)
		if t, ok := (*board)[tmpPos]; ok && t.cloister && t.meeple.playerIndex != -1 && t.meeple.sideIndex == SIDE_CENTER {
			if countSurroundingTiles(*board, tmpPos) == 8 {
				(*players)[t.meeple.playerIndex].score += 9
				(*players)[t.meeple.playerIndex].meeples += 1
				t.meeple = Meeple{-1, -1}
				(*board)[tmpPos] = t

				revMove.playerToBoardMeeple = append(revMove.playerToBoardMeeple, ReverseMeeplePlacement{t.meeple.playerIndex, tmpPos, SIDE_CENTER})
				revMove.awardedPoints = append(revMove.awardedPoints, ReversePlayerPoints{t.meeple.playerIndex, 9})
			}
		}
	}

	for side := 0; side < 4; side++ {
		if tile.hasConnectionAtSide(side) {

			// We always skip one side of a connection, so we only search each possible way once!
			ok := false
			for _, c := range g_connectionIndices {
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
					revMove.awardedPoints = append(revMove.awardedPoints, ReversePlayerPoints{playerIndex, score})
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
		tile := board[pos]
		side := tile.meeple.sideIndex

		delete(meeplePositions, pos)

		// Cloister tiles do not need to be calculated recursively. They can be short-cut
		if tile.cloister && tile.meeple.sideIndex == SIDE_CENTER {
			(*playerScores)[tile.meeple.playerIndex] += 1 + countSurroundingTiles(board, pos)
			continue
		}

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

func generateInitialBoard(playerCount int) GameState {
	startTile, tiles := getTiles()
	var players []Player
	for i := 0; i < playerCount; i++ {
		players = append(players, Player{i, 0, 6})
	}
	game := GameState{
		make(map[Pos]Tile),
		tiles,
		players,
		map[Pos]bool{Pos{-1, 0}: true, Pos{1, 0}: true, Pos{0, -1}: true, Pos{0, 1}: true},
		[]ReverseMove{},
	}

	game.board[Pos{0, 0}] = startTile

	return game
}

func main() {

	game := generateInitialBoard(3)

	for rounds := 0; rounds < 10; rounds++ {
		i := 0
		for i < len(game.tiles) {
			for _, player := range game.players {
				if i >= len(game.tiles) {
					break
				}
				tile := game.tiles[i]
				i += 1

				moves := generatePossibleMoves(game.board, []Tile{tile}, game.openPlacements, player)
				if len(moves) > 0 {
					move := moves[rand.Intn(len(moves))]

					revMove := ReverseMove{}
					placeTile(&game, move.tile, move.pos, &revMove)
					updateFinalPoints(&game.board, move.pos, &game.players, &revMove)
					game.lastMoves = append(game.lastMoves, revMove)
				}
			}
		}
	}

	drawField(game.board)

	for _, p := range game.players {
		fmt.Println(p)
	}

}
