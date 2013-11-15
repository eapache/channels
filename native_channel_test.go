package channels

import "testing"

func TestNativeChannels(t *testing.T) {
	var ch Channel

	ch = NewNativeChannel(NoBuffer)
	testChannel(t, "bufferless native channel", ch)

	ch = NewNativeChannel(NoBuffer)
	testChannelPair(t, "bufferless native channel", ch, ch)

	ch = NewNativeChannel(5)
	testChannel(t, "5-buffer native channel", ch)

	ch = NewNativeChannel(5)
	testChannelPair(t, "5-buffer native channel", ch, ch)
}
