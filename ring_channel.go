package channels

// RingChannel implements the Channel interface in a way that never blocks the writer.
// Specifically, if a value is written to a RingChannel when its buffer is full then the oldest
// value in the buffer is discarded to make room (just like a standard ring-buffer).
// Note that Go's scheduler can cause discarded values when they could be avoided, simply by scheduling
// the writer before the reader. This happens a lot particularly with small buffer sizes, so caveat emptor.
// For the opposite behaviour (discarding the newest element, not the oldest) see OverflowingChannel.
type RingChannel struct {
	input, output chan interface{}
	buffer        []interface{}
	size          BufferCap
}

func NewRingChannel(size BufferCap) *RingChannel {
	if size < 0 && size != Infinity {
		panic("channels: invalid negative size in NewRingChannel")
	}
	ch := &RingChannel{make(chan interface{}), make(chan interface{}), nil, size}
	if size == None {
		go ch.overflowingDirect()
	} else {
		go ch.ringBuffer()
	}
	return ch
}

func (ch *RingChannel) In() chan<- interface{} {
	return ch.input
}

func (ch *RingChannel) Out() <-chan interface{} {
	return ch.output
}

func (ch *RingChannel) Len() int {
	return len(ch.buffer)
}

func (ch *RingChannel) Cap() BufferCap {
	return ch.size
}

func (ch *RingChannel) Close() {
	close(ch.input)
}

func (ch *RingChannel) shutdown() {
	for _, elem := range ch.buffer {
		ch.output <- elem
	}
	close(ch.output)
}

// for entirely unbuffered cases
func (ch *RingChannel) overflowingDirect() {
	for elem := range ch.input {
		// if we can't write it immediately, drop it and move on
		select {
		case ch.output <- elem:
		default:
		}
	}
}

// for all buffered cases
func (ch *RingChannel) ringBuffer() {
	for {
		if len(ch.buffer) == 0 {
			elem, open := <-ch.input
			if open {
				ch.buffer = append(ch.buffer, elem)
			} else {
				ch.shutdown()
				return
			}
		} else {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer = append(ch.buffer, elem)
					if ch.size != Infinity && len(ch.buffer) > int(ch.size) {
						ch.buffer = ch.buffer[1:]
					}
				} else {
					ch.shutdown()
					return
				}
			case ch.output <- ch.buffer[0]:
				ch.buffer = ch.buffer[1:]
			}
		}
	}
}
