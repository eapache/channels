package channels

import "testing"

func BenchmarkQueueSerial(b *testing.B) {
	q := newQueue()
	for i := 0; i < b.N; i++ {
		q.add(i)
	}
	for i := 0; i < b.N; i++ {
		q.remove()
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := newQueue()
	for i := 0; i < b.N; i++ {
		q.add(i)
		q.remove()
	}
}
