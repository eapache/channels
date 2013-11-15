package channels

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

type NativeChannel chan interface{}

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
