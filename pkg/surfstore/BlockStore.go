package surfstore

import (
	context "context"
	sync "sync"
)

type BlockStore struct {
	BlockMap map[string]*Block
	Mtx      *sync.RWMutex

	UnimplementedBlockStoreServer
}

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {
	//panic("todo")
	bs.Mtx.RLock()
	defer bs.Mtx.RUnlock()
	block := bs.BlockMap[blockHash.Hash]
	return block, nil
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {
	//panic("todo")
	hash := GetBlockHashString(block.BlockData)

	bs.Mtx.Lock()
	defer bs.Mtx.Unlock()

	bs.BlockMap[hash] = block

	return &Success{Flag: true}, nil
}

// Given a list of hashes “in”, returns a list containing the
// subset of in that are stored in the key-value store
func (bs *BlockStore) HasBlocks(ctx context.Context, blockHashesIn *BlockHashes) (*BlockHashes, error) {
	//panic("todo")
	bs.Mtx.RLock()
	defer bs.Mtx.RUnlock()

	var hashout BlockHashes
	for _, hash := range blockHashesIn.Hashes {
		_, exists := bs.BlockMap[hash]
		if exists {
			hashout.Hashes = append(hashout.Hashes, hash)
		}
	}

	return &hashout, nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)

func NewBlockStore() *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
		Mtx:      &sync.RWMutex{},
	}
}
