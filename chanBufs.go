package chanBufs

// NewInfiniteBuffer produces a write-only channel and a read-only channel that
// behave exactly like the two ends of a single go channel with an infinite buffer in between.
// Be very careful using this, as no buffer is truly infinite - if the internal
// buffer grows too large your program will run out of memory and crash.
func NewInfiniteBuffer() (input chan<- interface{}, output <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})

	go func() {
		var buffer []interface{}
		for {
			if len(buffer) == 0 {
				elem, open := <-in
				if open {
					buffer = append(buffer, elem)
				} else {
					close(out)
					return
				}
			} else {
				select {
				case elem, open := <-in:
					if open {
						buffer = append(buffer, elem)
					} else {
						for elem := range buffer {
							out <- elem
						}
						close(out)
						return
					}
				case out <- buffer[0]:
					buffer = buffer[1:]
				}
			}
		}
	}()

	return in, out
}

// NewAdjustableBuffer produces a write-only channel and a read-only channel that
// behave exactly like the two ends of a single go channel with an adjustable buffer in between.
// The channel initially has a buffer size of 1, but can be adjusted by sending the new desired size
// down the third, 'adjust' channel.
//
// Setting a size of 0 would ideally produce an unbuffered, blocking channel but that does not appear
// to be possible with this approach, so we do the best we can which is a buffer of 1.
//
// Setting a negative size produces an infinite buffer.
//
// It is an error to close the adjust channel, it will be closed automatically when input is closed.
// It is an error to write to the adjust channel after closing the input channel.
func NewAdjustableBuffer() (input chan<- interface{}, output <-chan interface{}, adjust chan<- int) {
	in := make(chan interface{})
	out := make(chan interface{})
	adj := make(chan int)
	size := 1

	go func() {
		var buffer []interface{}
		for {
			if len(buffer) == 0 {
				select {
				case elem, open := <-in:
					if open {
						buffer = append(buffer, elem)
					} else {
						close(out)
						close(adj)
						return
					}
				case size = <-adj:
				}
			} else if size >= 0 && len(buffer) >= size {
				select {
				case out <- buffer[0]:
					buffer = buffer[1:]
				case size = <-adj:
				}
			} else {
				select {
				case elem, open := <-in:
					if open {
						buffer = append(buffer, elem)
					} else {
						for elem := range buffer {
							out <- elem
						}
						close(out)
						close(adj)
						return
					}
				case out <- buffer[0]:
					buffer = buffer[1:]
				case size = <-adj:
				}
			}
		}
	}()

	return in, out, adj
}
