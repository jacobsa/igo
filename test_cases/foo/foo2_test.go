package foo

import "testing"

func TestComputeThreeTimesThree(t *testing.T) {
	if ComputeThreeTimesThree() != 9 {
		t.Errorf("Expected nine.")
	}
}
