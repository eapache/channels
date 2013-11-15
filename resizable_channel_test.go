package channels

import "testing"

func TestResizableChannel(t *testing.T) {
	var ch *ResizableChannel

	ch = NewResizableChannel()
	testChannel(t, "default resizable channel", ch)

	ch = NewResizableChannel()
	testChannelPair(t, "default resizable channel", ch, ch)

	ch = NewResizableChannel()
	ch.Resize(InfiniteBuffer)
	testChannel(t, "infinite resizable channel", ch)

	ch = NewResizableChannel()
	ch.Resize(InfiniteBuffer)
	testChannelPair(t, "infinite resizable channel", ch, ch)

	ch = NewResizableChannel()
	ch.Resize(5)
	testChannel(t, "5-buffer resizable channel", ch)

	ch = NewResizableChannel()
	ch.Resize(5)
	testChannelPair(t, "5-buffer resizable channel", ch, ch)
}
