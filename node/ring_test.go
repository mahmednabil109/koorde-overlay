package node

import (
	"encoding/hex"
	"testing"

	"github.com/mahmednabil109/koorde-overlay/utils"
)

var (
	id = ID(utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582f29"))
	a  = ID(utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582f20"))
	b  = ID(utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582fff"))
	_b = ID(utils.ParseID("00000000000000000001"))
)

func Test_inlxrange(t *testing.T) {

	if !id.InLXRange(a, b) {
		t.Fatalf("id in the range but the function returns false")
	}

	if a.InLXRange(id, b) {
		t.Fatal("id is not in the range but the function returns true")
	}

	if !id.InLXRange(a, _b) {
		t.Fatal("id is in the range but the function return false")
	}

	if !id.InLXRange(ZERO_ID, MAX_ID) {
		t.Fatal("id in the range but the function returns false")
	}

	if ZERO_ID.InLXRange(ZERO_ID, MAX_ID) {
		t.Fatal("id is not in the range but the function returns true")
	}

	if !ZERO_ID.InLXRange(id, ZERO_ID) {
		t.Fatal("id is in range but the function returns false")
	}

	if ZERO_ID.InLXRange(ZERO_ID, _b) {
		t.Fatal("id is not in the range but the function returns true")
	}

	if !ZERO_ID.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !a.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !b.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !_b.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}
}

func Test_equal(t *testing.T) {
	if !a.Equal(a) {
		t.Fatal("equal returns false when it should return true")
	}

	if a.Equal(b) {
		t.Fatal("equal returns true when it should return false")
	}

	if ZERO_ID.Equal(_b) {
		t.Fatal("equal returns true when it should return false")

	}
}

func Test_leftshift(t *testing.T) {
	_id, _ := id.LeftShift()
	_id_str := hex.EncodeToString(_id)
	if _id_str != "09b0ce42948043810a1f2cc7e7079aec7582f290" {
		t.Fatal("leftshift dose not work correctly for hex base")
	}

	copy(_id, id)
	for i := 0; i < 40; i++ {
		_id, _ = _id.LeftShift()
	}
	t.Log(_id)
	if !_id.Equal(ZERO_ID) {
		t.Fatal("40 shift must result in ZERO_ID")
	}
}

func Test_topshift(t *testing.T) {
	_id := id.TopShift(id)
	_id_str := hex.EncodeToString(_id)
	if _id_str != "09b0ce42948043810a1f2cc7e7079aec7582f29a" {
		t.Fatal("topShift must replace lower digit by the top digit of the other id")
	}

	_a := make([]byte, len(a))

	copy(_id, id)
	copy(_a, a)
	for i := 0; i < 40; i++ {
		_id = _id.TopShift(_a)
		_a, _ = ID(_a).LeftShift()
	}

	if !_id.Equal(a) {
		t.Fatalf("after 40 topShifts the id must equal to the other one %v", _id)
	}
}
