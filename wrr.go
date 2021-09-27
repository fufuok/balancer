package balancer

import (
	"sync"
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
	b.Lock()
	defer b.Unlock()

	w := 1
	if len(weight) > 0 {
		w = weight[0]
	}
	b.add(item, w)
	b.all[item] = w
}

func (b *wrr) add(item string, weight int) {
	b.remove(item)
	b.items = append(b.items, &wrrItems{
		item:   item,
		weight: weight,
	})
	b.n++

	if weight > 0 {
		if b.gcd == 0 {
			b.i = -1
			b.cw = 0
			b.gcd = weight
			b.max = weight
		} else {
			b.gcd = gcd(b.gcd, weight)
			if b.max < weight {
				b.max = weight
			}
		}
	}
}

func (b *wrr) All() interface{} {
	return b.all
}

func (b *wrr) Name() string {
	return "WeightedRoundRobin"
}

func (b *wrr) Select(_ ...string) string {
	switch b.n {
	case 0:
		return ""
	case 1:
		if b.items[0].weight > 0 {
			return b.items[0].item
		}
		return ""
	default:
		return b.chooseNext().item
	}
}

func (b *wrr) chooseNext() *wrrItems {
	b.Lock()
	defer b.Unlock()

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
	for i := 0; i < b.n; i++ {
		if item == b.items[i].item {
			b.items = append(b.items[:i], b.items[i+1:]...)
			b.n--
			delete(b.all, item)
			ok = true
			return
		}
	}
	return
}

func (b *wrr) RemoveAll() {
	b.Lock()
	defer b.Unlock()

	b.removeAll()
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
	defer b.Unlock()

	b.i = -1
	b.cw = 0
}

func (b *wrr) Update(items interface{}) bool {
	v, ok := items.(map[string]int)
	if !ok {
		return false
	}

	b.Lock()
	defer b.Unlock()

	b.removeAll()
	b.all = v

	for i, w := range v {
		b.add(i, w)
	}

	return true
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
