package channels

import (
	"sort"
	"sync"
)

const minQueueLen = 16

type bufCache struct {
	cache [][]interface{}
	count int
	lck   sync.Mutex
}

var gBufCache = &bufCache{} // global cached buffors

func (bc *bufCache) get(n int) (result []interface{}) {
	bc.lck.Lock()
	defer bc.lck.Unlock()

	if n < minQueueLen {
		n = minQueueLen
	}

	cacheLen := len(bc.cache)
	if cacheLen == 0 {
		bc.count++
		return make([]interface{}, n)
	}

	pos := sort.Search(cacheLen, func(x int) bool { return len(bc.cache[x]) >= n })
	if pos == cacheLen {
		bc.count++
		return make([]interface{}, n)
	}

	result = bc.cache[pos]
	copy(bc.cache[pos:], bc.cache[pos+1:])
	bc.cache[cacheLen-1] = nil
	bc.cache = bc.cache[:cacheLen-1]

	return
}

func (bc *bufCache) put(buf []interface{}) {
	bc.lck.Lock()
	defer bc.lck.Unlock()

	bc.count--

	buf = buf[:cap(buf)]
	bufLen := len(buf)
	if bufLen == 0 {
		return
	}

	cacheLen := len(bc.cache)
	pos := sort.Search(cacheLen, func(x int) bool { return len(bc.cache[x]) >= bufLen })
	bc.cache = append(bc.cache, nil)
	copy(bc.cache[pos+1:], bc.cache[pos:])
	bc.cache[pos] = buf

	// if every requested buffer had came back,
	// assume we do not need cache anymore, flush now
	if bc.count == 0 && len(bc.cache) > 0 {
		bc.cache = nil
	}
}

// A fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
// Using this instead of a simple slice+append provides substantial memory and time
// benefits, and fewer GC pauses.
type queue struct {
	buf               []interface{}
	head, tail, count int
}

func newQueue() *queue {
	return &queue{buf: gBufCache.get(minQueueLen)}
}

func (q *queue) length() int {
	return q.count
}

func (q *queue) resize() {
	newBuf := gBufCache.get(q.count * 2)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else if q.count > 0 {
		copy(newBuf, q.buf[q.head:len(q.buf)])
		copy(newBuf[len(q.buf)-q.head:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	gBufCache.put(q.buf)
	q.buf = newBuf
}

func (q *queue) enqueue(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
}

func (q *queue) peek() interface{} {
	return q.buf[q.head]
}

func (q *queue) dequeue() {
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	if len(q.buf) > minQueueLen && q.count*4 <= len(q.buf) {
		q.resize()
	}
}
