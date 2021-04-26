package consistenthashing

import (
	"hash/crc32"
	"strconv"
	"testing"
)

type testCase struct {
	name   string
	addr   string
	weight int
	hits   int64
}

// test data
var (
	cases = []testCase{
		{name: "machine1", addr: "192.168.0.1:8080", weight: 24},
		{name: "machine2", addr: "192.168.0.2:8080", weight: 32},
		{name: "machine3", addr: "192.168.0.3:8080", weight: 48},
	}

	negativeWeightCase = testCase{
		name: "abc", addr: "efg", weight: -1,
	}
)

func TestNew(t *testing.T) {
	ch := NewConsistentHashing()

	var members []Member
	members = ch.GetMembers()

	if len(members) != 0 {
		t.Errorf("newly created object should have 0 member, but got %d", len(members))
	}
}

func TestNewWithOptions(t *testing.T) {
	ch := NewConsistentHashing(WithHashFunc(crc32.ChecksumIEEE))

	var members []Member
	members = ch.GetMembers()

	if len(members) != 0 {
		t.Errorf("newly created object should have 0 member, but got %d", len(members))
	}
}

func TestAdd(t *testing.T) {
	ch := NewConsistentHashing()

	var members []Member

	for c := range cases {
		ch.Add(Member{
			Name:   cases[c].name,
			Addr:   cases[c].addr,
			Weight: cases[c].weight,
		})

		members = ch.GetMembers()

		// the amount of members
		assertMembersAreEqual(t, len(members), c+1)

		// assess cases and members
		assertMembersMatchCases(t, cases, members)
	}
}

func TestAddWithNegativeWeight(t *testing.T) {
	ch := NewConsistentHashing()

	var members []Member

	ch.Add(Member{
		Name:   negativeWeightCase.name,
		Addr:   negativeWeightCase.addr,
		Weight: negativeWeightCase.weight,
	})

	members = ch.GetMembers()

	// the amount of members
	assertMembersAreEqual(t, len(members), 1)

	// assess the member
	if negativeWeightCase.name != members[0].Name {
		t.Errorf("member's Name was set as %s, but got %s", negativeWeightCase.name, members[0].Name)
	}
	if negativeWeightCase.addr != members[0].Addr {
		t.Errorf("member's Addr was set as %s, but got %s", negativeWeightCase.addr, members[0].Addr)
	}
	if members[0].Weight != 1 {
		t.Errorf("member's Weight should be reset as 1 when negative, but got %d", members[0].Weight)
	}
}

func TestAddThenGetMembers(t *testing.T) {
	ch := initConsistentHashing()

	var members []Member

	members = ch.GetMembers()

	assertMembersAreEqual(t, len(members), len(cases))

	// assess cases and members
	assertMembersMatchCases(t, cases, members)
}

func TestAddRemoveThenAdd(t *testing.T) {
	ch := initConsistentHashing()

	var members []Member

	members = ch.GetMembers()

	assertMembersAreEqual(t, len(members), len(cases))

	// remove machine2
	ch.Remove("machine2")
	members = ch.GetMembers()

	assertMembersAreEqual(t, len(members), len(cases)-1)

	newCases := append(cases[:1], cases[2])
	// assess cases and members
	assertMembersMatchCases(t, newCases, members)

	// add machine4
	ch.Add(Member{
		Name:   "machine4",
		Addr:   "192.168.0.4:8080",
		Weight: 1024,
	})
	members = ch.GetMembers()
	assertMembersAreEqual(t, len(members), len(cases)-1+1)

	newCases = append(newCases, testCase{
		name: "machine4", addr: "192.168.0.4:8080", weight: 1024,
	})
	// assess cases and members
	assertMembersMatchCases(t, newCases, members)
}

func TestAddExistingMember(t *testing.T) {
	ch := initConsistentHashing()

	var happened bool
	var members []Member
	var before, after int

	members = ch.GetMembers()
	before = len(members)

	happened = ch.Add(Member{
		Name:   "machine1",
		Addr:   "whatever",
		Weight: 1,
	})

	members = ch.GetMembers()
	after = len(members)

	if before != after || happened {
		t.Errorf("Adding existing member should be skipped")
	}
}

func TestRemoveNotExistMember(t *testing.T) {
	ch := initConsistentHashing()

	var happened bool
	var members []Member
	var before, after int

	members = ch.GetMembers()
	before = len(members)

	happened = ch.Remove("nothingexists")

	members = ch.GetMembers()
	after = len(members)

	if before != after || happened {
		t.Errorf("Removing not existing member should be skipped")
	}
}

func TestLookup(t *testing.T) {
	ch := initConsistentHashing()

	var members []Member
	var member Member
	var hit bool

	members = ch.GetMembers()

	for i := 0; i < 1000; i++ {
		hit = false

		key := "my_test_key_" + strconv.Itoa(i)
		member = ch.Lookup(key)

		// member must within the members
		for _, m := range members {
			if member.Name == m.Name && member.Addr == m.Addr && member.Weight == m.Weight {
				hit = true
			}
		}

		if !hit {
			t.Errorf("can't lookup proper member with key %s", key)
		}
	}
}

func initConsistentHashing() ConsistentHashing {
	ch := NewConsistentHashing()

	for _, member := range cases {
		ch.Add(Member{
			Name:   member.name,
			Addr:   member.addr,
			Weight: member.weight,
		})
	}

	return ch
}

func assertMembersAreEqual(t *testing.T, expected int, got int) {
	if expected != got {
		t.Errorf("members must be %d but got %d", expected, got)
	}
}

func assertMembersMatchCases(t *testing.T, cases []testCase, members []Member) {
	var member Member

	// the members retrieved from GetMembers may or may not be sequential
	for i := range cases {
		for j := range members {
			member = members[j]
			if cases[i].name == member.Name {
				if cases[i].name != member.Name {
					t.Errorf("member's Name was set as %s, but got %s", cases[i].name, member.Name)
				}
				if cases[i].addr != member.Addr {
					t.Errorf("member's Addr was set as %s, but got %s", cases[i].addr, member.Addr)
				}
				if cases[i].weight != member.Weight {
					t.Errorf("member's Weight was set as %d, but got %d", cases[i].weight, member.Weight)
				}
			}
		}
	}
}
