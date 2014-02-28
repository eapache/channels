package channels

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type InfiniteChannel struct {
	input, output     chan interface{}
	buffer            []interface{}
	head, tail, count int
}

func NewInfiniteChannel() *InfiniteChannel {
	ch := &InfiniteChannel{
		input:  make(chan interface{}),
		output: make(chan interface{}),
		buffer: make([]interface{}, 16),
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
	return ch.count
}

func (ch *InfiniteChannel) Cap() BufferCap {
	return Infinity
}

func (ch *InfiniteChannel) Close() {
	close(ch.input)
}

func (ch *InfiniteChannel) shutdown() {
	for ch.count > 0 {
		ch.output <- ch.dequeue()
	}
	close(ch.output)
}

func (ch *InfiniteChannel) enqueue(elem interface{}) {
	if ch.head == ch.tail && ch.count > 0 {
		buffer := make([]interface{}, len(ch.buffer)*2)

		copy(buffer, ch.buffer[ch.head:])
		copy(buffer[len(ch.buffer)-ch.head:], ch.buffer[:ch.head])

		ch.head = 0
		ch.tail = len(ch.buffer)
		ch.buffer = buffer
	}

	ch.buffer[ch.tail] = elem
	ch.tail = (ch.tail + 1) % len(ch.buffer)
	ch.count++
}

func (ch *InfiniteChannel) dequeue() interface{} {
	elem := ch.buffer[ch.head]
	ch.head = (ch.head + 1) % len(ch.buffer)
	ch.count--

	return elem
}

func (ch *InfiniteChannel) infiniteBuffer() {
	for {
		if ch.count == 0 {
			elem, open := <-ch.input
			if open {
				ch.enqueue(elem)
			} else {
				ch.shutdown()
				return
			}
		} else {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.enqueue(elem)
				} else {
					ch.shutdown()
					return
				}
			case ch.output <- ch.buffer[ch.head]: // can't use dequeue directly here
				ch.head = (ch.head + 1) % len(ch.buffer)
				ch.count--
			}
		}
	}
}
