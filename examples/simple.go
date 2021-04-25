package main

import (
	"fmt"
	"hash/crc32"
	"strconv"

	ch "github.com/brightzheng100/consistenthashing"
)

func main() {
	//ring := ch.NewConsistentHashing()
	ring := ch.NewConsistentHashing(ch.WithHashFunc(crc32.ChecksumIEEE))

	fmt.Println("------- Add nodes -------")

	ring.Add(ch.Member{
		Name:   "machine1",
		Addr:   "192.168.0.1:8080",
		Weight: 1032,
	})
	ring.Add(ch.Member{
		Name:   "machine2",
		Addr:   "192.168.0.2:8080",
		Weight: 1024,
	})
	ring.Add(ch.Member{
		Name:   "machine3",
		Addr:   "192.168.0.3:8080",
		Weight: 1032,
	})
	ring.Add(ch.Member{
		Name:   "machine4",
		Addr:   "192.168.0.4:8080",
		Weight: 1016,
	})

	fmt.Println("------- Lookup node by key -------")
	for i := 0; i < 1000000; i++ {
		key := "my_test_key_" + strconv.Itoa(i)
		_ = ring.Lookup(key)
		//fmt.Printf("key [%s] matches node [%s]\n", key, name)
	}

	fmt.Println("------- Stats -------")
	// stats
	for _, m := range ring.GetMembers() {
		fmt.Printf("node [%s] serving [%s] got [%d] hits\n", m.Name, m.Addr, m.Hits)
	}

	fmt.Println("------- Remove node: machine2 -------")

	ring.Remove("machine2")

	fmt.Println("------- Lookup node by key-------")
	for i := 0; i < 1000000; i++ {
		key := "my_test_key_" + strconv.Itoa(i)
		_ = ring.Lookup(key)
		//fmt.Printf("key [%s] matches node [%s]\n", key, name)
	}

	fmt.Println("------- Stats -------")
	// stats
	for _, m := range ring.GetMembers() {
		fmt.Printf("node [%s] serving [%s] got [%d] hits\n", m.Name, m.Addr, m.Hits)
	}

	fmt.Println("------- Add new node: machine5 -------")

	ring.Add(ch.Member{
		Name:   "machine5",
		Addr:   "192.168.0.5:8080",
		Weight: 1024,
	})

	fmt.Println("------- Lookup node by key-------")
	for i := 0; i < 1000000; i++ {
		key := "my_test_key_" + strconv.Itoa(i)
		_ = ring.Lookup(key)
		//fmt.Printf("key [%s] matches node [%s]\n", key, name)
	}

	fmt.Println("------- Stats -------")
	// stats
	for _, m := range ring.GetMembers() {
		fmt.Printf("node [%s] serving [%s] got [%d] hits\n", m.Name, m.Addr, m.Hits)
	}
}
