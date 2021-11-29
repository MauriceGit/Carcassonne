package main

import (
	"testing"
)

func TestSmallClosedCityPoints(t *testing.T) {

	game := generateInitialBoard(3)
	game.board[Pos{0, -1}] = Tile{11, [4]Area{AREA_GRASS, AREA_CITY, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{1, 2}}

	//drawField(game.board)

	expectedPoints := []int{0, 0, 4}
	revMove := ReverseMove{}
	game.updateFinalPoints(Pos{0, 0}, &revMove)

	for i, _ := range game.players {
		if game.players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, game.players[i].score, expectedPoints[i])
		}
	}

	//drawField(board)
}

func TestMediumClosedCityPoints(t *testing.T) {

	game := generateInitialBoard(3)

	conns := connectionsToUint16([]Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}})
	game.board[Pos{0, -1}] = Tile{10, [4]Area{AREA_CITY, AREA_CITY, AREA_CITY, AREA_CITY}, false, true, conns, Meeple{3, 1}}
	game.board[Pos{0, -2}] = Tile{11, [4]Area{AREA_GRASS, AREA_CITY, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	conns = connectionsToUint16([]Pos{Pos{0, 2}})
	game.board[Pos{1, -1}] = Tile{12, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, true, conns, Meeple{-1, -1}}
	game.board[Pos{2, -1}] = Tile{13, [4]Area{AREA_CITY, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	game.board[Pos{-1, -1}] = Tile{14, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, false, 0x104, Meeple{-1, -1}}
	game.board[Pos{-2, -1}] = Tile{15, [4]Area{AREA_GRASS, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 18, 0}
	revMove := ReverseMove{}
	game.updateFinalPoints(Pos{0, 0}, &revMove)

	for i, _ := range game.players {
		if game.players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, game.players[i].score, expectedPoints[i])
		}
	}

	//drawField(board)
}

func TestMediumOpenCityPoints(t *testing.T) {

	game := generateInitialBoard(3)

	conns := connectionsToUint16([]Pos{Pos{0, 1}, Pos{0, 2}, Pos{0, 3}, Pos{1, 2}, Pos{1, 3}, Pos{2, 3}})
	game.board[Pos{0, -1}] = Tile{10, [4]Area{AREA_CITY, AREA_CITY, AREA_CITY, AREA_CITY}, false, true, conns, Meeple{3, 1}}
	game.board[Pos{0, -2}] = Tile{11, [4]Area{AREA_GRASS, AREA_CITY, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	conns = connectionsToUint16([]Pos{Pos{0, 2}})
	game.board[Pos{1, -1}] = Tile{12, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, true, conns, Meeple{-1, -1}}
	game.board[Pos{2, -1}] = Tile{13, [4]Area{AREA_CITY, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0x0, Meeple{-1, -1}}
	game.board[Pos{-1, -1}] = Tile{14, [4]Area{AREA_CITY, AREA_GRASS, AREA_CITY, AREA_GRASS}, false, false, 0x104, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 0, 0}
	revMove := ReverseMove{}
	game.updateFinalPoints(Pos{0, 0}, &revMove)

	for i, _ := range game.players {
		if game.players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, game.players[i].score, expectedPoints[i])
		}
	}

	expectedPoints = []int{0, 8, 0}
	playerScores := []int{0, 0, 0}
	game.updateImmediatePoints(&playerScores)

	for i, score := range playerScores {
		if score != expectedPoints[i] {
			t.Errorf("Player %v has immediate point count. %v != %v (expected)", i, score, expectedPoints[i])
		}
	}

	//drawField(board)
}

func TestClosedCloisterPoints(t *testing.T) {
	game := generateInitialBoard(3)

	game.board[Pos{0, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, true, false, 0, Meeple{SIDE_CENTER, 2}}

	game.board[Pos{-1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{-1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{-1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	game.board[Pos{1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	game.board[Pos{0, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{0, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	//drawField(game.board)

	expectedPoints := []int{0, 0, 9}
	revMove := ReverseMove{}
	game.updateFinalPoints(Pos{-1, -1}, &revMove)

	for i, _ := range game.players {
		if game.players[i].score != expectedPoints[i] {
			t.Errorf("Player %v has wrong point count. %v != %v (expected)", i, game.players[i].score, expectedPoints[i])
		}
	}

}

func TestOpenCloisterPoints(t *testing.T) {
	game := generateInitialBoard(3)

	game.board[Pos{0, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, true, false, 0, Meeple{SIDE_CENTER, 2}}

	game.board[Pos{-1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{-1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{-1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	game.board[Pos{1, 0}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{1, -1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}
	game.board[Pos{1, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	game.board[Pos{0, 1}] = Tile{0, [4]Area{AREA_GRASS, AREA_GRASS, AREA_GRASS, AREA_GRASS}, false, false, 0, Meeple{-1, -1}}

	//drawField(board)

	expectedPoints := []int{0, 0, 8}
	playerScores := []int{0, 0, 0}
	game.updateImmediatePoints(&playerScores)

	for i, score := range playerScores {
		if score != expectedPoints[i] {
			t.Errorf("Player %v has immediate point count. %v != %v (expected)", i, score, expectedPoints[i])
		}
	}

}
