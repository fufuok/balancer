package balancer

import (
	"sync"
)

// smooth weighted round-robin balancing
// Ref: https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
type swrr struct {
	items []*swrrItems
	count int
	all   map[string]int

	sync.Mutex
}

type swrrItems struct {
	item          string
	weight        int
	currentWeight int
}

func NewSmoothWeightedRoundRobin(items ...map[string]int) (lb *swrr) {
	if len(items) == 0 {
		return &swrr{
			all: make(map[string]int),
		}
	}
	lb = &swrr{}
	lb.Update(items[0])
	return
}

func (b *swrr) Add(item string, weight int) {
	b.Lock()
	defer b.Unlock()

	b.add(item, weight)
	b.all[item] = weight
}

func (b *swrr) add(item string, weight int) {
	b.remove(item)
	b.items = append(b.items, &swrrItems{
		item:   item,
		weight: weight,
	})
	b.count++
}

func (b *swrr) All() interface{} {
	return b.all
}

func (b *swrr) Name() string {
	return "SmoothWeightedRoundRobin"
}

func (b *swrr) Select(_ ...string) string {
	switch b.count {
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

func (b *swrr) chooseNext() (choice *swrrItems) {
	b.Lock()
	defer b.Unlock()

	total := 0
	for i := range b.items {
		c := b.items[i]
		if c == nil {
			return nil
		}

		total += c.weight
		c.currentWeight += c.weight

		if choice == nil || c.currentWeight > choice.currentWeight {
			choice = c
		}
	}

	if choice == nil {
		return nil
	}

	choice.currentWeight -= total

	return choice
}

func (b *swrr) Remove(item string, _ ...bool) bool {
	b.Lock()
	defer b.Unlock()

	return b.remove(item)
}

func (b *swrr) remove(item string) (ok bool) {
	for i := 0; i < b.count; i++ {
		if item == b.items[i].item {
			b.items = append(b.items[:i], b.items[i+1:]...)
			b.count--
			delete(b.all, item)
			ok = true
			return
		}
	}
	return
}

func (b *swrr) RemoveAll() {
	b.Lock()
	defer b.Unlock()

	b.removeAll()
}

func (b *swrr) removeAll() {
	b.items = b.items[:0]
	b.count = 0
	b.all = make(map[string]int)
}

func (b *swrr) Reset() {
	b.Lock()
	defer b.Unlock()

	for i := range b.items {
		b.items[i].currentWeight = b.items[i].weight
	}
}

func (b *swrr) Update(items interface{}) bool {
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
