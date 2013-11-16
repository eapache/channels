package channels

import "testing"

func testChannel(t *testing.T, name string, ch Channel) {
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()
	for i := 0; i < 1000; i++ {
		val := <-ch.Out()
		if i != val.(int) {
			t.Fatal(name, "expected", i, "but got", val.(int))
		}
	}
}

func testChannelPair(t *testing.T, name string, in InChannel, out OutChannel) {
	go func() {
		for i := 0; i < 1000; i++ {
			in.In() <- i
		}
		in.Close()
	}()
	for i := 0; i < 1000; i++ {
		val := <-out.Out()
		if i != val.(int) {
			t.Fatal("pair", name, "expected", i, "but got", val.(int))
		}
	}
}

func TestPipe(t *testing.T) {
	a := NewNativeChannel(None)
	b := NewNativeChannel(None)

	Pipe(a, b)

	testChannelPair(t, "pipe", a, b)
}
