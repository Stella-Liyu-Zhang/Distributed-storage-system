package surfstore

import (
	context "context"
	"fmt"
	"sync"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var muMetaGet sync.Mutex
var muMetaUpdate sync.Mutex

type MetaStore struct {
	FileMetaMap    map[string]*FileMetaData
	BlockStoreAddr string
	Mutex          sync.Mutex
	UnimplementedMetaStoreServer
}

func (m *MetaStore) GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	return &FileInfoMap{FileInfoMap: m.FileMetaMap}, nil
}

func (m *MetaStore) UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	fileName := fileMetaData.Filename
	if meta, exists := m.FileMetaMap[fileName]; exists {
		if fileMetaData.Version-meta.Version < 0 {
			return &Version{
				Version: -1,
			}, fmt.Errorf("file version error: should be %d but %d", meta.Version+1, fileMetaData.Version)
		}
		(*meta).Version += 1
		(*meta).BlockHashList = fileMetaData.BlockHashList

	} else {
		m.FileMetaMap[fileName] = fileMetaData
	}

	return &Version{
		Version: m.FileMetaMap[fileName].Version,
	}, nil
}

func (m *MetaStore) GetBlockStoreAddr(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddr, error) {
	return &BlockStoreAddr{Addr: m.BlockStoreAddr}, nil
}

// This line guarantees all method for MetaStore are implemented
var _ MetaStoreInterface = new(MetaStore)

func NewMetaStore(blockStoreAddr string) *MetaStore {
	return &MetaStore{
		FileMetaMap:    map[string]*FileMetaData{},
		BlockStoreAddr: blockStoreAddr,
		Mutex:          sync.Mutex{},
	}
}
