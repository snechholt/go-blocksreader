go-blocksreader
===============

Implements io.Reader and reads a set of blocks from an underlaying io.ReadSeeker as if the 
blocks were a continous stream of bytes.

Imagine the wrapped io.ReadSeeker is a piece of paper with text and the blocks are text 
highlighted with a marker pen. The blocks reader returns all the highlighted text (blocks)
as a single stream.
