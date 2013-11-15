package channels

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
// Be very careful using this, as no buffer is truly infinite - if the internal buffer grows too large your
// program will run out of memory and crash.
type InfiniteChannel struct {
	input, output chan interface{}
	buffer        []interface{}
}

func NewInfiniteChannel() *InfiniteChannel {
	ch := &InfiniteChannel{make(chan interface{}), make(chan interface{}), nil}
	go ch.infiniteBuffer()
	return ch
}

func (ch *InfiniteChannel) In() chan<- interface{} {
	return ch.input
}

func (ch *InfiniteChannel) Out() <-chan interface{} {
	return ch.output
}

func (ch *InfiniteChannel) Len() int {
	return len(ch.buffer)
}

func (ch *InfiniteChannel) Cap() int {
	return InfiniteBuffer
}

func (ch *InfiniteChannel) Close() {
	close(ch.input)
}

func (ch *InfiniteChannel) shutdown() {
	for elem := range ch.buffer {
		ch.output <- elem
	}
	close(ch.output)
}

func (ch *InfiniteChannel) infiniteBuffer() {
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
