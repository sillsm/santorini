package Santorini

import(
  "fmt"
  "strings"
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

func render(p Position) string {
	var b1 int32 = p.B1
	var b2 int32 = p.B2
	var b3 int32 = p.B3
	var b4 int32 = p.B4
	var a int32 = p.A
	var b int32 = p.B
	var x int32 = p.X
	var y int32 = p.Y

	var out string
	const tileSize = 6
	const tileNum = 5
	const gridSize = tileSize * tileNum

	out += "|"

	var piecePlacement [4]int
	for i := 0; i < 25; i++ {
		tile := "0"
		if b1%2 == 1 {
			tile = "1"
		}
		if b2%2 == 1 {
			tile = "2"
		}
		if b3%2 == 1 {
			tile = "3"
		}
		if b4%2 == 1 {
			tile = "4"
		}

		if a%2 == 1 {
			piecePlacement[0] = i
		}
		if b%2 == 1 {
			piecePlacement[1] = i
		}
		if x%2 == 1 {
			piecePlacement[2] = i
		}
		if y%2 == 1 {
			piecePlacement[3] = i
		}

		out += tile
		b1 = b1 >> 1
		b2 = b2 >> 1
		b3 = b3 >> 1
		b4 = b4 >> 1
		a = a >> 1
		b = b >> 1
		x = x >> 1
		y = y >> 1
		continue
	}

	pieces := fmt.Sprintf("%02d", piecePlacement)
	pieces = strings.Replace(pieces, "[", "", -1)
	pieces = strings.Replace(pieces, "]", "", -1)
	pieces = strings.Replace(pieces, " ", "", -1)
	out += "|" + pieces
	out += "|"

	return out

}

func (p Position) String() string {
	return render(p)
}

// Stringified positions take the form
// |0400300002001303040111124|08050018|
// 25 ints from 0 to 4 representing the heights
// of the 25 board squares.
//
// Then 4 ints (of format 00), representing position of
// White's first and second piece, then
// Black's first and second piece.
// Note that white and black's pieces are interchangeable.
// We should probably require the position go from 
// low to high as another integrity check
func NewPosition(s string) (Position, error) {
  // Position
  if len (s) != 36{
    panic("string integrity check fail, position incorrect length")
  }
  p := Position{}

  for i := 1; i < 26; i++{
    height := s[i]
    var mask int32
    mask = 1 << (i -1)
    if height == '1' {
      p.B1 |= mask
    }
    if height == '2' {
      p.B2 |= mask
    }
    if height == '3' {
      p.B3 |= mask
    }
    if height == '4' {
      p.B4 |= mask
    }
  }
  var err error
  whiteOne, err := strconv.Atoi(s[27:29])
  whiteTwo, err := strconv.Atoi(s[29:31])
  blackOne, err := strconv.Atoi(s[31:33])
  blackTwo, err := strconv.Atoi(s[33:35])
	
  if err != nil{
    panic("string integrity check fail, can't decode piece position")
  }
  p.A = occupancy[whiteOne]
  p.B = occupancy[whiteTwo]
  p.X = occupancy[blackOne]
  p.Y = occupancy[blackTwo]
  
  return p, nil
}
