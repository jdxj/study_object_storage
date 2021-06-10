package reed_solomon

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

const (
	DataShards   = 4
	ParityShards = 2
	AllShards    = DataShards + ParityShards

	BlockPerShard = 8000
	BlockSize     = BlockPerShard * DataShards
)

func F() {
	sc, err := reedsolomon.NewStream(6, 4)
	sc.
}