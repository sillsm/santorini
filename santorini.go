package main

import (
	"fmt"
	"strconv"
)

type Position struct {
	B1 int32
	B2 int32
	B3 int32
	B4 int32
	A  int32 // First White Piece
	B  int32 // Second White Piece
	X  int32 // First Red Piece
	Y  int32 // Second Red Piece
}

var (
	// Column and Row Mask Constants
	c1 int32 = 0x00108421
	c2 int32 = c1 << 1
	c3 int32 = c1 << 2
	c4 int32 = c1 << 3
	c5 int32 = c1 << 4

	r1 int32 = 0x0000001F
	r2 int32 = r1 << 5
	r3 int32 = r1 << 10
	r4 int32 = r1 << 15
	r5 int32 = r1 << 20

	// King moves for center piece.
	center int32 = (r2 | r3 | r4) & (c2 | c3 | c4) &^ (c3 & r3)
)

// King move masks for each of the 25 locations.
// 0 to 24
var kingMasks []int32 = []int32{
	center >> 12 &^ c5,
	center >> 11,
	center >> 10,
	center >> 9,
	center >> 8 &^ c1,
	center >> 7 &^ c5,
	center >> 6,
	center >> 5,
	center >> 4,
	center >> 3 &^ c1,
	center >> 2 &^ c5,
	center >> 1,
	center, // 12
	center << 1,
	center << 2 &^ c1,
	center << 3 &^ c5,
	center << 4,
	center << 5,
	center << 6, //18
	center << 7 &^ c1,
	center << 8 &^ c5,
	center << 9,
	center << 10,
	center << 11,
	center << 12 &^ c1,
}

// Single bit piece occupancy masks.
var occupancy []int32 = []int32{
	1,
	1 << 1,
	1 << 2,
	1 << 3,
	1 << 4,
	1 << 5,
	1 << 6,
	1 << 7,
	1 << 8,
	1 << 9,
	1 << 10,
	1 << 11,
	1 << 12,
	1 << 13,
	1 << 14,
	1 << 15,
	1 << 16,
	1 << 17,
	1 << 18,
	1 << 19,
	1 << 20,
	1 << 21,
	1 << 22,
	1 << 23,
	1 << 24,
}

const tileSize = 6
const tileNum = 5
const gridSize = tileSize * tileNum

type moveBuild struct {
	move  int
	build int
	piece int32
}

// Renders square blocks to a bigger square.
func render(p Position) {
	var b1 int32 = p.B1
	var b2 int32 = p.B2
	var b3 int32 = p.B3
	var b4 int32 = p.B4
	var a int32 = p.A
	var b int32 = p.B
	var x int32 = p.X
	var y int32 = p.Y

	tile := [tileSize][tileSize]rune{
		{' ', '^', '^', '^', '^', ' '},
		{' ', '|', '1', ' ', '|', ' '}, // $ magic sub tile.
		{' ', '|', ' ', '4', '|', ' '},
		{' ', '5', '6', ' ', ' ', ' '},
		{' ', '^', '^', '^', '^', ' '},
		{' ', ' ', ' ', ' ', ' ', ' '},
	}

	// Tile the output.
	output := [gridSize][gridSize + 1]rune{}
	t := 0
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize+1; j++ {
			output[i][j] = tile[i%tileSize][j%tileSize]
			// Paint height of tiles.
			if output[i][j] == '1' {
				output[i][j] = '0'
				if b1%2 == 1 {
					output[i][j] = '1'
				}
				if b2%2 == 1 {
					output[i][j] = '2'
				}
				if b3%2 == 1 {
					output[i][j] = '3'
				}
				if b4%2 == 1 {
					output[i][j] = '4'
				}
				b1 = b1 >> 1
				b2 = b2 >> 1
				b3 = b3 >> 1
				b4 = b4 >> 1
				continue
			}
			// Paint piece names
			if output[i][j] == '4' {
				output[i][j] = ' '
				if a%2 == 1 {
					output[i][j] = 'A'
				}
				if b%2 == 1 {
					output[i][j] = 'B'
				}
				if x%2 == 1 {
					output[i][j] = 'X'
				}
				if y%2 == 1 {
					output[i][j] = 'Y'
				}
				a = a >> 1
				b = b >> 1
				x = x >> 1
				y = y >> 1
			}

			// Number Tiles
			if output[i][j] == '6' {
				s := strconv.Itoa(t)
				t = t + 1
				if len(s) == 1 {
					output[i][j] = rune(s[0])
					output[i][j-1] = ' '

				} else {
					output[i][j-1] = rune(s[0])
					output[i][j] = rune(s[1])
				}
			}

			if j == gridSize {
				output[i][j] = '\n'
			}
		}
	}

	// Print the output
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize+1; j++ {
			fmt.Printf(string(output[i][j]))
		}
	}

}

