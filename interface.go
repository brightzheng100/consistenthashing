package consistenthashing

// Member is the physical machine that joins for consistent hashing
type Member struct {
	Name   string                 // the name
	Addr   string                 // the network address, e.g. ":8080" or "https://server.com"
	Weight int                    // the weight this Member has
	Hits   int64                  // stats for how many hits this Member has after it joins the ring
	Config map[string]interface{} // extra config if any
}

// ConsistentHashing defines the interfaces
type ConsistentHashing interface {
	// Add adds a Member with desired logical/virtual "node"s to the Ring
	Add(member Member) bool
	// Remove removes the named Member from the Ring
	Remove(name string) bool
	// GetMembers gets all current Members
	GetMembers() []Member
	// Lookup looks up a Member by a given key
	Lookup(key string) Member
}
