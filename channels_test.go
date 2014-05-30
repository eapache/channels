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

func TestWeakPipe(t *testing.T) {
	a := NewNativeChannel(None)
	b := NewNativeChannel(None)

	WeakPipe(a, b)

	testChannelPair(t, "pipe", a, b)
}

func testMultiplex(t *testing.T, multi func(output SimpleInChannel, inputs ...SimpleOutChannel)) {
	a := NewNativeChannel(None)
	b := NewNativeChannel(None)

	multi(b, a)

	testChannelPair(t, "simple multiplex", a, b)

	a = NewNativeChannel(None)
	inputs := []Channel{
		NewNativeChannel(None),
		NewNativeChannel(None),
		NewNativeChannel(None),
		NewNativeChannel(None),
	}

	multi(a, inputs[0], inputs[1], inputs[2], inputs[3])

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

func TestMultiplex(t *testing.T) {
	testMultiplex(t, Multiplex)
}

func TestWeakMultiplex(t *testing.T) {
	testMultiplex(t, WeakMultiplex)
}

func testTee(t *testing.T, tee func(input SimpleOutChannel, outputs ...SimpleInChannel)) {
	a := NewNativeChannel(None)
	b := NewNativeChannel(None)

	tee(a, b)

	testChannelPair(t, "simple tee", a, b)

	a = NewNativeChannel(None)
	outputs := []Channel{
		NewNativeChannel(None),
		NewNativeChannel(None),
		NewNativeChannel(None),
		NewNativeChannel(None),
	}

	tee(a, outputs[0], outputs[1], outputs[2], outputs[3])

	go func() {
		for i := 0; i < 1000; i++ {
			a.In() <- i
		}
		a.Close()
	}()
	for i := 0; i < 1000; i++ {
		for _, output := range outputs {
			val := <-output.Out()
			if i != val.(int) {
				t.Fatal("teeing expected", i, "but got", val.(int))
			}
		}
	}
}

func TestTee(t *testing.T) {
	testTee(t, Tee)
}

func TestWeakTee(t *testing.T) {
	testTee(t, WeakTee)
}

func ExampleChannel() {
	var ch Channel

	ch = NewInfiniteChannel()

	for i := 0; i < 10; i++ {
		ch.In() <- nil
	}

	for i := 0; i < 10; i++ {
		<-ch.Out()
	}
}
