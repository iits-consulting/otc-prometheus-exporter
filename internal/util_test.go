package internal

import (
	"golang.org/x/exp/slices"
	"testing"
)

func TestNewSliceWindow(t *testing.T) {

	input := []int{1, 2, 3, 4, 5, 6}
	if _, err := NewSliceWindow(input, 0); err == nil {
		t.Errorf("It was expected to return an error because the windowSize was 0")
	}

	if _, err := NewSliceWindow(input, -1); err == nil {
		t.Errorf("It was expected to return an error because the windowSize was -1")
	}

	if _, err := NewSliceWindow(input, 5); err != nil {
		t.Errorf("It was expected not to return an error because the windowSize was valid")
	}
}

func TestSliceWindow(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	windowSize := 3
	expectedWindows := [][]int{
		{1, 2, 3}, {4, 5, 6},
	}

	sw, err := NewSliceWindow(input, windowSize)
	if err != nil {
		t.Fatalf("No SliceWindow has been returned because of error:" + err.Error())
	}

	sw.currentWindowNumber = 0
	i := 0
	for sw.HasNext() {
		window := sw.Window()
		if !slices.Equal(window, expectedWindows[i]) {
			t.Fatalf("Missmatch in window %d", i)
		}
		i++
		sw.NextWindow()
	}

}
