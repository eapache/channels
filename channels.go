/*
Package channels provides a collection of helper functions, interfaces and special types for
working with and extending the capabilities of golang's existing channels.

Several types in this package provide so-called "infinite" buffers. Be *very* careful using these, as no
buffer is truly infinite - if such a buffer grows too large your program will run out of memory and crash.
Caveat emptor.
*/
package channels

// BufferCap represents the capacity of the buffer backing a channel. Valid values consist of all
// positive integers, as well as the special values below.
type BufferCap int

const (
	None     BufferCap = 0
	Infinity BufferCap = -1
)

// InChannel is an interface representing a writeable channel.
type InChannel interface {
	In() chan<- interface{} // The writeable end of the channel.
	Len() int               // The number of elements currently buffered.
	Cap() BufferCap         // The size of the backing buffer.
	Close()                 // Closes the channel.
}

// OutChannel is an interface representing a readable channel.
type OutChannel interface {
	Out() <-chan interface{} // The readable end of the channel.
	Len() int                // The number of elements currently buffered.
	Cap() BufferCap          // The size of the backing buffer.
}

// Channel is an interface representing a channel that is both readable and writeable.
type Channel interface {
	In() chan<- interface{}  // The writeable end of the channel.
	Out() <-chan interface{} // The readable end of the channel.
	Len() int                // The number of elements currently buffered.
	Cap() BufferCap          // The size of the backing buffer.
	Close()                  // Closes the channel.
}

// Pipe connects the input channel to the output channel so that
// they behave as if a single channel.
func Pipe(input OutChannel, output InChannel) {
	go func() {
		for elem := range input.Out() {
			output.In() <- elem
		}
		output.Close()
	}()
}
