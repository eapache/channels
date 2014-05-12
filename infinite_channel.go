package channels

import "github.com/eapache/queue"

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type InfiniteChannel struct {
	input, output chan interface{}
	buffer        *queue.Queue
}

func NewInfiniteChannel() *InfiniteChannel {
	ch := &InfiniteChannel{make(chan interface{}), make(chan interface{}), queue.New()}
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
	return ch.buffer.Length()
}

func (ch *InfiniteChannel) Cap() BufferCap {
	return Infinity
}

func (ch *InfiniteChannel) Close() {
	close(ch.input)
}

func (ch *InfiniteChannel) shutdown() {
	for ch.buffer.Length() > 0 {
		ch.output <- ch.buffer.Peek()
		ch.buffer.Remove()
	}
	close(ch.output)
}

func (ch *InfiniteChannel) infiniteBuffer() {
	for {
		if ch.buffer.Length() == 0 {
			elem, open := <-ch.input
			if open {
				ch.buffer.Add(elem)
			} else {
				ch.shutdown()
				return
			}
		} else {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer.Add(elem)
				} else {
					ch.shutdown()
					return
				}
			case ch.output <- ch.buffer.Peek():
				ch.buffer.Remove()
			}
		}
	}
}
