package internal

import (
	"errors"
	"strings"
)

func WithPrefixIfNotPresent(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s
	}
	return p + s
}

// A slice window provides an iterator over a slice in form of windows of fixed width
// It is intended to be created by the NewSliceWindow function.
type SliceWindow[T any] struct {
	slice               []T
	currentWindowNumber int
	windowSize          int
}

// NewSliceWindow is a function which creates a new SliceWindow on the given slice with a window size.
// It returns a SliceWindow if the windowSize is positive and not zero. Valid values are 1, 2, 3, 4, ...
func NewSliceWindow[T any](slice []T, windowSize int) (SliceWindow[T], error) {
	if windowSize <= 0 {
		return SliceWindow[T]{}, errors.New("The window size must be positive and not zero!")
	}
	return SliceWindow[T]{
		slice:               slice,
		currentWindowNumber: 0,
		windowSize:          int(windowSize),
	}, nil
}

// Returns the current window on the underlying slice.
// This method returns an empty slice if you finished iterating over the array. Please use the HasNext method to check
// if you can continue iterating over it. And use NextWindow to continue to the next window in order.
func (sw *SliceWindow[T]) Window() []T {
	lowerWindowIndex := sw.currentWindowNumber * sw.windowSize
	upperWindowIndex := (sw.currentWindowNumber + 1) * sw.windowSize
	if len(sw.slice) < upperWindowIndex {
		upperWindowIndex = len(sw.slice)
	}

	if sw.HasNext() {
		return sw.slice[lowerWindowIndex:upperWindowIndex]
	}

	return make([]T, 0)
}

// Returns a boolean if there is a next window on the slice available or not.
func (sw *SliceWindow[T]) HasNext() bool {
	result := len(sw.slice) > (sw.windowSize * sw.currentWindowNumber)
	return result
}

// Moves the current window position by one.
func (sw *SliceWindow[T]) NextWindow() {
	sw.currentWindowNumber += 1
}
