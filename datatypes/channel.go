package datatypes

import (
	"errors"
	"fmt"
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
	ui.Debug(ui.ByteCodeLogger, "--> Created  %s", c.String())

	return c
}

// Send transmits an arbitrary data object through the channel, if it
// is open.
func (c *Channel) Send(datum interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isOpen {
		ui.Debug(ui.ByteCodeLogger, "--> Sending on %s", c.String())
		c.channel <- datum
		return nil
	}
	return errors.New(ChannelNotOpenError)
}

// Receive accepts an arbitrary data object through the channel, waiting
// if there is no information avaialble yet.
func (c *Channel) Receive() (interface{}, error) {

	if !c.isOpen {
		return nil, errors.New(ChannelNotOpenError)
	}
	ui.Debug(ui.ByteCodeLogger, "--> Receiving on %s", c.String())
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
	defer c.mutex.Unlock()

	ui.Debug(ui.ByteCodeLogger, "--> Closing %s", c.String())
	wasActive := c.isOpen
	close(c.channel)
	c.isOpen = false
	return wasActive
}

func (c *Channel) String() string {
	state := "open"
	if !c.isOpen {
		state = "closed"
	}
	return fmt.Sprintf("channel, size %d, %s, id %s",
		c.size, state, c.id)
}
