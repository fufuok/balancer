package balancer

// DefaultBalancer is an global balancer
var DefaultBalancer = NewWeightedRoundRobin()

// Add add an item to be selected.
func Add(item string, weight ...int) {
	DefaultBalancer.Add(item, weight...)
}

// All get all items.
func All() interface{} {
	return DefaultBalancer.All()
}

// Select gets next selected item.
func Select(_ ...string) string {
	return DefaultBalancer.Select()
}

// Name load balancer name.
func Name() string {
	return DefaultBalancer.Name()
}

// Remove remove an item.
func Remove(item string, _ ...bool) bool {
	return DefaultBalancer.Remove(item)
}

// RemoveAll remove all items.
func RemoveAll() {
	DefaultBalancer.RemoveAll()
}

// Reset reset the balancer.
func Reset() {
	DefaultBalancer.Reset()
}

func Update(items interface{}) bool {
	return DefaultBalancer.Update(items)
}
