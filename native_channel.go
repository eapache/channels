package channels

// NativeInChannel implements the InChannel interface by wrapping a native go write-only channel
type NativeInChannel chan<- interface{}

func (ch NativeInChannel) In() chan<- interface{} {
	return ch
}

func (ch NativeInChannel) Len() int {
	return len(ch)
}

func (ch NativeInChannel) Cap() int {
	return cap(ch)
}

func (ch NativeInChannel) Close() {
	close(ch)
}

// NativeOutChannel implements the OutChannel interface by wrapping a native go read-only channel
type NativeOutChannel <-chan interface{}

func (ch NativeOutChannel) Out() <-chan interface{} {
	return ch
}

func (ch NativeOutChannel) Len() int {
	return len(ch)
}

func (ch NativeOutChannel) Cap() int {
	return cap(ch)
}

// NativeChannel implements the Channel interface by wrapping a native go channel
type NativeChannel chan interface{}

// NewNativeChannel makes a new NativeChannel with the given buffer size. Just a convenience wrapper
// to avoid calling make() and then casting the result.
func NewNativeChannel(buffer int) NativeChannel {
	return make(chan interface{}, buffer)
}

func (ch NativeChannel) In() chan<- interface{} {
	return ch
}

func (ch NativeChannel) Out() <-chan interface{} {
	return ch
}

func (ch NativeChannel) Len() int {
	return len(ch)
}

func (ch NativeChannel) Cap() int {
	return cap(ch)
}

func (ch NativeChannel) Close() {
	close(ch)
}
