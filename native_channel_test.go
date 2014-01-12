package channels

import "testing"

func TestNativeChannels(t *testing.T) {
	var ch Channel

	ch = NewNativeChannel(None)
	testChannel(t, "bufferless native channel", ch)

	ch = NewNativeChannel(None)
	testChannelPair(t, "bufferless native channel", ch, ch)

	ch = NewNativeChannel(5)
	testChannel(t, "5-buffer native channel", ch)

	ch = NewNativeChannel(5)
	testChannelPair(t, "5-buffer native channel", ch, ch)
}

func TestDeadChannel(t *testing.T) {
	ch := NewDeadChannel()

	select {
	case <-ch.Out():
		t.Error("read from a dead channel")
	default:
	}

	select {
	case ch.In() <- nil:
		t.Error("wrote to a dead channel")
	default:
	}

	ch.Close()
}
