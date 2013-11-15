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
