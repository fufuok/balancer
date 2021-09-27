package balancer

type Balancer interface {
	// Add add an item to be selected.
	Add(item string, weight int)

	// All get all items.
	// RoundRobin/Random/ConsistentHash: []string
	// WeightedRoundRobin/SmoothWeightedRoundRobin: map[string]int
	All() interface{}

	// Select gets next selected item.
	// key is only used for ConsistentHash
	Select(key ...string) string

	// Name load balancer name.
	Name() string

	// Remove remove an item.
	// asClean: clean up or remove only one
	Remove(item string, asClean ...bool) bool

	// RemoveAll remove all items.
	RemoveAll()

	// Reset reset the balancer.
	Reset()

	// Update reinitialize the balancer items.
	// RoundRobin/Random/ConsistentHash: []string
	// WeightedRoundRobin/SmoothWeightedRoundRobin: map[string]int
	Update(items interface{}) bool
}

// Mode defines the selectable balancer algorithm.
type Mode int

const (
	// WeightedRoundRobin is the default balancer algorithm.
	WeightedRoundRobin Mode = iota
	SmoothWeightedRoundRobin
	ConsistentHash
	RoundRobin
	Random
)

// New create a balancer with or without items.
// RoundRobin/Random/ConsistentHash: []string
// WeightedRoundRobin/SmoothWeightedRoundRobin: map[string]int
func New(b Mode, items interface{}) Balancer {
	switch b {
	case SmoothWeightedRoundRobin:
		if v, ok := items.(map[string]int); ok {
			return NewSmoothWeightedRoundRobin(v)
		}
		return NewSmoothWeightedRoundRobin()
	case RoundRobin:
		if v, ok := items.([]string); ok {
			return NewRoundRobin(v)
		}
		return NewRoundRobin()
	case ConsistentHash:
		if v, ok := items.([]string); ok {
			return NewConsistentHash(v)
		}
		return NewConsistentHash()
	case Random:
		if v, ok := items.([]string); ok {
			return NewRandom(v)
		}
		return NewRandom()
	default:
		if v, ok := items.(map[string]int); ok {
			return NewWeightedRoundRobin(v)
		}
		return NewWeightedRoundRobin()
	}
}
