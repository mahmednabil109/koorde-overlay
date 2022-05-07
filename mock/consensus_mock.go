package mock

import "sync"

type Block struct{ Info string }
type Consensus struct {
	data     []Block
	data_mux sync.Mutex
}

func (c *Consensus) AddBlock(b Block) {
	c.data_mux.Lock()
	defer c.data_mux.Unlock()

	c.data = append(c.data, b)
}

func (c *Consensus) GetBlocks() []Block {
	c.data_mux.Lock()
	defer c.data_mux.Unlock()

	return c.data
}
