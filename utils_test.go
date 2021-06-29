package main

import (
	"testing"
)

func TestReverseList(t *testing.T) {
	actual := []int{1, 2, 3, 4, 5}
	expected := []int{5, 4, 3, 2, 1}
	ReverseList(actual)
	
	if len(actual) != len(expected) {
		t.Errorf("%v should have the same length as %v", actual, expected)
	}
	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("%v should be equal to %v", actual, expected)
		}
	}
}
