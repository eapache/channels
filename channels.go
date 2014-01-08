/*
Package channels provides a collection of helper functions, interfaces and special types for
working with and extending the capabilities of golang's existing channels.

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
