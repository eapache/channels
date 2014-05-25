package channels

import (
	"github.com/eapache/queue"
	"reflect"
)

//sharedBufferChannel implements SimpleChannel and is created by the public
//SharedBuffer type below
type sharedBufferChannel struct {
	in  chan interface{}
	out chan interface{}
	buf *queue.Queue
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
	cases []reflect.SelectCase   // 2n+1 of these; [0] is for control, [i%2==0] for send, [i%2==1] for recv
	chans []*sharedBufferChannel // n of these
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
		buf: queue.New(),
	}
	buf.in <- ch
	return ch
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
					//room in the buffer again, re-enable all recv cases
					for j := range buf.chans {
						buf.cases[(j*2)+1].Chan = reflect.ValueOf(buf.chans[j].in)
					}
				}
				buf.count--
				ch := buf.chans[(i-1)/2]
				if ch.buf.Length() > 0 {
					buf.cases[i].Send = reflect.ValueOf(ch.buf.Peek())
					ch.buf.Remove()
				} else {
					//nothing left for this channel to send, disable that case
					buf.cases[i].Chan = reflect.Value{}
					buf.cases[i].Send = reflect.Value{}
				}
			} else {
				if ok {
					//Receive
					buf.count++
					ch := buf.chans[i/2]
					if ch.buf.Length() == 0 && !buf.cases[i].Chan.IsValid() {
						//this channel now has something to send
						buf.cases[i].Chan = reflect.ValueOf(ch.out)
						buf.cases[i].Send = val
					} else {
						ch.buf.Add(val.Interface())
					}
					if buf.count == int(buf.size) {
						//buffer full, disable recv cases
						for j := range buf.chans {
							buf.cases[(j*2)+1].Chan = reflect.Value{}
						}
					}
				} else {
					//Close
					//TODO
				}
			}
		}
	}
}

func (buf *SharedBuffer) Len() int {
	return buf.count
}

func (buf *SharedBuffer) Cap() BufferCap {
	return buf.size
}
