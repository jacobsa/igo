package foo

import "testing"

func TestComputeTwoPlusTwo(t *testing.T) {
	if ComputeTwoPlusTwo() != 4 {
		t.Errorf("Expected four.")
	}
}
