package main

import (
	"encoding/csv"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/divyag9/protobuf/packages/protoTest"
	"github.com/golang/protobuf/proto"
)

func main() {
	c := make(chan *pb.MediaStore)
	go func() {
		for {
			message := <-c
			writeValuesToFile(message)
		}
	}()
	listener, err := net.Listen("tcp", "127.0.0.1:2110")
	checkError(err)
	for {
		if conn, err := listener.Accept(); err == nil {
			//If err is nil then that means that data is available for us so we take up this data and pass it to a new goroutine
			go handleProtoClient(conn, c)
		} else {
			continue
		}
	}
}

func writeValuesToFile(mediaStore *pb.MediaStore) {
	file, err := os.OpenFile("CSVValues.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	defer file.Close()
	writer := csv.NewWriter(file)
	medias := mediaStore.GetMedia()
	for _, media := range medias {
		client := media.Client
		username := media.Username
		version := strconv.Itoa(int(media.Version))

		mediaItems := media.GetMediaItem()
		for _, mediaItem := range mediaItems {
			record := []string{client, username, version, mediaItem.Guid, mediaItem.Mimetype, mediaItem.Ordernumber, mediaItem.Type.String(), string(mediaItem.ImageBytes)}
			writer.Write(record)
		}
	}
}

func handleProtoClient(conn net.Conn, c chan *pb.MediaStore) {
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	checkError(err)

	protodata := &pb.MediaStore{}
	err1 := proto.Unmarshal(data[0:n], protodata)
	checkError(err1)
	c <- protodata
}

func checkError(err error) {
	if err != nil {
		log.Fatal("FatalError: ", err)
	}
}
