package balancer

import (
	"sort"
	"sync"

	"github.com/fufuok/balancer/utils"
)

// WeightedRand
type wr struct {
	items   []*wrItem
	weights []int
	count   int
	max     uint32
	all     map[string]int

	sync.RWMutex
}

type wrItem struct {
	item   string
	weight int
}

func NewWeightedRand(items ...map[string]int) (lb *wr) {
	if len(items) > 0 && len(items[0]) > 0 {
		lb = &wr{}
		lb.Update(items[0])
		return
	}
	return &wr{
		all: make(map[string]int),
	}
}

func (b *wr) Add(item string, weight ...int) {
	w := 1
	if len(weight) > 0 {
		w = weight[0]
	}

	all := b.All().(map[string]int)
	all[item] = w
	b.Update(all)
}

func (b *wr) All() interface{} {
	all := make(map[string]int)

	b.RLock()
	for k, v := range b.all {
		all[k] = v
	}
	b.RUnlock()

	return all
}

func (b *wr) Name() string {
	return "WeightedRand"
}

func (b *wr) Select(_ ...string) (item string) {
	b.RLock()
	switch b.count {
	case 0:
		item = ""
	case 1:
		item = b.items[0].item
	default:
		r := utils.FastRandn(b.max) + 1
		i := utils.SearchInts(b.weights, int(r))
		item = b.items[i].item
	}
	b.RUnlock()

	return
}

func (b *wr) Remove(item string, _ ...bool) (ok bool) {
	b.RLock()
	_, ok = b.all[item]
	b.RUnlock()

	if ok {
		all := b.All().(map[string]int)
		delete(all, item)
		b.Update(all)
	}
	return
}

func (b *wr) RemoveAll() {
	b.Lock()
	b.items = b.items[:0]
	b.weights = b.weights[:0]
	b.count = 0
	b.max = 0
	b.all = make(map[string]int)
	b.Unlock()
}

func (b *wr) Reset() {}

func (b *wr) Update(items interface{}) bool {
	v, ok := items.(map[string]int)
	if !ok {
		return false
	}

	var (
		count int
		data  []*wrItem
	)
	for item, weight := range v {
		// Discard items with a weight less than 1
		if weight > 0 {
			data = append(data, &wrItem{
				item:   item,
				weight: weight,
			})
			count++
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].weight < data[j].weight
	})

	max := 0
	weights := make([]int, count)
	for i := range data {
		max += data[i].weight
		weights[i] = max
	}

	b.Lock()
	b.items = data
	b.weights = weights
	b.count = count
	b.max = uint32(max)
	b.all = v
	b.Unlock()

	return true
}
