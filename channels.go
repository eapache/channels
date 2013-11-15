/*
Package channels provides a collection of helper functions, interfaces and special types for
working with and extending the capabilities of golang's existing channels.

Several types in this package provide so-called "infinite" buffers. Be *very* careful using these, as no
buffer is truly infinite - if such a buffer grows too large your program will run out of memory and crash.
Caveat emptor.
*/
package channels

const (
	NoBuffer       int = 0
	InfiniteBuffer int = -1
)

// InChannel is an interface representing a writeable channel.
type InChannel interface {
	In() chan<- interface{}
	Len() int
	Cap() int
	Close()
}

// OutChannel is an interface representing a readable channel.
type OutChannel interface {
	Out() <-chan interface{}
	Len() int
	Cap() int
}

// Channel is an interface representing a channel that is both readable and writeable.
type Channel interface {
	InChannel
	Out() <-chan interface{}
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
