package balancer

import (
	"hash/fnv"
	"sync"

	"github.com/fufuok/balancer/internal/utils"

	"github.com/fufuok/balancer/internal/doublejump"
)

// JumpConsistentHash
type consistentHash struct {
	items []string
	count int
	h     *doublejump.Hash

	sync.RWMutex
}

func NewConsistentHash(items ...[]string) (lb *consistentHash) {
	if len(items) == 0 {
		return &consistentHash{
			h: doublejump.NewHash(),
		}
	}
	lb = &consistentHash{}
	lb.Update(items[0])
	return
}

func (b *consistentHash) Add(item string, _ int) {
	b.Lock()
	defer b.Unlock()

	b.add(item)
}

func (b *consistentHash) add(item string) {
	b.items = append(b.items, item)
	b.h.Add(item)
	b.count++
}

func (b *consistentHash) All() interface{} {
	return b.items
}

func (b *consistentHash) Name() string {
	return "ConsistentHash"
}

func (b *consistentHash) Select(key ...string) string {
	switch b.count {
	case 0:
		return ""
	case 1:
		return b.items[0]
	default:
		return b.chooseNext(key)
	}
}

func (b *consistentHash) chooseNext(key []string) (choice string) {
	b.RLock()
	defer b.RUnlock()

	hash := b.hash(utils.AddString(key...))
	choice, _ = b.h.Get(hash).(string)
	return
}

func (b *consistentHash) hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
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
	defer b.Unlock()

	b.removeAll()
}

func (b *consistentHash) removeAll() {
	b.items = b.items[:0]
	b.count = 0
	b.h = doublejump.NewHash()
}

func (b *consistentHash) Reset() {
}

func (b *consistentHash) Update(items interface{}) bool {
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