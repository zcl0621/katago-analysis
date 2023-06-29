package wq

import (
	"bytes"
	"fmt"
	"os"
)

var HoshiString = "." // Can be changed. Used when printing the board.

// A Board contains information about a Go position. It is possible to generate
// boards from nodes in an SGF tree, but modifying boards created in this way
// has no effect on the SGF nodes themselves. Creating boards from nodes is
// relatively costly, and should probably be avoided if batch processing many
// files.

type BoardMove struct {
	X      int
	Y      int
	C      Colour
	Type   int // 1 move 2 stones 3 pass
	Commit MoveCommit
}

type MoveCommit struct {
	HasCommit bool  `json:"-"`
	TimeStamp int64 `json:"t"`
	BTime     int   `json:"b_t,omitempty"`
	WTime     int   `json:"w_t,omitempty"`
	BByo      int   `json:"b_b,omitempty"`
	BByoT     int   `json:"b_bt,omitempty"`
	WByo      int   `json:"w_b,omitempty"`
	WByoT     int   `json:"w_bt,omitempty"`
}

type Board struct {
	Size             int
	Player           Colour
	Ko               string
	KM               float64
	Step             int
	State            [][]Colour
	CapturesBy       map[Colour]int
	Paused           bool
	Move             []BoardMove
	BContinuePass    int
	WContinuePass    int
	ContinuePass     int
	BScore           float64
	WScore           float64
	ControversyCount int
	OwnerShip        []float64
}

// NewBoard returns an empty board of specified size.
func NewBoard(sz int) *Board {

	if sz < 1 || sz > 52 {
		panic(fmt.Sprintf("NewBoard(): bad size %d", sz))
	}

	board := new(Board)

	board.Size = sz
	board.Player = BLACK
	board.ClearKo()

	board.State = make([][]Colour, sz)
	for x := 0; x < sz; x++ {
		board.State[x] = make([]Colour, sz)
	}

	board.CapturesBy = make(map[Colour]int)
	board.CapturesBy[BLACK] = 0 // Not strictly
	board.CapturesBy[WHITE] = 0 // necessary...

	return board
}

// Equals returns true if the two boards are the same, including ko status,
// captures, and next player to move.
func (self *Board) Equals(other *Board) bool {
	if self.Size != other.Size || self.Player != other.Player || self.Ko != other.Ko {
		return false
	}
	if self.CapturesBy[BLACK] != other.CapturesBy[BLACK] || self.CapturesBy[WHITE] != other.CapturesBy[WHITE] {
		return false
	}
	for x := 0; x < self.Size; x++ {
		for y := 0; y < self.Size; y++ {
			if self.State[x][y] != other.State[x][y] {
				return false
			}
		}
	}
	return true
}

// Copy returns a deep copy of the board.
func (self *Board) Copy() *Board {

	ret := new(Board)

	// Easy stuff...

	ret.Size = self.Size
	ret.Player = self.Player
	ret.Ko = self.Ko
	ret.Paused = self.Paused
	// State...

	ret.State = make([][]Colour, ret.Size)
	for x := 0; x < ret.Size; x++ {
		ret.State[x] = make([]Colour, ret.Size)
		for y := 0; y < ret.Size; y++ {
			ret.State[x][y] = self.State[x][y]
		}
	}

	// Captures...

	ret.CapturesBy = make(map[Colour]int)
	ret.CapturesBy[BLACK] = self.CapturesBy[BLACK]
	ret.CapturesBy[WHITE] = self.CapturesBy[WHITE]

	return ret
}

// Get returns the colour at the specified point. The argument should be an SGF
// coordinate, e.g. "dd".
func (self *Board) Get(p string) Colour {
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return EMPTY
	}
	return self.State[x][y]
}

// get_fast is for trusted input.
func (self *Board) get_fast(p string) Colour {
	x := int(p[0]) - 97
	y := int(p[1]) - 97
	if p[0] <= 'Z' {
		x = int(p[0]) - 39
	}
	if p[1] <= 'Z' {
		y = int(p[1]) - 39
	}
	return self.State[x][y]
}

// HasKo returns true if the board has a ko square, on which capture by the
// current player to move is prohibited.
func (self *Board) HasKo() bool {
	return self.Ko != ""
}

// Dump prints the board, and some information about captures and next player.
func (self *Board) Dump() {

	// In all our printing, try to build up the whole
	// string first to avoid jerky printouts...

	s := self.String()
	s += fmt.Sprintf("Captures: %d by Black - %d by White\n", self.CapturesBy[BLACK], self.CapturesBy[WHITE])
	s += fmt.Sprintf("Next to play: %v\n", self.Player.Word())

	fmt.Printf(s)
	os.Stdout.Sync() // Same reasoning.
}

// DumpBoard prints the board; it is equivalent to fmt.Printf(board.String()).
func (self *Board) DumpBoard() {
	fmt.Printf(self.String())
}

// String returns an ASCII representation of the board.
func (self *Board) String() string {

	var b bytes.Buffer

	ko_x, ko_y, _ := ParsePoint(self.Ko, self.Size) // Usually -1, -1

	for y := 0; y < self.Size; y++ {
		for x := 0; x < self.Size; x++ {
			c := self.State[x][y]
			if c == BLACK {
				b.WriteString(" X")
			} else if c == WHITE {
				b.WriteString(" O")
			} else if ko_x == x && ko_y == y {
				b.WriteString(" :")
			} else {
				if IsStarPoint(Point(x, y), self.Size) {
					b.WriteString(" ")
					b.WriteString(HoshiString)
				} else {
					b.WriteString(" .")
				}
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (self *Board) ko_square_finder(p string) string {

	// Only called when we know there is indeed a ko.
	// Argument is the location of the capturing stone that caused it.

	var hits []string

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.Get(a) == EMPTY {
			hits = append(hits, a)
		}
	}

	if len(hits) != 1 {
		panic(fmt.Sprintf("ko_square_finder(): got %d hits", len(hits)))
	}

	return hits[0]
}
