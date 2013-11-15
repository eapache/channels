package chanBufs

// NewInfiniteBuffer produces a write-only channel and a read-only channel that
// behave exactly like the two ends of a single go channel with an infinite buffer in between.
// Be very careful using this, as no buffer is truly infinite - if the internal
// buffer grows too large your program will run out of memory and crash.
func NewInfiniteBuffer() (chan<- interface{}, <-chan interface{}) {
	input := make(chan interface{})
	output := make(chan interface{})

	go func() {
		var buffer []interface{}
		for {
			if len(buffer) == 0 {
				elem, open := <-input
				if open {
					buffer = append(buffer, elem)
				} else {
					close(output)
					return
				}
			} else {
				select {
				case elem, open := <-input:
					if open {
						buffer = append(buffer, elem)
					} else {
						for elem := range buffer {
							output <- elem
						}
						close(output)
						return
					}
				case output <- buffer[0]:
					buffer = buffer[1:]
				}
			}
		}
	}()

	return input, output
}