func clearScreen() {
	//print("\033[H\033[2J")
}

func extractBits(i int32) []int {
	ret := []int{}
	for j := 0; j < 25; j++ {
		if i%2 == 1 {
			ret = append(ret, j)
		}
		i = i >> 1
	}
	return ret
}

// Given position and piece, returns legal moves from that spot.
func legalMoves(p Position, piece int32) []int {
	i := extractBits(piece)
	if len(i) != 1 {
		panic("legalMoves assumes piece has single bit set.")
	}
	moves := kingMasks[i[0]] &^ (p.A | p.B | p.X | p.Y | p.B4)
	// If you're not at least 1 high, can't go to 2 or 3
	if piece&p.B1 == 0 {
		moves = moves &^ (p.B2 | p.B3)
	}
	// If you're not at least 2 high, can't go to 3
	if piece&p.B2 == 0 {
		moves = moves &^ (p.B3)
	}
	return extractBits(moves)
}

// Given position and piece, returns locations that can be built from spot.
func legalBuilds(p Position, piece int32) []int {
	i := extractBits(piece)
	if len(i) != 1 {
		panic("legalBuilds assumes piece has single bit set.")
	}
	// Get kingmoves - occupied spots
	builds := kingMasks[i[0]] &^ (p.A | p.B | p.X | p.Y | p.B4)
	return extractBits(builds)
}

func legalMoveBuilds(p Position, piece int32) []moveBuild {
	var ret []moveBuild
	moves := legalMoves(p, piece)

	for _, val := range moves {
		testPosition := p
		if p.A == piece {
			testPosition.A = occupancy[val]
		}
		if p.B == piece {
			testPosition.B = occupancy[val]
		}
		if p.X == piece {
			testPosition.X = occupancy[val]
		}
		if p.Y == piece {
			testPosition.Y = occupancy[val]
		}

		builds := legalBuilds(testPosition, occupancy[val])
		for _, build := range builds {
			ret = append(ret, moveBuild{val, build, piece})
		}

	}
	return ret
}

func updatePosition(p Position, m moveBuild) Position {
	location := occupancy[m.move]
	buildBit := occupancy[m.build]
	// logic to move the piece
	switch {
	case p.A&m.piece > 0:
		p.A = location
	case p.B&m.piece > 0:
		p.B = location
	case p.X&m.piece > 0:
		p.X = location
	case p.Y&m.piece > 0:
		p.Y = location
	}
	// logic to build
	switch {
	case p.B1&buildBit == 0:
		p.B1 |= buildBit
	case p.B2&buildBit == 0:
		p.B2 |= buildBit
	case p.B3&buildBit == 0:
		p.B3 |= buildBit
	case p.B4&buildBit == 0:
		p.B4 |= buildBit
	}
	return p
}

func printABMoves(position Position) Position {
	movesA := legalMoveBuilds(position, position.A)
	movesB := legalMoveBuilds(position, position.B)

	// Test
	fmt.Printf("\nMove A:\n")
	for i, move := range movesA {
		fmt.Printf("%v: %v|%v\t", i+100, move.move, move.build)
		if i%5 == 0 && i > 0 {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("\nMove B:\n")
	for i, move := range movesB {
		fmt.Printf("%v: %v|%v\t", i+200, move.move, move.build)
		if i%5 == 0 && i > 0 {
			fmt.Printf("\n")
		}
	}

	var input string
	var i int
	for {
		fmt.Scanln(&input)
		var err error
		i, err = strconv.Atoi(input)
		if err != nil {
            panic(err)
		}
		break
	}

	if i >= 200 {
		i -= 200
		return updatePosition(position, movesB[i])

	} else {
		i -= 100
		return updatePosition(position, movesA[i])
	}
    fmt.Printf("Returned from print moves")
	return position
}

func main() {
	startPosition := Position{
		0x1FAB212, 0x182A212, 0x102A012, 0x1020002,
		occupancy[11],
		occupancy[12],
		occupancy[16],
		occupancy[17],
	}

	// repl
	var position Position = startPosition
	//var init bool
	for {
		clearScreen()
		render(position)
		// Move piece A logic
		position = printABMoves(position)
	}
}
