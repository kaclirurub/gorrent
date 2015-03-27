package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/missingsix/bencode"
)

type torrent struct {
	announce string
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

// TODO : Probably easier to parse the whole thing then create the torrent struct...
func loadTorrent(dict map[string]interface{}) (torrent, error) {
	announce := dict["announce"].(string)

	decodedInfo := dict["info"].(map[string]interface{})
	//	fmt.Println(info)
	tempFiles := reflect.ValueOf(decodedInfo["files"].(interface{}))
	//	s := reflect.ValueOf(files)

	// TODO: Append files to the torrent.file struct
	// TODO: Need to get lengths from files as well
	files := make([]file, tempFiles.Len())
	for i := 0; i < tempFiles.Len(); i++ {
		length := tempFiles.Index(i).Interface().(map[string]interface{})["length"].(int64)
		fmt.Println("length->", length)
		tempPath := tempFiles.Index(i).Interface().(map[string]interface{})["path"].([]interface{})

		paths := make([]string, len(tempPath))
		for p := range tempPath {
			paths[p] = tempPath[p].(string)
		}

		for _, pathName := range paths {
			fmt.Println("pathName->", pathName)
		}

		file := file{length, paths}
		fmt.Println("file->", file)
		files[i] = file
	}

	fmt.Println("files->", files)

	name := decodedInfo["name"].(string)
	fmt.Println(name)
	//	fmt.Println("pieceLength->", reflect.TypeOf(info["piece length"]))
	pieceLength := decodedInfo["piece length"].(int64)
	fmt.Println(pieceLength)
	//torrent.info.files = files.(map[int]string)

	pieces := decodedInfo["pieces"].(string)

	torrent := torrent{announce, info{files, name, pieceLength, pieces}}

	fmt.Println("torrent->", torrent)

	return torrent, nil
}

func main() {
	//file, err := os.Open(os.Args[1])
	file, err := os.Open("/home/ctglover/Downloads/Simpsons.torrent")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dict, err := bencode.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	torrent, err := loadTorrent(dict)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n" + torrent.announce)

	/*for key, value := range dict {
		fmt.Printf(key+"->%v\n", value)
	}*/
	//fmt.Printf("string: %s\n", dict["string key"].(string))
	//fmt.Printf("string: %s\n", dict["int key"].(int64))
}
