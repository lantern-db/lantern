package model

type IncrementalStatMap map[Key]uint32

func (d IncrementalStatMap) Increment(key Key) {
	if _, ok := d[key]; ok {
		d[key] += 1
	} else {
		d[key] = 1
	}
}

func (d IncrementalStatMap) Decrement(key Key) {
	if _, ok := d[key]; ok {
		d[key] -= 1
		if d[key] == 0 {
			delete(d, key)
		}
	}
}
