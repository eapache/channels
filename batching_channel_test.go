package channels

import "testing"

func TestBatchingChannel(t *testing.T) {
	ch := NewBatchingChannel(Infinity)
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()

	i := 0
	for val := range ch.Out() {
		for _, elem := range val.([]interface{}) {
			if i != elem.(int) {
				t.Fatal("batching channel expected", i, "but got", elem.(int))
			}
			i++
		}
	}
}
