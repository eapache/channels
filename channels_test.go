package channels

import (
	"math/rand"
	"testing"
	"time"
)

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

func TestMultiplex(t *testing.T) {
	a := NewNativeChannel(None)
	b := NewNativeChannel(None)

	Multiplex(b, a)

	testChannelPair(t, "simple multiplex", a, b)

	a = NewNativeChannel(None)
	inputs := []Channel{
		NewNativeChannel(None),
		NewNativeChannel(None),
		NewNativeChannel(None),
		NewNativeChannel(None),
	}

	Multiplex(a, inputs[0], inputs[1], inputs[2], inputs[3])

	go func() {
		rand.Seed(time.Now().Unix())
		for i := 0; i < 1000; i++ {
			inputs[rand.Intn(len(inputs))].In() <- i
		}
		for i := range inputs {
			inputs[i].Close()
		}
	}()
	for i := 0; i < 1000; i++ {
		val := <-a.Out()
		if i != val.(int) {
			t.Fatal("multiplexing expected", i, "but got", val.(int))
		}
	}
}
