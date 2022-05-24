package stack

import (
	"testing"
)

func TestStack(t *testing.T) {
	var s Stack
	s.Push("abcd")
	for !s.IsEmpty() {
		result := s.Pop()
		if result != "abcd" {
			t.Errorf("Stack get error, get %s", result)
		}
	}
	r := s.Pop()
	if r != "" {
		t.Errorf("Stack get not empty, get %s", r)
	}
}
