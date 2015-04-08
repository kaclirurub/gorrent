package gorrent

import (
	"reflect"
)

type Torrent struct {
	Announce string
	info     info
}

type info struct {
	files       []file
	name        string
	pieceLength int64
	pieces      string // SHA1 hash of piece (20 bytes long)
}

type file struct {
	length int64
	path   []string
}

func (t Torrent) GetConnectionString() string {
	return t.Announce
}

func LoadTorrent(dict map[string]interface{}) (Torrent, error) {
	announce := dict["announce"].(string)

	decodedInfo := dict["info"].(map[string]interface{})
	tempFiles := reflect.ValueOf(decodedInfo["files"].(interface{}))

	files := make([]file, tempFiles.Len())
	for i := 0; i < tempFiles.Len(); i++ {
		length := tempFiles.Index(i).Interface().(map[string]interface{})["length"].(int64)
		tempPath := tempFiles.Index(i).Interface().(map[string]interface{})["path"].([]interface{})

		paths := make([]string, len(tempPath))
		for p := range tempPath {
			paths[p] = tempPath[p].(string)
		}

		file := file{length, paths}
		files[i] = file
	}

	name := decodedInfo["name"].(string)
	pieceLength := decodedInfo["piece length"].(int64)

	pieces := decodedInfo["pieces"].(string)

	torrent := Torrent{announce, info{files, name, pieceLength, pieces}}

	return torrent, nil
}
