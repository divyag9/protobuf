package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/divyag9/protobuf/packages/protoTest"
	"github.com/golang/protobuf/proto"
)

//Headers of the csv file
type Headers []string

func main() {
	fileName := "test.csv"
	dest := "127.0.0.1:2110"
	mediaStore := &pb.MediaStore{}
	media, err := retrieveDataFromFile(fileName)
	checkError(err)
	mediaStore.Media = append(mediaStore.Media, media)
	out, err := proto.Marshal(mediaStore)
	if err != nil {
		log.Fatalln("Failed to encode media store:", err)
	}
	sendDataToDest(out, dest)
}

func checkError(err error) {
	if err != nil {
		log.Fatal("FatalError: ", err)
	}
}

func retrieveDataFromFile(fname string) (*pb.Media, error) {
	file, err := os.Open(fname)
	checkError(err)
	defer file.Close()
	csvReader := csv.NewReader(file)

	var headers Headers
	headers, err = csvReader.Read()
	checkError(err)
	clientIndex := headers.getHeaderIndex("client")
	usernameIndex := headers.getHeaderIndex("username")
	versionIndex := headers.getHeaderIndex("version")
	guidIndex := headers.getHeaderIndex("guid")
	ordernumberIndex := headers.getHeaderIndex("ordernumber")
	loantypeIndex := headers.getHeaderIndex("loantype")
	mimetypeIndex := headers.getHeaderIndex("mimetype")
	imageBytesIndex := headers.getHeaderIndex("imagebytes")

	var media *pb.Media
	for {
		record, err := csvReader.Read()
		if err != io.EOF {
			checkError(err)
		} else {
			break
		}

		media = &pb.Media{}
		media.Client = record[clientIndex]
		media.Username = record[usernameIndex]
		versionValue, err := strconv.Atoi(record[versionIndex])
		media.Version = int64(versionValue)

		mediaItem := &pb.Media_MediaItem{}
		mediaItem.Guid = record[guidIndex]
		mediaItem.Ordernumber = record[ordernumberIndex]
		mediaItem.Mimetype = record[mimetypeIndex]
		mediaItem.ImageBytes = []byte(record[imageBytesIndex])
		//protoMedia.Client = record[clientIndex]
		switch record[loantypeIndex] {
		case "type1":
			mediaItem.Type = pb.Media_TYPE1
		case "type2":
			mediaItem.Type = pb.Media_TYPE2
		case "type3":
			mediaItem.Type = pb.Media_TYPE3
		case "type4":
			mediaItem.Type = pb.Media_TYPE4
		default:
			fmt.Printf("Unknown phone type %q.  Using default.\n", record[loantypeIndex])
		}

		media.MediaItem = append(media.MediaItem, mediaItem)
	}
	return media, nil
}

func sendDataToDest(data []byte, dst string) {
	conn, err := net.Dial("tcp", dst)
	checkError(err)
	n, err := conn.Write(data)
	checkError(err)
	fmt.Println("Sent " + strconv.Itoa(n) + " bytes")
}

func (h Headers) getHeaderIndex(headername string) int {
	for index, s := range h {
		if s == headername {
			return index
		}
	}
	return -1
}
