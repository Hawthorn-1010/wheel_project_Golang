package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 因为测试，自己写hash，直接用转int
	m := New(func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	}, 3)

	m.Add("2", "4", "6")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	m.Add("8")
	testCases["27"] = "8"

	for k, v := range testCases {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
