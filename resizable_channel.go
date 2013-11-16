package channels

// ResizableChannel implements the Channel interface with a resizable buffer between the input and the output.
// The channel initially has a buffer size of 1, but can be resized by calling Resize().
//
// Resizing to a buffer capacity of None is, unfortunately, not supported and will panic.
// Resizing back and forth between a finite and infinite buffer is fully supported.
type ResizableChannel struct {
	input, output chan interface{}
	resize        chan BufferCap
	size          BufferCap
	buffer        []interface{}
}

func NewResizableChannel() *ResizableChannel {
	ch := &ResizableChannel{make(chan interface{}), make(chan interface{}), make(chan BufferCap), 1, nil}
	go ch.magicBuffer()
	return ch
}

func (ch *ResizableChannel) In() chan<- interface{} {
	return ch.input
}

func (ch *ResizableChannel) Out() <-chan interface{} {
	return ch.output
}

func (ch *ResizableChannel) Len() int {
	return len(ch.buffer)
}

func (ch *ResizableChannel) Cap() BufferCap {
	return ch.size
}

func (ch *ResizableChannel) Close() {
	close(ch.input)
}

func (ch *ResizableChannel) Resize(newSize BufferCap) {
	if newSize == None {
		panic("channels: ResizableChannel does not support unbuffered behaviour")
	}
	if newSize < 0 && newSize != Infinity {
		panic("channels: invalid negative size trying to resize channel")
	}
	ch.resize <- newSize
}

func (ch *ResizableChannel) shutdown() {
	for _, elem := range ch.buffer {
		ch.output <- elem
	}
	close(ch.output)
	close(ch.resize)
}

func (ch *ResizableChannel) magicBuffer() {
	for {
		if len(ch.buffer) == 0 {
			select {
			case elem, open := <-ch.input:
				if open {
					ch.buffer = append(ch.buffer, elem)
				} else {
					ch.shutdown()
					return
				}
			case ch.size = <-ch.resize:
			}
		} else if ch.size != Infinity && len(ch.buffer) >= int(ch.size) {
			select {
			case ch.output <- ch.buffer[0]:
				ch.buffer = ch.buffer[1:]
			case ch.size = <-ch.resize:
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
			case ch.size = <-ch.resize:
			}
		}
	}
}
