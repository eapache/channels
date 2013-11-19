package channels

// OverflowingChannel implements the Channel interface in a way that never blocks.
// Specifically, if a value is written to an OverflowingChannel when its buffer is full
// (or, in an unbuffered case, when the recipient is not ready) then that value is simply discarded.
// Note that Go's scheduler can cause discarded values when they could be avoided, simply by scheduling
// the writer before the reader. This happens a lot particularly with small buffer sizes, so caveat emptor.
type OverflowingChannel struct {
	input, output chan interface{}
	buffer        []interface{}
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
	return len(ch.buffer)
}

func (ch *OverflowingChannel) Cap() BufferCap {
	return ch.size
}

func (ch *OverflowingChannel) Close() {
	close(ch.input)
}

func (ch *OverflowingChannel) shutdown() {
	for _, elem := range ch.buffer {
		ch.output <- elem
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
		if len(ch.buffer) == 0 {
			elem, open := <-ch.input
			if open {
				ch.buffer = append(ch.buffer, elem)
			} else {
				ch.shutdown()
				return
			}
		} else if ch.size != Infinity && len(ch.buffer) >= int(ch.size) {
			select {
			case _, open := <-ch.input: // discard new inputs
				if !open {
					ch.shutdown()
					return
				}
			case ch.output <- ch.buffer[0]:
				ch.buffer = ch.buffer[1:]
			}
		} else {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer = append(ch.buffer, elem)
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
