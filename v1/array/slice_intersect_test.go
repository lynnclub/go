package array

import (
	"reflect"
	"testing"
)

func TestIntersectDeduplicatesElementsWithinEachSlice(t *testing.T) {
	got := Intersect(
		[]int{1, 1, 2, 3},
		[]int{1, 2, 2, 4},
		[]int{1, 1, 2, 5},
	)

	if len(got) != 2 || !In(got, 1) || !In(got, 2) {
		t.Fatalf("Intersect returned %v, want the set {1, 2}", got)
	}
}

func TestIntersectDoesNotModifyInputSliceHeaders(t *testing.T) {
	b := make([][]int, 2, 3)
	b[0] = []int{2, 3}
	b[1] = []int{3, 4}
	before := append([][]int(nil), b...)
	backing := b[:cap(b)]
	backing[2] = []int{99}

	got := Intersect([]int{1, 2, 3}, b...)

	if len(got) != 1 || got[0] != 3 {
		t.Fatalf("Intersect returned %v, want [3]", got)
	}
	if !reflect.DeepEqual(b, before) {
		t.Fatalf("Intersect modified the input slice headers: got %v, want %v", b, before)
	}
	if !reflect.DeepEqual(backing[2], []int{99}) {
		t.Fatalf("Intersect modified the unused backing-array element: got %v", backing[2])
	}
}
