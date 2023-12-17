package main

import (
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {
	tests := []struct {
		actions         []int
		expectedAnswers []int
	}{
		{[]int{}, []int{}},
		{[]int{1, 2, 3, -1, -1, -1}, []int{1, 2, 3}},
		{[]int{1, 2, 3, -1, 4, 5, 6, 7, 8, 9, 10, 11, 12, -1, -1, -1}, []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		q := NewRingBufferQueue()
		answers := []int{}
		for _, action := range tt.actions {
			if action > 0 {
				q.push(action)
			}

			if action == -1 {
				answers = append(answers, q.pop())
			}
		}
		if !reflect.DeepEqual(answers, tt.expectedAnswers) {
			t.Errorf("%v; want %v", answers, tt.expectedAnswers)
		}
	}
}
