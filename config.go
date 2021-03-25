package gorums

import (
	"fmt"

	"github.com/relab/gorums/ordering"
)

// Configuration represents a static set of nodes on which quorum calls may be invoked.
type Configuration []*Node

// NewConfiguration returns a configuration based on the provided list of nodes.
// Nodes can be supplied using WithNodeMap or WithNodeList or WithNodeIDs.
func NewConfiguration(mgr *Manager, opt NodeListOption) (nodes Configuration, err error) {
	if opt == nil {
		return nil, ConfigCreationError(fmt.Errorf("missing required node list"))
	}
	return opt.newConfig(mgr)
}

// Channels returns the default channels for each node in the configuration.
func (c Configuration) Channels() []*Channel {
	channels := make([]*Channel, len(c))
	for i, node := range c {
		channels[i] = node.Channel()
	}
	return channels
}

// NewChannels creates new channels for each node in the configuration.
func (c Configuration) NewChannels() (channels []*Channel, err error) {
	channels = make([]*Channel, len(c))
	for i, node := range c {
		channels[i], err = node.NewChannel()
		if err != nil {
			return nil, err
		}
	}
	return channels, nil
}

// NodeIDs returns a slice of this configuration's Node IDs.
func (c Configuration) NodeIDs() []uint32 {
	ids := make([]uint32, len(c))
	for i, node := range c {
		ids[i] = node.ID()
	}
	return ids
}

// Nodes returns the nodes in this configuration.
func (c Configuration) Nodes() []*Node {
	return c
}

// Size returns the number of nodes in this configuration.
func (c Configuration) Size() int {
	return len(c)
}

// Equal returns true if configurations b and c have the same set of nodes.
func (c Configuration) Equal(b Configuration) bool {
	if len(c) != len(b) {
		return false
	}
	for i := range c {
		if c[i].ID() != b[i].ID() {
			return false
		}
	}
	return true
}

// newCall returns unique metadata for a method call.
func (c Configuration) newCall(method string) (md *ordering.Metadata) {
	// Note that we just use the first node's newCall method since all nodes
	// associated with the same manager use the same receiveQueue instance.
	return c[0].newCall(method)
}

// newReply returns a channel for receiving replies
// and a done function to be called for clean up.
func (c Configuration) newReply(md *ordering.Metadata, maxReplies int) (replyChan chan *response, done func()) {
	// Note that we just use the first node's newReply method since all nodes
	// associated with the same manager use the same receiveQueue instance.
	return c[0].newReply(md, maxReplies)
}
