package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestPerm(t *testing.T) {
	res := rand.Perm(10)
	for i, v := range res {
		fmt.Printf("%2d: %10d\n", i, v)
	}
}
