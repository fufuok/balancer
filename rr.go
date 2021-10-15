package balancer

import (
	"sync"
)

// RoundRobin
type rr struct {
	items   []string
	count   int
	current int

	sync.Mutex
}

func NewRoundRobin(items ...[]string) (lb *rr) {
	lb = &rr{}
	if len(items) > 0 && len(items[0]) > 0 {
		lb.Update(items[0])
	}
	return
}

func (b *rr) Add(item string, _ ...int) {
	b.Lock()
	b.items = append(b.items, item)
	b.count++
	b.Unlock()
}

func (b *rr) All() interface{} {
	all := make([]string, b.count)

	b.Lock()
	for i, v := range b.items {
		all[i] = v
	}
	b.Unlock()

	return all
}

func (b *rr) Name() string {
	return "RoundRobin"
}

func (b *rr) Select(_ ...string) (item string) {
	b.Lock()
	switch b.count {
	case 0:
		item = ""
	case 1:
		item = b.items[0]
	default:
		item = b.items[b.current]
		b.current = (b.current + 1) % b.count
	}
	b.Unlock()

	return
}

func (b *rr) Remove(item string, asClean ...bool) (ok bool) {
	b.Lock()
	defer b.Unlock()

	clean := len(asClean) > 0 && asClean[0]
	for i := 0; i < b.count; i++ {
		if item == b.items[i] {
			b.items = append(b.items[:i], b.items[i+1:]...)
			b.count--
			ok = true
			// remove all or remove one
			if !clean {
				return
			}
			i--
		}
	}
	return
}

func (b *rr) RemoveAll() {
	b.Lock()
	b.items = b.items[:0]
	b.count = 0
	b.current = 0
	b.Unlock()
}

func (b *rr) Reset() {
	b.Lock()
	b.current = 0
	b.Unlock()
}

func (b *rr) Update(items interface{}) bool {
	v, ok := items.([]string)
	if !ok {
		return false
	}

	b.Lock()
	b.items = v
	b.count = len(v)
	b.current = 0
	b.Unlock()

	return true
}
