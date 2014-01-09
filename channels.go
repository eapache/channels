/*
Package channels provides a collection of helper functions, interfaces and implementations for
working with and extending the capabilities of golang's existing channels.

The general interface provided is Channel, though sub-interfaces are also provided for cases
where the full Channel interface cannot be met (for example, InChannel for write-only channels).

Helper functions include Pipe and Tee (which behave much like their Unix namesakes), as well as Multiplex.
Weak versions of these functions also exist, which do not close their output channel on completion.

A simple wrapper type called NativeChannel is included for wrapping native golang channels in the appropriate
interface. Several special implementations of the Channel interface are also provided, including channels backed
by special buffers (resizable, infinite, ring buffers, etc) and other useful types. A black hole channel for
discarding unwanted values (similar in purpose to ioutil.Discard or /dev/null) rounds out the set.

Several types in this package provide so-called "infinite" buffers. Be *very* careful using these, as no
buffer is truly infinite - if such a buffer grows too large your program will run out of memory and crash.
Caveat emptor.
*/
package channels

import "reflect"

// BufferCap represents the capacity of the buffer backing a channel. Valid values consist of all
// positive integers, as well as the special values below.
type BufferCap int

const (
	None     BufferCap = 0
	Infinity BufferCap = -1
)

// Buffer is an interface for any channel that provides access to query the state of its buffer.
// Even unbuffered channels can implement this interface by simply returning 0 from Len() and None from Cap().
type Buffer interface {
	Len() int       // The number of elements currently buffered.
	Cap() BufferCap // The maximum number of elements that can be buffered.
}

// SimpleInChannel is an interface representing a writeable channel that does not necessarily
// implement the Buffer interface.
type SimpleInChannel interface {
	In() chan<- interface{} // The writeable end of the channel.
	Close()                 // Closes the channel. It is an error to write to In() after calling Close().
}

// InChannel is an interface representing a writeable channel with a buffer.
type InChannel interface {
	SimpleInChannel
	Buffer
}

// SimpleOutChannel is an interface representing a readable channel that does not necessarily
// implement the Buffer interface.
type SimpleOutChannel interface {
	Out() <-chan interface{} // The readable end of the channel.
}

// OutChannel is an interface representing a readable channel implementing the Buffer interface.
type OutChannel interface {
	SimpleOutChannel
	Buffer
}

// SimpleChannel is an interface representing a channel that is both readable and writeable,
// but does not necessarily implement the Buffer interface.
type SimpleChannel interface {
	SimpleInChannel
	SimpleOutChannel
}

// Channel is an interface representing a channel that is readable, writeable and implements
// the Buffer interface
type Channel interface {
	SimpleChannel
	Buffer
}

// Pipe connects the input channel to the output channel so that
// they behave as if a single channel.
func Pipe(input SimpleOutChannel, output SimpleInChannel) {
	go func() {
		for elem := range input.Out() {
			output.In() <- elem
		}
		output.Close()
	}()
}

// Multiplex takes an arbitrary number of input channels and multiplexes their output into a single output
// channel. When all input channels have been closed, the output channel is closed. Multiplex with a single
// input channel is equivalent to Pipe (though slightly less efficient).
func Multiplex(output SimpleInChannel, inputs ...SimpleOutChannel) {
	if len(inputs) == 0 {
		panic("channels: Multiplex requires at least one input")
	}
	go func() {
		inputCount := len(inputs)
		cases := make([]reflect.SelectCase, inputCount)
		for i := range cases {
			cases[i].Dir = reflect.SelectRecv
			cases[i].Chan = reflect.ValueOf(inputs[i].Out())
		}
		for inputCount > 0 {
			chosen, recv, recvOK := reflect.Select(cases)
			if recvOK {
				output.In() <- recv.Interface()
			} else {
				cases[chosen].Chan = reflect.ValueOf(nil)
				inputCount -= 1
			}
		}
		output.Close()
	}()
}

// Tee (like its Unix namesake) takes a single input channel and an arbitrary number of output channels
// and duplicates each input into every output. When the input channel is closed, all outputs channels are closed.
// Tee with a single output channel is equivalent to Pipe (though slightly less efficient).
func Tee(input SimpleOutChannel, outputs ...SimpleInChannel) {
	if len(outputs) == 0 {
		panic("channels: Tee requires at least one output")
	}
	go func() {
		cases := make([]reflect.SelectCase, len(outputs))
		for i := range cases {
			cases[i].Dir = reflect.SelectSend
		}
		for elem := range input.Out() {
			for i := range cases {
				cases[i].Chan = reflect.ValueOf(outputs[i].In())
				cases[i].Send = reflect.ValueOf(elem)
			}
			for remaining := len(cases); remaining > 0; remaining -= 1 {
				chosen, _, _ := reflect.Select(cases)
				cases[chosen].Chan = reflect.ValueOf(nil)
			}
		}
		for i := range outputs {
			outputs[i].Close()
		}
	}()
}

// WeakPipe behaves like Pipe (connecting the two channels) except that it does not close
// the output channel when the input channel is closed.
func WeakPipe(input SimpleOutChannel, output SimpleInChannel) {
	go func() {
		for elem := range input.Out() {
			output.In() <- elem
		}
	}()
}

// WeakMultiplex behaves like Multiplex (multiplexing multiple inputs into a single output) except that it does not close
// the output channel when the input channels are closed.
func WeakMultiplex(output SimpleInChannel, inputs ...SimpleOutChannel) {
	if len(inputs) == 0 {
		panic("channels: WeakMultiplex requires at least one input")
	}
	go func() {
		inputCount := len(inputs)
		cases := make([]reflect.SelectCase, inputCount)
		for i := range cases {
			cases[i].Dir = reflect.SelectRecv
			cases[i].Chan = reflect.ValueOf(inputs[i].Out())
		}
		for inputCount > 0 {
			chosen, recv, recvOK := reflect.Select(cases)
			if recvOK {
				output.In() <- recv.Interface()
			} else {
				cases[chosen].Chan = reflect.ValueOf(nil)
				inputCount -= 1
			}
		}
	}()
}

// WeakTee behaves like Tee (duplicating a single input into multiple outputs) except that it does not close
// the output channels when the input channel is closed.
func WeakTee(input SimpleOutChannel, outputs ...SimpleInChannel) {
	if len(outputs) == 0 {
		panic("channels: WeakTee requires at least one output")
	}
	go func() {
		cases := make([]reflect.SelectCase, len(outputs))
		for i := range cases {
			cases[i].Dir = reflect.SelectSend
		}
		for elem := range input.Out() {
			for i := range cases {
				cases[i].Chan = reflect.ValueOf(outputs[i].In())
				cases[i].Send = reflect.ValueOf(elem)
			}
			for remaining := len(cases); remaining > 0; remaining -= 1 {
				chosen, _, _ := reflect.Select(cases)
				cases[chosen].Chan = reflect.ValueOf(nil)
			}
		}
	}()
}
