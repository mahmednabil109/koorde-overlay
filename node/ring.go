package node

import "encoding/hex"

// Fixed Identifer space of 160 bits
//* [Note] this works for hex base ids
// the identifier is a []byte each entry contains 2 digits
//TODO make it works for generic id radixes
var (
	MAX_ID  ID = []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	ZERO_ID ID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

type ID []byte

// function checks if the id is in left execlusive range
// id is in the range (a, b]
func (id ID) InLXRange(a, b ID) bool {
	if a.Less(b) {
		return a.Less(id) && (id.Less(b) || id.Equal(b))
	} else {
		return id.InLXRange(a, MAX_ID) || id.Equal(ZERO_ID) || id.InLXRange(ZERO_ID, b)
	}
}

// checks that id is befor a in the ring
func (id ID) Less(a ID) bool {
	if len(id) == len(a) {
		for idx, digit := range id {
			if digit != a[idx] {
				return digit < a[idx]
			}
		}
	}
	return len(id) < len(a)
}

func (id ID) Equal(a ID) bool {
	for idx, digit := range id {
		if digit != a[idx] {
			return false
		}
	}
	return true
}

func (id ID) LeftShift() (ID, byte) {
	var (
		_id   []byte = make([]byte, len(id))
		carry byte
	)

	copy(_id, id)
	carry = ((_id[0] >> 0x04) & 0x0f)
	for i, digit := range _id {
		_id[i] = ((digit << 0x04) & 0xf0)
		if i != len(_id)-1 {
			_id[i] |= ((_id[i+1] >> 0x04) & 0x0f)
		}
	}
	return _id, carry
}

func (id ID) TopShift(a ID) ID {
	_id, _ := id.LeftShift()
	_id[len(_id)-1] |= ((a[0] >> 0x04) & 0x0f)
	return _id
}

func (id ID) String() string {
	return hex.EncodeToString(id)
}
