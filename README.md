# A simple consistent hashing design and implementation

## The interface defined in `interface.go`

```go
// Member is the physical machine that joins for consistent hashing
type Member struct {
	Name   string                 // the name
	Addr   string                 // the network address, e.g. ":8080" or "https://server.com"
	Weight int                    // the weight this Member has
	Hits   int64                  // stats for how many hits this Member has after it joins the ring
	Config map[string]interface{} // extra config if any
}

// ConsistantHashing defines the interfaces
type ConsistantHashing interface {
	// Add adds a Member with desired logical/virtual "node"s to the Ring
	Add(member Member) bool
	// Remove removes the named Member from the Ring
	Remove(name string) bool
	// GetMembers gets all current Members
	GetMembers() []Member
	// Lookup looks up a Member by a given key
	Lookup(key string) Member
}
```

## The sample & outcome

I've provided an example under `/examples`.

```sh
$ go run simple.go
```

The output may look like: 

```log
------- Add nodes -------
------- Lookup node by key -------
------- Stats -------
node [machine1] serving [192.168.0.1:8080] got [282536] hits
node [machine2] serving [192.168.0.2:8080] got [220394] hits
node [machine3] serving [192.168.0.3:8080] got [251943] hits
node [machine4] serving [192.168.0.4:8080] got [245127] hits
------- Remove node: machine2 -------
------- Lookup node by key-------
------- Stats -------
node [machine1] serving [192.168.0.1:8080] got [621093] hits
node [machine3] serving [192.168.0.3:8080] got [597263] hits
node [machine4] serving [192.168.0.4:8080] got [561250] hits
------- Add new node: machine5 -------
------- Lookup node by key-------
------- Stats -------
node [machine3] serving [192.168.0.3:8080] got [854567] hits
node [machine4] serving [192.168.0.4:8080] got [800514] hits
node [machine1] serving [192.168.0.1:8080] got [901406] hits
node [machine5] serving [192.168.0.5:8080] got [223119] hits
```
