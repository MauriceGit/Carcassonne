from enum import Enum
from recordtype import recordtype
from random import shuffle


class Area(Enum):
    Road  = 0
    Grass = 1
    City  = 2

# sides:       (left, down, right, up) -> (road, grass, city)
# cloister:    (True, False)
# connections: ((sides)] with sides -> (0..3] as indices to sides
#              connections are for roads, cities and grass areas
# meeple:      (side_index, player_index) with -1 for no_side and no_player
Tile = recordtype("Tile", "sides cloister emblem connections meeple")


# returns (start_tile, all_other_tiles)
def get_tiles(shuffled=False):
    tiles = []
    # WinA
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Grass, Area.Grass), True, False, (), (-1, -1)) for i in range(2)])
    # WinB
    tiles.extend([Tile((Area.Grass, Area.Grass, Area.Grass, Area.Grass), False, False, (), (-1, -1)) for i in range(4)])
    # WinC
    tiles.extend([Tile((Area.City, Area.City, Area.City, Area.City), False, True, ((0,1,2,3),), (-1, -1))])
    # WinD
    tiles.extend([Tile((Area.Road, Area.Grass, Area.Road, Area.City), False, False, ((0,2),), (-1, -1)) for i in range(3)])
    # WinE
    tiles.extend([Tile((Area.Grass, Area.Grass, Area.Grass, Area.City), False, False, (), (-1, -1)) for i in range(5)])
    # WinF
    tiles.extend([Tile((Area.City, Area.Grass, Area.City, Area.Grass), False, True, ((0,2),), (-1, -1)) for i in range(2)])
    # WinG
    tiles.extend([Tile((Area.City, Area.Grass, Area.City, Area.Grass), False, False, ((0,2),), (-1, -1))])
    # WinH
    tiles.extend([Tile((Area.Grass, Area.City, Area.Grass, Area.City), False, False, (), (-1, -1)) for i in range(3)])
    # WinI
    tiles.extend([Tile((Area.Grass, Area.Grass, Area.City, Area.City), False, False, (), (-1, -1)) for i in range(2)])
    # WinJ
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Road, Area.City), False, False, ((1,2),), (-1, -1)) for i in range(3)])
    # WinK
    tiles.extend([Tile((Area.Road, Area.Road, Area.Grass, Area.City), False, False, ((0,1),), (-1, -1)) for i in range(3)])
    # WinL
    tiles.extend([Tile((Area.Road, Area.Road, Area.Road, Area.City), False, False, (), (-1, -1)) for i in range(3)])
    # WinM
    tiles.extend([Tile((Area.City, Area.Grass, Area.Grass, Area.City), False, True, ((0,3),), (-1, -1)) for i in range(2)])
    # WinN
    tiles.extend([Tile((Area.City, Area.Grass, Area.Grass, Area.City), False, False, ((0,3),), (-1, -1)) for i in range(3)])
    # WinO
    tiles.extend([Tile((Area.City, Area.Road, Area.Road, Area.City), False, True, ((0,3), (1,2)), (-1, -1)) for i in range(2)])
    # WinP
    tiles.extend([Tile((Area.City, Area.Road, Area.Road, Area.City), False, False, ((0,3), (1,2)), (-1, -1)) for i in range(3)])
    # WinQ
    tiles.extend([Tile((Area.City, Area.Grass, Area.Grass, Area.City), False, True, ((0,2,3),), (-1, -1))])
    # WinR
    tiles.extend([Tile((Area.City, Area.Grass, Area.Grass, Area.City), False, False, ((0,2,3),), (-1, -1)) for i in range(3)])
    # WinS
    tiles.extend([Tile((Area.City, Area.Road, Area.Road, Area.City), False, True, ((0,2,3),), (-1, -1)) for i in range(2)])
    # WinT
    tiles.extend([Tile((Area.City, Area.Road, Area.Road, Area.City), False, False, ((0,2,3),), (-1, -1))])
    # WinU
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Grass, Area.Road), False, False, ((1,3),), (-1, -1)) for i in range(8)])
    # WinV
    tiles.extend([Tile((Area.Road, Area.Road, Area.Grass, Area.Grass), False, False, ((0,1),), (-1, -1)) for i in range(9)])
    # WinW
    tiles.extend([Tile((Area.Road, Area.Road, Area.Road, Area.Grass), False, False, (), (-1, -1)) for i in range(4)])
    # WinX
    tiles.extend([Tile((Area.Road, Area.Road, Area.Road, Area.Road), False, False, (), (-1, -1))])
    # WinTierA
    tiles.extend([Tile((Area.Road, Area.Road, Area.Road, Area.Grass), False, False, (), (-1, -1))])
    # WinTierB
    tiles.extend([Tile((Area.Road, Area.Grass, Area.Grass, Area.City), False, False, (), (-1, -1))])
    # WinTierC
    tiles.extend([Tile((Area.Grass, Area.Grass, Area.Road, Area.City), False, False, (), (-1, -1))])
    # WinTierD
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Grass, Area.Grass), True, False, (), (-1, -1))])
    # WinTierE
    tiles.extend([Tile((Area.Road, Area.Grass, Area.City, Area.City), False, False, ((2,3),), (-1, -1))])
    # WinTierF
    tiles.extend([Tile((Area.Road, Area.Grass, Area.Road, Area.Grass), True, False, (), (-1, -1))])
    # WinTierG
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Road, Area.City), False, False, ((1,2),), (-1, -1))])
    # WinTierH
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Grass, Area.City), False, False, (), (-1, -1))])
    # WinTierI
    tiles.extend([Tile((Area.Road, Area.Road, Area.Road, Area.Road), False, False, ((0,3), (1,2)), (-1, -1))])
    # WinTierJ
    tiles.extend([Tile((Area.Grass, Area.Road, Area.City, Area.City), False, False, ((2,3),), (-1, -1))])
    # WinTierK
    tiles.extend([Tile((Area.Grass, Area.Road, Area.Grass, Area.Grass), False, False, (), (-1, -1))])
    # WinTierL
    tiles.extend([Tile((Area.Road, Area.Road, Area.Grass, Area.City), False, False, ((0,1),), (-1, -1))])

    if shuffled:
        shuffle(tiles)

    start_tile = Tile((Area.Road, Area.Grass, Area.Road, Area.City), False, False, ((0,2),), (-1, -1))

    return start_tile, tiles
