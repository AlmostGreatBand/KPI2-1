package main

import "testing"

func TestPlus(t *testing.T) {
	res := plus(2, 3)

	if res != 5 {
		t.Fatalf("[error] expected 5 but was %d", res)
	}
}
