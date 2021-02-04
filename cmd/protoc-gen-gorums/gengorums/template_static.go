// Code generated by protoc-gen-gorums. DO NOT EDIT.
// Source files can be found in: ./cmd/protoc-gen-gorums/dev

package gengorums

// pkgIdentMap maps from package name to one of the package's identifiers.
// These identifiers are used by the Gorums protoc plugin to generate
// appropriate import statements.
var pkgIdentMap = map[string]string{
	"fmt":                             "Errorf",
	"github.com/relab/gorums":         "ContentSubtype",
	"google.golang.org/grpc/encoding": "GetCodec",
	"sort":                            "Search",
}

// reservedIdents holds the set of Gorums reserved identifiers.
// These identifiers cannot be used to define message types in a proto file.
var reservedIdents = []string{
	"Configuration",
	"Manager",
	"NewManager",
	"Node",
	"QuorumSpec",
}

var staticCode = `// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	id    uint32
	nodes []*gorums.Node
	mgr   *Manager
	qspec QuorumSpec
	errs  chan gorums.Error
}

// NewConfig returns a configuration for the given node addresses and quorum spec.
// The returned func() must be called to close the underlying connections.
// This is an experimental API.
func NewConfig(qspec QuorumSpec, opts ...gorums.ManagerOption) (*Configuration, func(), error) {
	man, err := NewManager(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create manager: %v", err)
	}
	c, err := man.NewConfiguration(man.NodeIDs(), qspec)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create configuration: %v", err)
	}
	return c, func() { man.Close() }, nil
}

// ID reports the identifier for the configuration.
func (c *Configuration) ID() uint32 {
	return c.id
}

// NodeIDs returns a slice containing the local ids of all the nodes in the
// configuration. IDs are returned in the same order as they were provided in
// the creation of the Configuration.
func (c *Configuration) NodeIDs() []uint32 {
	ids := make([]uint32, len(c.nodes))
	for i, node := range c.nodes {
		ids[i] = node.ID()
	}
	return ids
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Configuration.
func (c *Configuration) Nodes() []*Node {
	nodes := make([]*Node, 0, len(c.nodes))
	for _, n := range c.nodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// Size returns the number of nodes in the configuration.
func (c *Configuration) Size() int {
	return len(c.nodes)
}

func (c *Configuration) String() string {
	return fmt.Sprintf("config-%d", c.id)
}

// Equal returns a boolean reporting whether a and b represents the same
// configuration.
func Equal(a, b *Configuration) bool { return a.id == b.id }

// SubError returns a channel for listening to individual node errors. Currently
// only a single listener is supported.
func (c *Configuration) SubError() <-chan gorums.Error {
	return c.errs
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

func NewManager(opts ...gorums.ManagerOption) (mgr *Manager, err error) {
	mgr = &Manager{}
	mgr.Manager, err = gorums.NewManager(opts...)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

type Manager struct {
	*gorums.Manager
}

func (m *Manager) NewConfiguration(ids []uint32, qspec QuorumSpec) (*Configuration, error) {
	if len(ids) == 0 {
		return nil, gorums.IllegalConfigError("need at least one node")
	}

	var nodes []*gorums.Node
	unique := make(map[uint32]struct{})
	for _, nid := range ids {
		// ensure that identical IDs are only counted once
		if _, duplicate := unique[nid]; duplicate {
			continue
		}
		unique[nid] = struct{}{}

		node, found := m.Node(nid)
		if !found {
			return nil, gorums.NodeNotFoundError(nid)
		}

		i := sort.Search(len(nodes), func(i int) bool {
			return node.ID() < nodes[i].ID()
		})
		nodes = append(nodes, nil)
		copy(nodes[i+1:], nodes[i:])
		nodes[i] = node
	}

	c := &Configuration{
		nodes: nodes,
		mgr:   m,
		qspec: qspec,
	}
	return c, nil
}

type Node struct {
	*gorums.Node
}

`
