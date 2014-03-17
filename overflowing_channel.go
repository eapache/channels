package channels

// OverflowingChannel implements the Channel interface in a way that never blocks the writer.
// Specifically, if a value is written to an OverflowingChannel when its buffer is full
// (or, in an unbuffered case, when the recipient is not ready) then that value is simply discarded.
// Note that Go's scheduler can cause discarded values when they could be avoided, simply by scheduling
// the writer before the reader. This happens a lot particularly with small buffer sizes, so caveat emptor.
// For the opposite behaviour (discarding the oldest element, not the newest) see RingChannel.
type OverflowingChannel struct {
	input, output chan interface{}
	buffer        *queue
	size          BufferCap
}

func NewOverflowingChannel(size BufferCap) *OverflowingChannel {
	if size < 0 && size != Infinity {
		panic("channels: invalid negative size in NewOverflowingChannel")
	}
	ch := &OverflowingChannel{make(chan interface{}), make(chan interface{}), nil, size}
	if size == None {
		go ch.overflowingDirect()
	} else {
		ch.buffer = newQueue()
		go ch.overflowingBuffer()
	}
	return ch
}

func (ch *OverflowingChannel) In() chan<- interface{} {
	return ch.input
}

func (ch *OverflowingChannel) Out() <-chan interface{} {
	return ch.output
}

func (ch *OverflowingChannel) Len() int {
	return ch.buffer.length()
}

func (ch *OverflowingChannel) Cap() BufferCap {
	return ch.size
}

func (ch *OverflowingChannel) Close() {
	close(ch.input)
}

func (ch *OverflowingChannel) shutdown() {
	for ch.buffer.length() > 0 {
		ch.output <- ch.buffer.peek()
		ch.buffer.remove()
	}
	close(ch.output)
}

// for entirely unbuffered cases
func (ch *OverflowingChannel) overflowingDirect() {
	for elem := range ch.input {
		// if we can't write it immediately, drop it and move on
		select {
		case ch.output <- elem:
		default:
		}
	}
}

// for all buffered cases
func (ch *OverflowingChannel) overflowingBuffer() {
	for {
		if ch.buffer.length() == 0 {
			elem, open := <-ch.input
			if open {
				ch.buffer.add(elem)
			} else {
				ch.shutdown()
				return
			}
		} else if ch.size != Infinity && ch.buffer.length() >= int(ch.size) {
			select {
			// Prefer to write if possible, which is surprisingly effective in reducing
			// dropped elements due to overflow. The naive read/write select chooses randomly
			// when both channels are ready, which produces unnecessary drops 50% of the time.
			case ch.output <- ch.buffer.peek():
				ch.buffer.remove()
			default:
				select {
				case _, open := <-ch.input: // discard new inputs
					if !open {
						ch.shutdown()
						return
					}
				case ch.output <- ch.buffer.peek():
					ch.buffer.remove()
				}
			}
		} else {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer.add(elem)
				} else {
					ch.shutdown()
					return
				}
			case ch.output <- ch.buffer.peek():
				ch.buffer.remove()
			}
		}
	}
}
