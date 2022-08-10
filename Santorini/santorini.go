package Santorini

import (
	"fmt"
  "sort"
	"strconv"
	"strings"
)

type Position struct {
	B1             int32
	B2             int32
	B3             int32
	B4             int32
	A              int32 // First White Piece
	B              int32 // Second White Piece
	X              int32 // First Red Piece
	Y              int32 // Second Red Piece
	Representation string
	Ply            bool // False for White, which moves first, True for Black.
}

type MoveBuild struct {
	Move  int32
	Build int32
	Ply   bool // Whose turn is it
	Piece bool // Does their first or second piece move
}

// Interface for the game tree search, so we can make mocks against it.
type GameNode interface {
  Children() []GameNode
  String() string
  // 'W', 'B', or '?'
  Outcome() rune
  WhichPly() bool //who's turn is it
  Score()(int,int) //position score, for the player whose turn it is
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

// map where keys are positionally encoded pieces
// like 0b0000001000000000 for a piece in square 6
// and values are list of positially encoded possible
// king moves if it were a chessboard.
var kingMoves = map[int32][]int32{
	occupancy[0]:  []int32{1 << 1, 1 << 5, 1 << 6},
	occupancy[1]:  []int32{1 << 0, 1 << 2, 1 << 5, 1 << 6, 1 << 7},
	occupancy[2]:  []int32{1 << 1, 1 << 3, 1 << 6, 1 << 7, 1 << 8},
	occupancy[3]:  []int32{1 << 2, 1 << 4, 1 << 7, 1 << 8, 1 << 9},
	occupancy[4]:  []int32{1 << 3, 1 << 8, 1 << 9},
	occupancy[5]:  []int32{1 << 0, 1 << 1, 1 << 6, 1 << 10, 1 << 11},
	occupancy[6]:  []int32{1 << 0, 1 << 1, 1 << 2, 1 << 5, 1 << 7, 1 << 10, 1 << 11, 1 << 12},
	occupancy[7]:  []int32{1 << 1, 1 << 2, 1 << 3, 1 << 6, 1 << 8, 1 << 11, 1 << 12, 1 << 13},
	occupancy[8]:  []int32{1 << 2, 1 << 3, 1 << 4, 1 << 7, 1 << 9, 1 << 12, 1 << 13, 1 << 14},
	occupancy[9]:  []int32{1 << 3, 1 << 4, 1 << 8, 1 << 13, 1 << 14},
	occupancy[10]: []int32{1 << 5, 1 << 6, 1 << 11, 1 << 15, 1 << 16},
	occupancy[11]: []int32{1 << 5, 1 << 6, 1 << 7, 1 << 10, 1 << 12, 1 << 15, 1 << 16, 1 << 17},
	occupancy[12]: []int32{1 << 6, 1 << 7, 1 << 8, 1 << 11, 1 << 13, 1 << 16, 1 << 17, 1 << 18},
	occupancy[13]: []int32{1 << 7, 1 << 8, 1 << 9, 1 << 12, 1 << 14, 1 << 17, 1 << 18, 1 << 19},
	occupancy[14]: []int32{1 << 8, 1 << 9, 1 << 13, 1 << 18, 1 << 19},
	occupancy[15]: []int32{1 << 10, 1 << 11, 1 << 16, 1 << 20, 1 << 21},
	occupancy[16]: []int32{1 << 10, 1 << 11, 1 << 12, 1 << 15, 1 << 17, 1 << 20, 1 << 21, 1 << 22},
	occupancy[17]: []int32{1 << 11, 1 << 12, 1 << 13, 1 << 16, 1 << 18, 1 << 21, 1 << 22, 1 << 23},
	occupancy[18]: []int32{1 << 12, 1 << 13, 1 << 14, 1 << 17, 1 << 19, 1 << 22, 1 << 23, 1 << 24},
	occupancy[19]: []int32{1 << 13, 1 << 14, 1 << 18, 1 << 23, 1 << 24},
	occupancy[20]: []int32{1 << 15, 1 << 16, 1 << 21},
	occupancy[21]: []int32{1 << 15, 1 << 16, 1 << 17, 1 << 20, 1 << 22},
	occupancy[22]: []int32{1 << 16, 1 << 17, 1 << 18, 1 << 21, 1 << 23},
	occupancy[23]: []int32{1 << 17, 1 << 18, 1 << 19, 1 << 22, 1 << 24},
	occupancy[24]: []int32{1 << 18, 1 << 19, 1 << 23},
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
	//if p.Representation == "" {
	return render(p)
	// }
	// return p.Representation
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
	if len(s) != 36 {
		panic("string integrity check fail, position incorrect length")
	}
	p := Position{}

	parity := false // asume it's white's turn

	for i := 1; i < 26; i++ {
		height := s[i]
		var mask int32
		mask = 1 << (i - 1)
		switch height {
		case '1':
			p.B1 |= mask
			parity = !parity // flip parity bit for odd tiles
		case '2':
			p.B1 |= mask
			p.B2 |= mask
		case '3':
			p.B1 |= mask
			p.B2 |= mask
			p.B3 |= mask
			parity = !parity // flip parity bit for odd tiles
		case '4':
			p.B1 |= mask
			p.B2 |= mask
			p.B3 |= mask
			p.B4 |= mask
		}

	}
	var err error
	whiteOne, err := strconv.Atoi(s[27:29])
	whiteTwo, err := strconv.Atoi(s[29:31])
	blackOne, err := strconv.Atoi(s[31:33])
	blackTwo, err := strconv.Atoi(s[33:35])

	// swap occupancies to keep the lowest numbered piece first
	if whiteOne > whiteTwo {
		whiteOne, whiteTwo = whiteTwo, whiteOne
	}
	if blackOne > blackTwo {
		blackOne, blackTwo = blackTwo, blackOne
	}

	if err != nil {
		panic("string integrity check fail, can't decode piece position")
	}
	p.A = occupancy[whiteOne]
	p.B = occupancy[whiteTwo]
	p.X = occupancy[blackOne]
	p.Y = occupancy[blackTwo]
	p.Ply = parity

	return p, nil
}

func UpdatePosition(p Position, b MoveBuild) Position {
	// update strings here too
	switch {
	case !b.Ply && !b.Piece:
		p.A = b.Move
	case !b.Ply && b.Piece:
		p.B = b.Move
	case b.Ply && !b.Piece:
		p.X = b.Move
	case b.Ply && b.Piece:
		p.Y = b.Move
	}
	// logic to build
	switch {
	case p.B1&b.Build == 0:
		p.B1 |= b.Build
	case p.B2&b.Build == 0:
		p.B2 |= b.Build
	case p.B3&b.Build == 0:
		p.B3 |= b.Build
	case p.B4&b.Build == 0:
		p.B4 |= b.Build
	}

	// make sure
	if p.A > p.B {
		p.A, p.B = p.B, p.A
	}
	if p.X > p.Y {
		p.X, p.Y = p.Y, p.X
	}
  //Flip whose turn it is
  p.Ply = !p.Ply
	return p
}

func legalMoves2(p Position, piece int32) []int32 {
	mask := ^(p.A | p.B | p.X | p.Y | p.B4)
	// If you're not at least 1 high, can't go to 2 or 3
	if piece&p.B1 == 0 {
		mask = mask &^ (p.B2 | p.B3)
	}
	// If you're not at least 2 high, can't go to 3
	if piece&p.B2 == 0 {
		mask = mask &^ (p.B3)
	}
	//fmt.Printf("MASK %v\n", mask)
	var ret []int32

	km, ok := kingMoves[piece]
	if !ok {
		panic("Map error")
	}

	// this will be at most 8 ops, often fewer.
	for _, move := range km {
		if test := move & mask; test != 0 {
			ret = append(ret, test)
		}
	}
	return ret
}

// legalBuilds assumes the piece don't move.
// so applying to it any but one of the 4 actual pieces
// returns nonsense.
func legalBuilds(p Position, piece int32) []int32 {
	// Assume you can build anywhere
	// minus where any of the pieces are, or on a 4 tile.
	//var mask int32 = 1<<25
	mask := ((1 << 25) - 1) &^ (p.A | p.B | p.X | p.Y | p.B4)
	//fmt.Printf("At Position %v, got mask \n%v\n", p.String(), mask)

	var ret []int32

	km, ok := kingMoves[piece]
	if !ok {
		panic("Map error")
	}

	// this will be at most 8 ops, often fewer.
	for _, move := range km {
		if test := move & mask; test != 0 {
			ret = append(ret, test)
		}
	}
	return ret
}

func legalBuildMoves(p Position) []MoveBuild {
	var ret []MoveBuild
	// If it's white's turn to move
	var piece1, piece2 int32
	if !p.Ply {
		piece1 = p.A
		piece2 = p.B
	} else {
		piece1 = p.X
		piece2 = p.Y
	}

	// Don't forget to update the position before you try to get the builds.
	for _, m := range legalMoves2(p, piece1) {
		testPiece := p //here it is
		if !p.Ply {
			testPiece.A = m
		} else {
			testPiece.X = m
		}
    legalbuilds := legalBuilds(testPiece, m)
    // Don't consider a move if there are no builds from it.
    if len(legalbuilds) == 0 {
      continue
    }
		for _, b := range legalBuilds(testPiece, m) {
			ret = append(ret, MoveBuild{m, b, p.Ply, false})
		}
	}
	for _, m := range legalMoves2(p, piece2) {
		testPiece := p
		if !p.Ply {
			testPiece.B = m
		} else {
			testPiece.Y = m
		}
    // Don't consider a move if there are no builds from it.
    legalbuilds := legalBuilds(testPiece, m)
    if len(legalbuilds) == 0 {
      continue
    }
		for _, b := range legalBuilds(testPiece, m) {
			ret = append(ret, MoveBuild{m, b, p.Ply, true})
		}
	}
	return ret
}

/*
 Make sure the position implements the GameNode interface
*/

func (p Position) Children()[]GameNode{
  var pp []GameNode
  for _, mb := range legalBuildMoves(p) {
    updated := UpdatePosition(p, mb)
    // If any of the moves are a win, that's it. Return no children.
    if (mb.Move & updated.B3) > 0 {
      return nil
    }
    pp = append(pp, updated)
  }
  return pp
}

func (p Position) Outcome() rune {
  // if it's my turn, one of my pieces is on a 2, and can move to a 3, I win.
  // Or, if either side can either not build or not move, I win.
  // If I can't build or move, I lose.
  lbm := legalBuildMoves(p)
  if len(lbm) == 0 {
    //panic("No legal moves")
    if !p.Ply {
      return 'B'
    } else {
      return 'W'
    }
  }

  for _, move := range lbm{
    if (move.Move & p.B3) > 0 {
      if !p.Ply {
        return 'W'
      } else {
        return 'B'
      }
    }
  }
  return '?'
}

func (p Position) WhichPly() bool{
  return p.Ply
}

// Give a score for White and Black
func (p Position)Score() (int, int){
  // See out the first 1000 outcomes from here.
  m, depth := ShallowExploreNode(p, 100)
  w := 0
  b := 0
  //fmt.Printf("The score sheet I saw %v\n", m)
  for _, v := range m{
    if v == 'W'{
      w++
    }
    if v == 'B'{
      b++
    }
  }
  // NOTE! Usually a ply is bool 0 for White to move and 1 for Black to move.
  // However, we usually check children of a node for their scores, and in
  // that case, they should have the same ply as their parent node, until
  // they've *actually been chosen* as a path to go down.
  if p.Ply {
    return w-b, depth
  }
  return b-w, depth
}

/*
Shallow ExploreNode is a BFS.
Find the first leafLimit bfs outcomes.
*/
// 21020 00100 44344 00231 00010 |01081924| false ?
func ShallowExploreNode(gn GameNode, leafLimit int) (map[string]rune, int){
  depth := 0
  shallowestDepth := 100
  currentGeneration := 1
  nextGeneration := 0

  exploreLimit := 10000
  explored := 0
  m := make(map[string]rune)
  var toExplore []GameNode
  toExplore = append(toExplore, gn)

  for (len(toExplore)>0 ){
     var pop GameNode
     explored++
     pop, toExplore = toExplore[0], toExplore[1:]

     // Now, get the deets on the head node.
     if o := pop.Outcome(); (o == 'W') || (o == 'B'){
      //fmt.Printf("\nHit a leaf %v, with outcome %v", pop.String(), string(o))
     //fmt.Printf("\n%v", n.String()+string(o))
       m[pop.String()]=o
       if depth < shallowestDepth{
         shallowestDepth = depth
       }

     }

     // See if it's time to bail
     if len(m) > leafLimit{
       break
     }
     if explored > exploreLimit{
       break
     }

     // Otherwise, get the children and add them to the list to inspect
     for _, child := range pop.Children(){
       toExplore = append(toExplore, child)
       nextGeneration++
     }

     currentGeneration = currentGeneration - 1
     if (currentGeneration == 0){
        depth++
        currentGeneration = nextGeneration
      }
  }

  // After search stuff
  return m, shallowestDepth
}

/*
Code to explore the game tree.
*/
func ExploreNode(gn GameNode, leafLimit int) map[string]rune{
  m := make(map[string]rune)
  // Don't forget, this is PLY depth. so White and black moves combined
  //  maxDepth := 6

  limit := 0

  var f func(n GameNode, d int)rune
  f = func(n GameNode, d int)rune{
     if limit % 100000 == 0{
       //fmt.Printf("%v\tsolved:%v\tshortest:%v\n",limit, len(m), shortest)
       fmt.Printf("Limit: %v Depth %v \t State:%v|%v\n", limit, d, n.String(),string(n.Outcome()))
     }
       // Leaf Limit
       if leafLimit > 0 {
         if len(m) > leafLimit {
           return '?'
         }
       }


     // Did I see this state before? If so, stop exploring it and descendants and
     // return what I know about it.
     if val,ok := m[n.String()]; ok{
       return val
     }
     limit++


    // Check my outcome and return it if I'm a leaf node
    if o := n.Outcome(); (o == 'W') || (o == 'B'){
      //fmt.Printf("\nHit a leaf %v, with outcome %v", n.String(), string(o))
    //fmt.Printf("\n%v", n.String()+string(o))
      m[n.String()]=o
      return o
    }
    // So I'm not a leaf node.
    childOutcomes := make(map[rune]int)

    children := n.Children()
    if leafLimit == 0 {

        //fmt.Printf("About to consider my children, %v\n", len(n.Children()))
        sort.Slice(children, func(p, q int) bool {
          pScore, pmindepth := children[p].Score()
          qScore, qmindepth := children[q].Score()
          if (pScore == qScore){
              return pmindepth < qmindepth
          }
          return pScore > qScore})
        //fmt.Printf("------------------------\n")
        //for _, l := range children{
          //l=l
          //fmt.Printf("Child %v, score %v, %v\n" , l.String(), l.Score(), l.WhichPly())
        }


    //fmt.Printf("Candidate moves and score:")
  //  for _, c := range children{
  //    fmt.Printf("\n%v %v", c, c.Score())

    //}
    //fmt.Printf("\n")


    for i, c := range children{
      //if (i == 0) && (leafLimit ==0){
      //  fmt.Printf("I CHOOSE %v | %v | %v | %v \n", c.String(), c.WhichPly(), len(m), string(c.Outcome()))
      //}
      // Best-only choice
      if i > 0 && (leafLimit == 0){
        continue
      }
      fmt.Printf("Going to: %v %v %v\n", c.String(), c.WhichPly(), string(c.Outcome()))
      // Recurse here
      outcome := f(c, d+1)
      //fmt.Printf("State:%v\n", c.String())
    //  fmt.Printf("\nOutcome came back as %v, ply:%v", string(outcome), !n.WhichPly())
      //fmt.Printf("\nSanity check %v %v, ", !n.WhichPly(), (outcome == 'B'))
      childOutcomes[outcome] += 1


      // If this child's outcome is one I would certainly
      // pick because it's my turn and I can to win, assume I will.
      // Stop searching
      if (!n.WhichPly() && (outcome == 'W')) || (n.WhichPly() && (outcome == 'B')){
        m[n.String()]=outcome
        return outcome
      }

      //outcome := c.Outcome()
      //rep     := c.String()
      //fmt.Printf("%v%v|%v|%v\n", rep, string(outcome), c.Score(), c.WhichPly())
      // Did I see this state before? If so, stop exploring it and descendants.
      //if _,ok := m[c.String()]; ok{
      //  continue
      //}

      // If this child's outcome is one I would certainly
      // pick because of parity, make it my outcome,
      // and stop exploring.

      // If this child's outcome is one that I would certainly
      // not pick, don't explore it further.
      //if (n.WhichPly() && (outcome == 'W')) || (!n.WhichPly() && (outcome == 'B')){
      //  m[n.String()]=outcome
      //  continue
      //}
      /*

      if (!n.WhichPly() && (outcome == 'W')) || (n.WhichPly() && (outcome == 'B')){
      if (!n.WhichPly() && (outcome == 'W')) {
        panic("FUCK YOU")
        m[n.String()]=outcome
        return outcome
      }

      // If this child's outcome is one that I would certainly
      // not pick, don't explore it further.
      if (n.WhichPly() && (outcome == 'W')) || (!n.WhichPly() && (outcome == 'B')){
        m[n.String()]=outcome
        continue
      }*/
    }
    //fmt.Printf("\n\n\n\n\n\nChild outcomes for this level %v\n", n.String())
    for k, _ := range childOutcomes{
      //fmt.Printf("\n%v|%v", string(k), v)
      // If this is a clear win or loss
      if len(childOutcomes) == 1{
        return k
      }
    }
    return 'Z'
  }
  //fmt.Printf("\nStart node:\n%v\n", gn.String())
  f(gn, 0)
  //SUMMARY
  //fmt.Printf("Final summary\n")
  //for k, v := range m {
    //fmt.Printf("\n%v %v", k, string(v) )
  //}
  return m
}
