package cache

import . "github.com/lantern-db/lantern/graph/model"

type DocumentFrequency map[Key]int

func (d DocumentFrequency) Increment(key Key) {
	if _, ok := d[key]; ok {
		d[key] += 1
	} else {
		d[key] = 1
	}
}

func (d DocumentFrequency) Decrement(key Key) {
	if _, ok := d[key]; ok {
		d[key] -= 1
		if d[key] == 0 {
			delete(d, key)
		}
	}
}

func (d DocumentFrequency) Get(key Key) (int, bool) {
	v, ok := d[key]
	return v, ok
}
