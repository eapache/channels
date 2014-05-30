package channels

import "testing"

func TestSharedBufferSingleton(t *testing.T) {
	buf := NewSharedBuffer(3)

	ch := buf.NewChannel()
	for i := 0; i < 5; i++ {
		ch.In() <- (*int)(nil)
		ch.In() <- (*int)(nil)
		ch.In() <- (*int)(nil)
		select {
		case ch.In() <- (*int)(nil):
			t.Error("Wrote to full shared-buffer")
		default:
		}

		<-ch.Out()
		<-ch.Out()
		<-ch.Out()
		select {
		case <-ch.Out():
			t.Error("Read from empty shared-buffer")
		default:
		}
	}

	ch.Close()
	buf.Close()
}
