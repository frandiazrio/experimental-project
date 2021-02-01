package chord

import (
	"testing"
)

func TestHashEquals(t *testing.T) {
	a := "A"
	b := "B"
	eq := isEqual([]byte(a), []byte(b))
	if eq {
		t.Errorf("Method is incorrect: got A and B are %t when it is false", eq)
	}
}

func TestInBetween(t *testing.T) {
	a := "a"
	b := "b"
	c := "c"

	bet := isBetween([]byte(b), []byte(a), []byte(c))

	if !bet {
		t.Error("hash of b should be between a and c ")
	}

	if isBetween([]byte(a), []byte(b), []byte(b)) {
		t.Errorf("hash a cannot be between two identical ends of a set")
	}

	if !isBetween([]byte(b), []byte(c), []byte(a)) {
		t.Errorf("hash of b should be between c and a (reversed)")
	}
}

func TestPowerOfTwo(t *testing.T) {
	num := 1
	res := true
	for num <= 1024 {
		res = res && isPowerOfTwo(num)
		num *= 2
	}

	if !res {
		t.Errorf("Error in power of two function")
	}
}
