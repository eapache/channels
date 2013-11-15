package channels

// BatchingChannel implements the Channel interface, with the change that instead of producing individual elements
// on Out(), it batches together the entire internal buffer each time. Trying to construct an unbuffered batching channel
// will panic, that configuration is not supported (and provides no benefit over an unbuffered NativeChannel).
type BatchingChannel struct {
	input, output chan interface{}
	buffer        []interface{}
	size          int
}

func NewBatchingChannel(size int) *BatchingChannel {
	if size == NoBuffer {
		panic("channels: BatchingChannel does not support NoBuffer")
	}
	if size < 0 && size != InfiniteBuffer {
		panic("channels: invalid negative size in NewBatchingChannel")
	}
	ch := &BatchingChannel{make(chan interface{}), make(chan interface{}), nil, size}
	go ch.batchingBuffer()
	return ch
}

func (ch *BatchingChannel) In() chan<- interface{} {
	return ch.input
}

func (ch *BatchingChannel) Out() <-chan interface{} {
	return ch.output
}

func (ch *BatchingChannel) Len() int {
	return len(ch.buffer)
}

func (ch *BatchingChannel) Cap() int {
	return ch.size
}

func (ch *BatchingChannel) Close() {
	close(ch.input)
}

func (ch *BatchingChannel) shutdown() {
	ch.output <- ch.buffer
	close(ch.output)
}

func (ch *BatchingChannel) batchingBuffer() {
	for {
		if len(ch.buffer) == 0 {
			elem, open := <-ch.input
			if open {
				ch.buffer = append(ch.buffer, elem)
			} else {
				ch.shutdown()
				return
			}
		} else if ch.size != InfiniteBuffer && len(ch.buffer) >= ch.size {
			ch.output <- ch.buffer
			ch.buffer = nil
		} else {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer = append(ch.buffer, elem)
				} else {
					ch.shutdown()
					return
				}
			case ch.output <- ch.buffer:
				ch.buffer = nil
			}
		}
	}
}
