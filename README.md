# A simple consistent hashing design and implementation

Consistant hashing is quite common in distributed systems, like `memcached`, to distribute workloads to different worker nodes.

The idea is to assign each joined worker node to a `ring` structure so that the workload can lookup for the near worker node, from a clockwise direction, based on its hash key value.

## The design

From design perspective, we define a struct `Member` to represent a physical machine that joins for consistant hashing.

The `Member` will have the name, network address, weight, hits and some other configurable properties.

And we define 4 key methods to be implemented as the interface.

```go
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
```


## The implementation

The implementation follows a typical `ring` structure, with a configurable hash function as of now.
Current default hash function is `crc32.ChecksumIEEE`.

Each `Member` has a weight which represents the virtual/logical `node`s distributed to the `ring`.
I'd recommend to have a proper amount of virtual/logical `node`s as that will lead to a better distribution result.

It's a thread safe implementation.


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

## The test

The code has 100% test coverage:

```sh
$ go test -coverprofile cp.out
PASS
coverage: 100.0% of statements
ok  	github.com/brightzheng100/consistenthashing	2.389s
```
