package channels

// BlackHole implements the InChannel interface and provides an analogue for the "Discard" variable in
// the ioutil package - it never blocks, and simply discards every value it reads. The number of items
// discarded in this way is counted and returned from Len.
type BlackHole struct {
	input   chan interface{}
	stopper chan struct{}
	count   int
}

func NewBlackHole() *BlackHole {
	ch := &BlackHole{make(chan interface{}), make(chan struct{}), 0}
	go ch.discard()
	return ch
}

func (ch *BlackHole) In() chan<- interface{} {
	return ch.input
}

func (ch *BlackHole) Len() int {
	return ch.count
}

func (ch *BlackHole) Cap() int {
	return InfiniteBuffer
}

func (ch *BlackHole) Close() {
	close(ch.input)
	<-ch.stopper
}

func (ch *BlackHole) discard() {
	for _ = range ch.input {
		ch.count += 1
	}
	close(ch.stopper)
}
