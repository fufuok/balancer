package balancer

import (
	"sync"
)

// Smooth weighted round-robin balancing
// Ref: https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
type swrr struct {
	items []*swrrItem
	count int
	all   map[string]int

	sync.Mutex
}

type swrrItem struct {
	item          string
	weight        int
	currentWeight int
}

func NewSmoothWeightedRoundRobin(items ...map[string]int) (lb *swrr) {
	if len(items) > 0 && len(items[0]) > 0 {
		lb = &swrr{}
		lb.Update(items[0])
		return
	}
	return &swrr{
		all: make(map[string]int),
	}
}

func (b *swrr) Add(item string, weight ...int) {
	w := 1
	if len(weight) > 0 {
		w = weight[0]
	}

	b.Lock()
	b.add(item, w)
	b.Unlock()
}

func (b *swrr) add(item string, weight int) {
	b.remove(item)
	b.items = append(b.items, &swrrItem{
		item:   item,
		weight: weight,
	})
	b.count++
	b.all[item] = weight
}

func (b *swrr) All() interface{} {
	all := make(map[string]int)

	b.Lock()
	for k, v := range b.all {
		all[k] = v
	}
	b.Unlock()

	return all
}

func (b *swrr) Name() string {
	return "SmoothWeightedRoundRobin"
}

func (b *swrr) Select(_ ...string) (item string) {
	b.Lock()
	switch b.count {
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

func (b *swrr) chooseNext() (choice *swrrItem) {
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
	b.items = b.items[:0]
	b.count = 0
	b.all = make(map[string]int)
	b.Unlock()
}

func (b *swrr) Reset() {
	b.Lock()
	for i := range b.items {
		b.items[i].currentWeight = b.items[i].weight
	}
	b.Unlock()
}

func (b *swrr) Update(items interface{}) bool {
	v, ok := items.(map[string]int)
	if !ok {
		return false
	}

	count := len(v)
	data := make([]*swrrItem, count)
	i := 0
	for item, weight := range v {
		data[i] = &swrrItem{
			item:   item,
			weight: weight,
		}
		i++
	}

	b.Lock()
	b.count = count
	b.all = v
	b.items = data
	b.Unlock()

	return true
}
