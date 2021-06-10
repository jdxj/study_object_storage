package reed_solomon

import (
	"fmt"
	"testing"

	"github.com/klauspost/reedsolomon"
)

func TestReedSolomon(t *testing.T) {
	encoder, err := reedsolomon.New(4, 2)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	data := []byte("123")

	shards, err := encoder.Split(data)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	err = encoder.Encode(shards)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("len1: %d\n", len(shards))
	for _, v := range shards {
		fmt.Printf("len2: %d\n", len(v))
	}
}
