package channels

// Pipe connects the input channel to the output channel so that
// they behave as if a single channel.
func Pipe(input OutChannel, output InChannel) {
	go func() {
		for elem := range input.Out() {
			output.In() <- elem
		}
		out.Close()
	}()
}
