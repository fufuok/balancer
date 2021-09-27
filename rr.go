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
	if len(items) > 0 {
		lb.Update(items[0])
	}
	return
}

func (b *rr) Add(item string, _ int) {
	b.Lock()
	defer b.Unlock()

	b.add(item)
}

func (b *rr) add(item string) {
	b.items = append(b.items, item)
	b.count++
}

func (b *rr) All() interface{} {
	return b.items
}

func (b *rr) Name() string {
	return "RoundRobin"
}

func (b *rr) Select(_ ...string) string {
	switch b.count {
	case 0:
		return ""
	case 1:
		return b.items[0]
	default:
		return b.chooseNext()
	}
}

func (b *rr) chooseNext() (choice string) {
	b.Lock()
	defer b.Unlock()

	choice = b.items[b.current]
	b.current = (b.current + 1) % b.count
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
	defer b.Unlock()

	b.removeAll()
}

func (b *rr) removeAll() {
	b.items = b.items[:0]
	b.count = 0
	b.current = 0
}

func (b *rr) Reset() {
	b.current = 0
}

func (b *rr) Update(items interface{}) bool {
	v, ok := items.([]string)
	if !ok {
		return false
	}

	b.Lock()
	defer b.Unlock()

	b.removeAll()

	for _, x := range v {
		b.add(x)
	}

	return true
}
