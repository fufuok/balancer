package balancer

type Balancer interface {
	// Add add an item to be selected.
	// weight is only used for WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand, default: 1
	Add(item string, weight ...int)

	// All get all items.
	// RoundRobin/Random/ConsistentHash: []string
	// WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand: map[string]int
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
	// WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand: map[string]int
	Update(items interface{}) bool
}

// Mode defines the selectable balancer algorithm.
type Mode int

const (
	// WeightedRoundRobin is the default balancer algorithm.
	WeightedRoundRobin Mode = iota
	SmoothWeightedRoundRobin
	WeightedRand
	ConsistentHash
	RoundRobin
	Random
)

// New create a balancer with or without items.
// RoundRobin/Random/ConsistentHash: []string
// WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand: map[string]int
func New(b Mode, itemsMap map[string]int, itemsList []string) Balancer {
	switch b {
	case SmoothWeightedRoundRobin:
		return NewSmoothWeightedRoundRobin(itemsMap)
	case WeightedRand:
		return NewWeightedRand(itemsMap)
	case ConsistentHash:
		return NewConsistentHash(itemsList)
	case RoundRobin:
		return NewRoundRobin(itemsList)
	case Random:
		return NewRandom(itemsList)
	default:
		return NewWeightedRoundRobin(itemsMap)
	}
}
