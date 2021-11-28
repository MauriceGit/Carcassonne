package main

import (
	"testing"
)

func TestSmallClosedCityPoints(t *testing.T) {

	startTile, _ := getTiles()

	board := make(map[Pos]Tile)
	board[Pos{0, 0}] = startTile
	board[Pos{0, -1}] = Tile{11, [4]Area{AREA_GRASS, AREA_CITY, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{1, 2}}

	//drawField(board)

	expectedPoints := []int{0, 0, 4}
	players := []Player{Player{0, 0, 6}, Player{1, 0, 6}, Player{2, 0, 6}}
	revMove := ReverseMove{}
	updateFinalPoints(&board, Pos{0, 0}, &players, &revMove)

	for i, _ := range players {
		if players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, players[i].score, expectedPoints[i])
		}
	}

	//drawField(board)
}

func TestMediumClosedCityPoints(t *testing.T) {

	startTile, _ := getTiles()

	board := make(map[Pos]Tile)
	board[Pos{0, 0}] = startTile

	conns := connectionsToUint16([]Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}})
	board[Pos{0, -1}] = Tile{10, [4]Area{AREA_CITY, AREA_CITY, AREA_CITY, AREA_CITY}, false, true, conns, Meeple{3, 1}}
	board[Pos{0, -2}] = Tile{11, [4]Area{AREA_GRASS, AREA_CITY, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	conns = connectionsToUint16([]Pos{Pos{0, 2}})
	board[Pos{1, -1}] = Tile{12, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, true, conns, Meeple{-1, -1}}
	board[Pos{2, -1}] = Tile{13, [4]Area{AREA_CITY, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	board[Pos{-1, -1}] = Tile{14, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, false, 0x104, Meeple{-1, -1}}
	board[Pos{-2, -1}] = Tile{15, [4]Area{AREA_GRASS, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 18, 0}
	players := []Player{Player{0, 0, 6}, Player{1, 0, 6}, Player{2, 0, 6}}
	revMove := ReverseMove{}
	updateFinalPoints(&board, Pos{0, 0}, &players, &revMove)

	for i, _ := range players {
		if players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, players[i].score, expectedPoints[i])
		}
	}

	//drawField(board)
}

func TestMediumOpenCityPoints(t *testing.T) {

	startTile, _ := getTiles()

	board := make(map[Pos]Tile)
	board[Pos{0, 0}] = startTile

	conns := connectionsToUint16([]Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}})
	board[Pos{0, -1}] = Tile{10, [4]Area{AREA_CITY, AREA_CITY, AREA_CITY, AREA_CITY}, false, true, conns, Meeple{3, 1}}
	board[Pos{0, -2}] = Tile{11, [4]Area{AREA_GRASS, AREA_CITY, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	conns = connectionsToUint16([]Pos{Pos{0, 2}})
	board[Pos{1, -1}] = Tile{12, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, true, conns, Meeple{-1, -1}}
	board[Pos{2, -1}] = Tile{13, [4]Area{AREA_CITY, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	board[Pos{-1, -1}] = Tile{14, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, false, 0x104, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 0, 0}
	players := []Player{Player{0, 0, 6}, Player{1, 0, 6}, Player{2, 0, 6}}
	revMove := ReverseMove{}
	updateFinalPoints(&board, Pos{0, 0}, &players, &revMove)

	for i, _ := range players {
		if players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, players[i].score, expectedPoints[i])
		}
	}

	expectedPoints = []int{0, 8, 0}
	playerScores := []int{0, 0, 0}
	updateImmediatePoints(board, &playerScores)

	for i, score := range playerScores {
		if score != expectedPoints[i] {
			t.Errorf("Player %v has immediate point count. %v != %v (expected)", i, score, expectedPoints[i])
		}
	}

	//drawField(board)
}

func TestClosedCloisterPoints(t *testing.T) {

	board := make(map[Pos]Tile)
	board[Pos{0, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, true, false, 0, Meeple{SIDE_CENTER, 2}}

	board[Pos{-1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{-1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{-1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	board[Pos{1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	board[Pos{0, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{0, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 0, 9}
	players := []Player{Player{0, 0, 6}, Player{1, 0, 6}, Player{2, 0, 6}}
	revMove := ReverseMove{}
	updateFinalPoints(&board, Pos{-1, -1}, &players, &revMove)

	for i, _ := range players {
		if players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, players[i].score, expectedPoints[i])
		}
	}

}

func TestOpenCloisterPoints(t *testing.T) {

	board := make(map[Pos]Tile)
	board[Pos{0, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, true, false, 0, Meeple{SIDE_CENTER, 2}}

	board[Pos{-1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{-1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{-1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	board[Pos{1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	board[Pos{1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	board[Pos{0, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 0, 8}
	playerScores := []int{0, 0, 0}
	updateImmediatePoints(board, &playerScores)

	for i, score := range playerScores {
		if score != expectedPoints[i] {
			t.Errorf("Player %v has immediate point count. %v != %v (expected)", i, score, expectedPoints[i])
		}
	}

}
