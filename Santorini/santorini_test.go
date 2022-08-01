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
	want := "St|040030000200A1B3030X4Y0111124|End"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
