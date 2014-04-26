package channels

import "reflect"

//sharedBufferChannel implements SimpleChannel and is created by the public
//SharedBuffer type below
type sharedBufferChannel struct {
	in  chan interface{}
	out chan interface{}
	buf *queue
}

func (sch *sharedBufferChannel) In() chan<- interface{} {
	return sch.in
}

func (sch *sharedBufferChannel) Out() <-chan interface{} {
	return sch.out
}

func (sch *sharedBufferChannel) Close() {
	close(sch.in)
}

//SharedBuffer ... TODO
type SharedBuffer struct {
	cases []reflect.SelectCase
	chans []*sharedBufferChannel
	count int
	size  BufferCap
	in    chan *sharedBufferChannel
}

func NewSharedBuffer(size BufferCap) *SharedBuffer {
	if size < 0 && size != Infinity {
		panic("channels: invalid negative size in NewSharedBuffer")
	} else if size == None {
		panic("channels: SharedBuffer does not support unbuffered behaviour")
	}
	buf := &SharedBuffer{
		size: size,
		in:   make(chan *sharedBufferChannel),
	}
	buf.cases = append(buf.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(buf.in),
	})
	go buf.mainLoop()
	return buf
}

func (buf *SharedBuffer) NewChannel() SimpleChannel {
	ch := &sharedBufferChannel{
		in:  make(chan interface{}),
		out: make(chan interface{}),
		buf: newQueue(),
	}
	buf.in <- ch
	return ch
}

func (buf *SharedBuffer) Len() int {
	return buf.count
}

func (buf *SharedBuffer) Cap() BufferCap {
	return buf.size
}

func (buf *SharedBuffer) Close() {
	// TODO: what if there are still active channels using this buffer?
	close(buf.in)
}

func (buf *SharedBuffer) mainLoop() {
	for {
		i, val, ok := reflect.Select(buf.cases)

		if i == 0 {
			if ok {
				//NewChannel was called on the SharedBuffer
				ch := val.Interface().(*sharedBufferChannel)
				buf.chans = append(buf.chans, ch)
				buf.cases = append(buf.cases,
					reflect.SelectCase{Dir: reflect.SelectRecv},
					reflect.SelectCase{Dir: reflect.SelectSend},
				)
				if buf.size == Infinity || buf.count < int(buf.size) {
					buf.cases[len(buf.cases)-2].Chan = reflect.ValueOf(ch.in)
				}
			} else {
				//Close was called on the SharedBuffer itself
				//TODO
			}
		} else {
			if i%2 == 0 {
				//Send
				if buf.count == int(buf.size) {
					//room in the buffer again
					//TODO re-enable recv cases
				}
				buf.count--
				ch := buf.chans[(i-1)/2]
				if ch.buf.length() > 0 {
					buf.cases[i].Send = reflect.ValueOf(ch.buf.peek())
					ch.buf.remove()
				} else {
					buf.cases[i].Chan = reflect.Value{}
					buf.cases[i].Send = reflect.Value{}
				}
			} else {
				//Receive or Close
				//TODO
			}
		}
	}
}
