package Santorini

import (
	"testing"
)

func TestPositionToString(t *testing.T) {
	testPosition1 := Position{
		0x1FAB212, 0x182A212, 0x102A012, 0x1020002,
		occupancy[11],
		occupancy[12],
		occupancy[16],
		occupancy[17],
	}

	got := testPosition1.String()
	want := "|0400300002001303040111124|11121617|"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestNewPosition(t *testing.T) {
  position, e := NewPosition("|0400300002001303040111124|08050018|")
  if e != nil{
    t.Errorf("Error forming position")
  }
  got := position.String()
  want := "|0400300002001303040111124|08050018|"
  if got != want {
	t.Errorf("\ngot\t\t\t%q\n, wanted \t%q", got, want)
  }
}
