package balancer

import (
	"sync"

	"github.com/fufuok/balancer/internal/doublejump"
	"github.com/fufuok/balancer/utils"
)

// JumpConsistentHash
type consistentHash struct {
	items []string
	count int
	h     *doublejump.Hash

	sync.RWMutex
}

func NewConsistentHash(items ...[]string) (lb *consistentHash) {
	if len(items) > 0 && len(items[0]) > 0 {
		lb = &consistentHash{}
		lb.Update(items[0])
		return
	}
	return &consistentHash{
		h: doublejump.NewHash(),
	}
}

func (b *consistentHash) Add(item string, _ ...int) {
	b.Lock()
	b.items = append(b.items, item)
	b.h.Add(item)
	b.count++
	b.Unlock()
}

func (b *consistentHash) All() interface{} {
	all := make([]string, b.count)

	b.Lock()
	for i, v := range b.items {
		all[i] = v
	}
	b.Unlock()

	return all
}

func (b *consistentHash) Name() string {
	return "ConsistentHash"
}

func (b *consistentHash) Select(key ...string) (item string) {
	b.RLock()
	switch b.count {
	case 0:
		item = ""
	case 1:
		item = b.items[0]
	default:
		hash := utils.HashString(key...)
		item, _ = b.h.Get(hash).(string)
	}
	b.RUnlock()

	return
}

func (b *consistentHash) Remove(item string, asClean ...bool) (ok bool) {
	b.Lock()
	defer b.Unlock()

	clean := len(asClean) > 0 && asClean[0]
	for i := 0; i < b.count; i++ {
		if item == b.items[i] {
			b.items = append(b.items[:i], b.items[i+1:]...)
			b.count--
			b.h.Remove(item)
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

func (b *consistentHash) RemoveAll() {
	b.Lock()
	b.items = b.items[:0]
	b.count = 0
	b.h = doublejump.NewHash()
	b.Unlock()
}

func (b *consistentHash) Reset() {}

func (b *consistentHash) Update(items interface{}) bool {
	v, ok := items.([]string)
	if !ok {
		return false
	}

	h := doublejump.NewHash()
	for _, x := range v {
		h.Add(x)
	}

	b.Lock()
	b.count = len(v)
	b.items = v
	b.h = h
	b.Unlock()

	return true
}
