package channels

import "testing"

// General warning: Go's benchmark utility (go test -bench .) increases the number of
// iterations until the benchmarks take a reasonable amount of time to run; memory usage
// is *NOT* considered. On my machine, these benchmarks hit around ~1GB before they've had
// enough, but if you have less than that available and start swapping, then all bets are off.

func BenchmarkQueueSerial(b *testing.B) {
	q := newQueue()
	for i := 0; i < b.N; i++ {
		q.add(nil)
	}
	for i := 0; i < b.N; i++ {
		q.remove()
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := newQueue()
	for i := 0; i < b.N; i++ {
		q.add(nil)
		q.remove()
	}
}
