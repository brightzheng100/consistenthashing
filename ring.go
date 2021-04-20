package consistenthashing

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

var _ ConsistantHashing = &ring{}

var (
	defaultHashFunc = crc32.ChecksumIEEE
)

type option struct {
	hashFunc HashFunc
}
type Option func(o *option)

// ring is the model of how a consistnet hash looks like
// Within the ring, there are many nodes
type ring struct {
	nodes nodes // the logical/virtual "node"s in the ring

	sync.RWMutex

	nodesMap  map[string]*node   // the logical/virtual "node"s by vname -> node
	memberMap map[string]*Member // the physical "nodes"s by name -> Member
	options   *option            // the configurable options
}

// HashFunc is the injectable hashing function
type HashFunc func(key []byte) uint32

// node is a logical/virtual "node" on the ring.
// A node must belone to one physical node, or "pnode"
type node struct {
	name   string  // logical/virtualnode name
	key    uint32  // the hash key this node has
	member *Member // the member this node belongs to
}

type nodes []node

func (n nodes) Len() int           { return len(n) }
func (n nodes) Less(i, j int) bool { return n[i].key < n[j].key }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

// NewRing is to initialize an empty Ring to start with
func NewRing(opts ...Option) *ring {
	options := &option{
		hashFunc: defaultHashFunc,
	}
	for _, opt := range opts {
		opt(options)
	}
	return &ring{
		nodes:     []node{},
		nodesMap:  make(map[string]*node),
		memberMap: make(map[string]*Member),
		options:   options,
	}
}

// WithHashFunc explicitly sets the hashing func
func WithHashFunc(hashFunc HashFunc) Option {
	return func(o *option) {
		o.hashFunc = hashFunc
	}
}

// Add adds a Member into the ring
// with desired weight which will represent the amount of the logical/virtual "node"s
func (r *ring) Add(member Member) bool {
	var name string = member.Name
	var weight int = member.Weight
	var vname string

	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	// if the Member with this name exists, return without a need to add
	if _, ok := r.memberMap[name]; ok {
		fmt.Printf("The member %s already exists", name)
		return false
	}

	// add Member into the ring's memberMap
	r.memberMap[name] = &member

	// iterate to create desired number of vnodes for this node
	if weight < 1 {
		weight = 1
	}
	for i := 0; i < weight; i++ {
		vname = name + "__" + strconv.Itoa(i)
		node := node{
			name:   vname,
			key:    r.options.hashFunc([]byte(vname)),
			member: &member,
		}
		r.nodes = append(r.nodes, node)
		r.nodesMap[vname] = &node
	}
	sort.Sort(r.nodes)
	return true
}

// Remove removes the member by name from the ring
func (r *ring) Remove(name string) bool {
	var vname string
	var weight int

	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	member, ok := r.memberMap[name]
	if !ok {
		fmt.Printf("The member %s doesn't exist", name)
		return false
	}
	// delete it from the memberMap
	delete(r.memberMap, name)

	// iterate to remove the logical/virtual "node"s on the ring
	weight = member.Weight
	if weight < 1 {
		weight = 1
	}
	for i := 0; i < weight; i++ {
		vname = name + "__" + strconv.Itoa(i)

		// delete it from the nodeMap
		delete(r.nodesMap, vname)

		// delete it from the nodes
		for idx, n := range r.nodes {
			if n.name == vname {
				r.nodes = append(r.nodes[:idx], r.nodes[idx+1:]...)
			}
		}

	}
	return true
}

func (r *ring) GetMembers() []Member {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	members := make([]Member, 0, len(r.memberMap))
	for _, m := range r.memberMap {
		members = append(members, *m)
	}

	return members
}

func (r *ring) Lookup(key string) Member {
	var hashKey uint32
	var hitIndex int

	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	hashKey = r.options.hashFunc([]byte(key))
	index := sort.Search(len(r.nodes), func(i int) bool { return r.nodes[i].key >= hashKey })
	if index == len(r.nodes) {
		//fmt.Println("search to the end!!!")
		hitIndex = 0
	} else {
		hitIndex = index
	}

	node := r.nodes[hitIndex]
	node.member.Hits++

	return *node.member
}
