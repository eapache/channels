package channels

import "testing"

func TestBlackHole(t *testing.T) {
	discard := NewBlackHole()

	for i := 0; i < 1000; i++ {
		discard.In() <- i
	}

	discard.Close()

	if discard.Len() != 1000 {
		t.Error("blackhole expected 1000 was", discard.Len())
	}
}
