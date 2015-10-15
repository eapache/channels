package channels

import "github.com/eapache/queue"

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type InfiniteChannel struct {
	input, output chan interface{}
	length        chan int
	buffer        *queue.Queue
}

func NewInfiniteChannel() *InfiniteChannel {
	ch := &InfiniteChannel{
		input:  make(chan interface{}),
		output: make(chan interface{}),
		length: make(chan int),
		buffer: queue.New(),
	}
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
	return <-ch.length
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
	close(ch.length)
}

func (ch *InfiniteChannel) infiniteBuffer() {
	for {
		if ch.buffer.Length() == 0 {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer.Add(elem)
				} else {
					ch.shutdown()
					return
				}
			case ch.length <- ch.buffer.Length():
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
			case ch.length <- ch.buffer.Length():
			}
		}
	}
}
