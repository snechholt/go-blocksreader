package blocksreader

import (
	"errors"
	"io"
)

type Block struct {
	Offset int64
	Length int
}

func NewReader(r io.ReadSeeker, blocks []Block) *Reader {
	return &Reader{
		r:      r,
		blocks: blocks,
	}
}

// Reader implements io.Reader and reads a set of blocks from an underlaying io.ReadSeeker
// as if the blocks were a continous stream of bytes.
type Reader struct {
	r          io.ReadSeeker // The wrapped reader
	blocks     []Block       // The blocks to read from
	blockIndex int           // Index of the current block
	blockPos   int           // Position within the current block
}

func (this *Reader) Read(p []byte) (int, error) {
	if this.blockIndex >= len(this.blocks) {
		return 0, io.EOF
	}
	block := this.blocks[this.blockIndex]
	if block.Offset < 0 {
		return 0, errors.New("blockReader: invalid block - offset < 0")
	}
	if block.Length < 1 {
		return 0, errors.New("blockReader: invalid block - length < 1")
	}
	if this.blockPos >= block.Length {
		this.blockPos = 0
		this.blockIndex += 1
		return this.Read(p)
	}
	offset := block.Offset + int64(this.blockPos)
	if _, err := this.r.Seek(offset, 0); err != nil {
		return 0, err
	}
	n := block.Length - int(this.blockPos)
	if len(p) < n {
		n = len(p)
	}
	b := make([]byte, n)
	n, err := this.r.Read(b)
	copy(p, b)
	this.blockPos += n
	return n, err
}
