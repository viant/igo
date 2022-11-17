package exec

import "math/bits"

type Uint64s []uint64

func (o Uint64s) Matches(u Uint64s) bool {
	switch len(u) {
	case 0:
		return false
	case 1:
		return o[0]&u[0] != 0
	case 2:
		return o[0]&u[0] != 0 || o[1]&u[1] != 0
	case 3:
		return o[0]&u[0] != 0 || o[1]&u[1] != 0 || o[2]&u[2] != 0
	}
	for i, w := range u {
		if o[i]&w != 0 {
			return true
		}
	}
	return false
}

//HasBit returns true if a bit at position in set
func (o Uint64s) HasBit(pos uint16) bool {
	idx := index(int(pos))
	result := (o[idx] & (1 << (pos % 64))) != 0
	return result
}

//ClearBit clears bit at position in set
func (o Uint64s) ClearBit(pos int) {
	o[index(pos)] &= ^(1 << (pos % 64))
}

//Count return bit population count
func (o Uint64s) Count() int {
	result := 0
	for _, w := range o {
		result += bits.OnesCount(uint(w))
	}
	return result
}

//Or applies logical OR
func (o Uint64s) Or(set Uint64s) Uint64s {
	var result = make([]uint64, len(set))
	for i, w := range set {
		result[i] = o[i] | w
	}
	return result
}

//SetBits sets bit at position in set
func (o *Uint64s) SetBits(positions []uint16) {
	for _, pos := range positions {
		o.SetBit(int(pos))
	}
}

//SetBit sets bit at position in set
func (o *Uint64s) SetBit(pos int) {
	idx := index(pos)
	for i := len(*o); i <= idx; i++ {
		*o = append(*o, 0)
	}
	(*o)[idx] |= (1 << (pos % 64))
}

func index(pos int) int {
	return pos / 64
}
