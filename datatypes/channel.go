package datatypes

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/ui"
)

// Structure of an Ego channel wrapper around Go channels
type Channel struct {
	channel chan interface{}
	mutex   sync.Mutex
	size    int
	isOpen  bool
	id      string
}

const (
	ChannelNotOpenError = "channel not open"
)

// Create a mew instance of an Ego channel. The size passed indicates
// the buffer size, which is 1 unless size is greater than 1, in which
// case it is set to the given size.
func NewChannel(size int) *Channel {
	if size < 1 {
		size = 1
	}
	c := &Channel{
		isOpen: true,
		size:   size,
		mutex:  sync.Mutex{},
		id:     uuid.New().String(),
	}
	c.channel = make(chan interface{}, size)
	ui.Debug(ui.ByteCodeLogger, "--> Created channel %s", c.id)

	return c
}

// Send transmits an arbitrary data object through the channel, if it
// is open.
func (c *Channel) Send(datum interface{}) error {
	if c.isOpen {
		ui.Debug(ui.ByteCodeLogger, "--> Sending on channel %s", c.id)
		c.channel <- datum
		return nil
	}
	return errors.New(ChannelNotOpenError)
}

// Receive accepts an arbitrary data object through the channel, waiting
// if there is no information avaialble yet.
func (c *Channel) Receive() (interface{}, error) {
	ui.Debug(ui.ByteCodeLogger, "--> Receiving on channel %s", c.id)
	datum := <-c.channel
	return datum, nil
}

// Return a boolean value indicating if this channel is still open for
// business.
func (c *Channel) IsOpen() bool {
	return c.isOpen
}

func (c *Channel) GetSize() int {
	return c.size
}

// Close the channel so no more sends are permitted to the channel, and
// the receiver can test for channel completion
func (c *Channel) Close() bool {
	c.mutex.Lock()
	wasActive := c.isOpen
	close(c.channel)
	c.isOpen = false
	c.mutex.Unlock()
	return wasActive
}
