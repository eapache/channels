package channels

import "testing"

func TestInfiniteChannel(t *testing.T) {
	var ch Channel

	ch = NewInfiniteChannel()
	testChannel(t, "infinite channel", ch)

	ch = NewInfiniteChannel()
	testChannelPair(t, "infinite channel", ch, ch)
}
