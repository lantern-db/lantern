package model

type DocumentFrequency map[Key]uint32

func NewDocumentFrequency() DocumentFrequency {
	return make(map[Key]uint32)
}

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
