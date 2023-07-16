package surfstore

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
)

// Implement the logic for a client syncing with the server here.
func ClientSync(client RPCClient) {
	fmt.Printf("[Client %s] start func ClientSync\n", client.BaseDir)

	// lookup the index.txt file. if not, create one
	indexFilePath := ConcatPath(client.BaseDir, "/index.txt")
	if _, err := os.Stat(indexFilePath); errors.Is(err, os.ErrNotExist) {
		indexFile, _ := os.Create(indexFilePath)
		indexFile.Close()
	}
	// read the directory and get file info of the client into the map
	files, err := ioutil.ReadDir(client.BaseDir)
	if err != nil {
		log.Fatalf("error while reading directory : %v", err)
	}
	/*
		Read the index.txt to get the current client matadata info
	*/

	//Load the index.txt file
	metaDataMap, err := LoadMetaFromMetaFile(client.BaseDir)
	if err != nil {
		log.Fatalf("error while reading index.txt, %v", err)
	}
	remoteFileMetaMap := make(map[string]*FileMetaData)
	err = client.GetFileInfoMap(&remoteFileMetaMap)
	if err != nil {
		log.Fatalf("error while getting server info map, %v", err)
	}

	fileMap := make(map[string][]string)
	//iterate through all the file in the map and compare
	for _, file := range files {
		if file.Name() == "index.txt" {
			continue
		}

		var num int = int(math.Ceil(float64(file.Size()) / float64(client.BlockSize)))
		readFile, err := os.Open(client.BaseDir + "/" + file.Name())
		if err != nil {
			log.Println("error in opening file", err)
		}
		//Read bytes from file in base directory
		for i := 0; i < num; i++ {
			slice := make([]byte, client.BlockSize)
			len, err := readFile.Read(slice)
			if err != nil {
				log.Println(err)
			}
			slice = slice[:len]
			hash := GetBlockHashString(slice)
			fileMap[file.Name()] = append(fileMap[file.Name()], hash)
		}

		if val, exists := metaDataMap[file.Name()]; exists {
			if !reflect.DeepEqual(fileMap[file.Name()], val.BlockHashList) {
				metaDataMap[file.Name()].BlockHashList = fileMap[file.Name()]
				metaDataMap[file.Name()].Version++
			}
		} else {
			// Make a new file of version 1
			meta := FileMetaData{
				Filename:      file.Name(),
				Version:       1,
				BlockHashList: fileMap[file.Name()]}
			metaDataMap[file.Name()] = &meta
		}
	}

	// iterate through all the indexfile and compare
	for fileName, metaData := range metaDataMap {
		if _, exists := fileMap[fileName]; !exists {
			if len(metaData.BlockHashList) != 1 || metaData.BlockHashList[0] != "0" {
				metaData.Version++
				metaData.BlockHashList = []string{"0"}
			} else {
				fmt.Print("The file is unchanged")
			}
		}
		fmt.Printf("[Client %s] finished syncing from local\n", client.BaseDir)
	}

	var address string
	if err := client.GetBlockStoreAddr(&address); err != nil {
		log.Println(err)
	}

	remoteIndex := make(map[string]*FileMetaData)
	if err := client.GetFileInfoMap(&remoteIndex); err != nil {
		log.Println(err)
	}

	for file, meta := range metaDataMap {
		if remoteMetaData, exists := remoteIndex[file]; exists {
			if meta.Version > remoteMetaData.Version {
				upload(client, meta, address)
			}
		} else {
			upload(client, meta, address)
		}
	}

	//Check for updates on server, download
	for file, meta := range remoteIndex {
		if localMetaData, exists := metaDataMap[file]; exists {
			if localMetaData.Version < meta.Version {
				download(client, localMetaData, meta, address)
			} else if localMetaData.Version == meta.Version && !reflect.DeepEqual(localMetaData.BlockHashList, meta.BlockHashList) {
				download(client, localMetaData, meta, address)
			}
		} else {
			metaDataMap[file] = &FileMetaData{}
			localMetaData := metaDataMap[file]
			download(client, localMetaData, meta, address)
		}
	}

	WriteMetaFile(metaDataMap, client.BaseDir)
}

func download(client RPCClient, local *FileMetaData, remote *FileMetaData, address string) error {
	filePath := client.BaseDir + "/" + remote.Filename
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	*local = *remote

	if len(remote.BlockHashList) == 1 && remote.BlockHashList[0] == "0" {
		if err := os.Remove(filePath); err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	temp := ""
	for _, hash := range remote.BlockHashList {
		var b Block
		if err := client.GetBlock(hash, address, &b); err != nil {
			log.Println(err)
		}

		temp += string(b.BlockData)
	}
	file.WriteString(temp)

	return nil
}

func upload(client RPCClient, meta *FileMetaData, address string) error {
	fmt.Printf("[Client %s] start func upload\n", client.BaseDir)

	filePath := client.BaseDir + "/" + meta.Filename
	var mostRecentVer int32
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err = client.UpdateFile(meta, &mostRecentVer)
		if err != nil {
			log.Println(err)
		}
		meta.Version = mostRecentVer
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	fileStat, _ := os.Stat(filePath)
	var n int = int(math.Ceil(float64(fileStat.Size()) / float64(client.BlockSize)))
	for j := 0; j < n; j++ {
		slice := make([]byte, client.BlockSize)
		len, err := file.Read(slice)
		if err != nil && err != io.EOF {
			log.Println(err)
		}
		slice = slice[:len]

		block := Block{BlockData: slice, BlockSize: int32(len)}

		var succ bool
		if err := client.PutBlock(&block, address, &succ); err != nil {
			log.Println(err)
		}
	}

	if err := client.UpdateFile(meta, &mostRecentVer); err != nil {
		log.Println(err)
		meta.Version = -1
	}
	meta.Version = mostRecentVer

	return nil
}
