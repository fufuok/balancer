package balancer

import (
	"sync"

	"github.com/fufuok/balancer/utils"
)

// Weighted Round-Robin Scheduling
// Ref: http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
type wrr struct {
	items []*wrrItems
	i     int
	n     int
	cw    int
	gcd   int
	max   int
	all   map[string]int

	sync.Mutex
}

type wrrItems struct {
	item   string
	weight int
}

func NewWeightedRoundRobin(items ...map[string]int) (lb *wrr) {
	if len(items) > 0 && len(items[0]) > 0 {
		lb = &wrr{}
		lb.Update(items[0])
		return
	}
	return &wrr{
		all: make(map[string]int),
	}
}

func (b *wrr) Add(item string, weight ...int) {
	w := 1
	if len(weight) > 0 {
		w = weight[0]
	}

	b.Lock()
	b.add(item, w)
	b.Unlock()
}

func (b *wrr) add(item string, weight int) {
	b.remove(item)

	b.items = append(b.items, &wrrItems{
		item:   item,
		weight: weight,
	})
	b.n++
	b.all[item] = weight

	b.addSettings(weight)
}

func (b *wrr) addSettings(weight int) {
	if weight > 0 {
		if b.gcd == 0 {
			b.i = -1
			b.cw = 0
			b.gcd = weight
			b.max = weight
		} else {
			b.gcd = utils.GCD(b.gcd, weight)
			if b.max < weight {
				b.max = weight
			}
		}
	}
}

func (b *wrr) All() interface{} {
	all := make(map[string]int)

	b.Lock()
	for k, v := range b.all {
		all[k] = v
	}
	b.Unlock()

	return all
}

func (b *wrr) Name() string {
	return "WeightedRoundRobin"
}

func (b *wrr) Select(_ ...string) (item string) {
	b.Lock()
	switch b.n {
	case 0:
		item = ""
	case 1:
		if b.items[0].weight > 0 {
			item = b.items[0].item
		}
	default:
		item = b.chooseNext().item
	}
	b.Unlock()

	return
}

func (b *wrr) chooseNext() *wrrItems {
	for {
		b.i = (b.i + 1) % b.n
		if b.i == 0 {
			b.cw = b.cw - b.gcd
			if b.cw <= 0 {
				b.cw = b.max
				if b.cw == 0 {
					return nil
				}
			}
		}

		if b.items[b.i].weight >= b.cw {
			return b.items[b.i]
		}
	}
}

func (b *wrr) Remove(item string, _ ...bool) bool {
	b.Lock()
	defer b.Unlock()

	return b.remove(item)
}

func (b *wrr) remove(item string) (ok bool) {
	maxWeight := 0
	for i := 0; i < b.n; i++ {
		if item == b.items[i].item {
			if b.max == b.items[i].weight {
				b.max = maxWeight
			}
			b.items = append(b.items[:i], b.items[i+1:]...)
			b.n--
			delete(b.all, item)
			ok = true
			return
		}
		if b.items[i].weight > maxWeight {
			maxWeight = b.items[i].weight
		}
	}
	return
}

func (b *wrr) RemoveAll() {
	b.Lock()
	b.removeAll()
	b.Unlock()
}

func (b *wrr) removeAll() {
	b.items = b.items[:0]
	b.n = 0
	b.i = -1
	b.cw = 0
	b.gcd = 0
	b.max = 0
	b.all = make(map[string]int)
}

func (b *wrr) Reset() {
	b.Lock()
	b.i = -1
	b.cw = 0
	b.Unlock()
}

func (b *wrr) Update(items interface{}) bool {
	v, ok := items.(map[string]int)
	if !ok {
		return false
	}

	b.Lock()
	defer b.Unlock()

	b.n = len(v)
	b.i = -1
	b.cw = 0
	b.gcd = 0
	b.max = 0
	b.all = v

	if cap(b.items) >= b.n {
		b.items = b.items[:b.n]
	} else {
		b.items = make([]*wrrItems, b.n)
	}

	i := 0
	for item, weight := range v {
		b.items[i] = &wrrItems{
			item:   item,
			weight: weight,
		}
		b.addSettings(weight)
		i++
	}

	return true
}
