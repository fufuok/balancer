package balancer

import (
	"sync"

	"github.com/fufuok/balancer/utils"
)

type random struct {
	items []string
	count uint32

	sync.RWMutex
}

func NewRandom(items ...[]string) (lb *random) {
	lb = &random{}
	if len(items) > 0 && len(items[0]) > 0 {
		lb.Update(items[0])
	}
	return
}

func (b *random) Add(item string, _ ...int) {
	b.Lock()
	b.items = append(b.items, item)
	b.count++
	b.Unlock()
}

func (b *random) All() interface{} {
	all := make([]string, b.count)

	b.Lock()
	for i, v := range b.items {
		all[i] = v
	}
	b.Unlock()

	return all
}

func (b *random) Name() string {
	return "Random"
}

func (b *random) Select(_ ...string) (item string) {
	b.RLock()
	switch b.count {
	case 0:
		item = ""
	case 1:
		item = b.items[0]
	default:
		item = b.items[utils.FastRandn(b.count)]
	}
	b.RUnlock()

	return
}

func (b *random) Remove(item string, asClean ...bool) (ok bool) {
	b.Lock()
	defer b.Unlock()

	clean := len(asClean) > 0 && asClean[0]
	for i := uint32(0); i < b.count; i++ {
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

func (b *random) RemoveAll() {
	b.Lock()
	b.items = b.items[:0]
	b.count = 0
	b.Unlock()
}

func (b *random) Reset() {}

func (b *random) Update(items interface{}) bool {
	v, ok := items.([]string)
	if !ok {
		return false
	}

	b.Lock()
	b.items = v
	b.count = uint32(len(v))
	b.Unlock()

	return true
}
