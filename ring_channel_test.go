package channels

import "testing"

func TestRingChannel(t *testing.T) {
	var ch Channel

	ch = NewRingChannel(Infinity) // yes this is rather silly, but it should work
	testChannel(t, "infinite ring-buffer channel", ch)

	ch = NewRingChannel(1)
	testChannel(t, "single-element ring-buffer channel", ch)

	ch = NewRingChannel(10)
	for i := 0; i < 1000; i++ {
		ch.In() <- i
	}
	ch.Close()
	for i := 990; i < 1000; i++ {
		val := <-ch.Out()
		if i != val.(int) {
			t.Fatal("ring channel expected", i, "but got", val.(int))
		}
	}
	if val := <-ch.Out(); val != nil {
		t.Fatal("ring channel expected closed but got", val)
	}
}
