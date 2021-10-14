package balancer

import (
	"sync"

	"github.com/fufuok/balancer/utils"
)

type random struct {
	items []string
	count uint32

	sync.Mutex
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
	defer b.Unlock()

	b.add(item)
}

func (b *random) add(item string) {
	b.items = append(b.items, item)
	b.count++
}

func (b *random) All() interface{} {
	return b.items
}

func (b *random) Name() string {
	return "Random"
}

func (b *random) Select(_ ...string) string {
	switch b.count {
	case 0:
		return ""
	case 1:
		return b.items[0]
	default:
		return b.chooseNext()
	}
}

func (b *random) chooseNext() string {
	b.Lock()
	defer b.Unlock()

	return b.items[utils.FastRandn(b.count)]
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
	defer b.Unlock()

	b.removeAll()
}

func (b *random) removeAll() {
	b.items = b.items[:0]
	b.count = 0
}

func (b *random) Reset() {}

func (b *random) Update(items interface{}) bool {
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
