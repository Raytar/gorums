package dev

import "reflect"

func appendIfNotPresent(set []uint32, x uint32) []uint32 {
	for _, y := range set {
		if y == x {
			return set
		}
	}
	return append(set, x)
}

// Future is an object that will return a result in the future
type Future interface {
	getc() chan struct{}
}

// Select blocks until one of the futures is ready, and returns the index of that future
func Select(futures []Future) int {
	selectCases := make([]reflect.SelectCase, 0, len(futures))
	for _, fut := range futures {
		selectCases = append(selectCases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(fut.getc()),
		})
	}
	i, _, _ := reflect.Select(selectCases)
	return i
}
