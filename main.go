package main

import (
	//"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/missingsix/bencode"
	"github.com/missingsix/torrent"
)

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

	torrent, err := torrent.LoadTorrent(dict)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("torrent->", torrent)

	fmt.Println(torrent.GetConnectionString())
	ServerAddr, err := net.ResolveUDPAddr("udp", "open.demonii.com:1337")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Resolved!->", ServerAddr)
	}

	conn, err := net.DialUDP("udp", nil, ServerAddr)
	if err != nil {
		fmt.Println(err)
	}

	rand.Seed(time.Now().UnixNano())

	var connection_id int64 = 0x41727101980
	var action int32 = 0x0
	var transaction_id int32 = rand.Int31n(255)

	fmt.Println(connection_id)
	fmt.Println(action)
	fmt.Println(transaction_id)

	connectionRequest := new(bytes.Buffer)
	binary.Write(connectionRequest, binary.BigEndian, connection_id)
	binary.Write(connectionRequest, binary.BigEndian, action)
	binary.Write(connectionRequest, binary.BigEndian, transaction_id)

	fmt.Println("connectionRequest->", connectionRequest.Bytes())
	conn.Write(connectionRequest.Bytes())
	//	fmt.Fprintf(conn, buf.Bytes())

	connectResponse := make([]byte, 16)
	//status, err := bufio.NewReader(conn).Read(p)
	status, _, err := conn.ReadFromUDP(connectResponse)

	fmt.Println("status->", status)
	fmt.Println("connectResponse->", connectResponse)

	action = 0x1

	announceRequest := new(bytes.Buffer)

	binary.Write(announceRequest, binary.BigEndian, connectResponse[8:])
	binary.Write(announceRequest, binary.BigEndian, action)
	binary.Write(announceRequest, binary.BigEndian, transaction_id)

	fmt.Println("announceRequest->", announceRequest.Bytes())

}
