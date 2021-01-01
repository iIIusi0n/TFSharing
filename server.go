package main

import (
	"context"
	"fmt"
	"github.com/cretz/bine/tor"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
)

type File struct {
	fileName string
	fileSize int64
	content []byte
}

var (
	channels []chan File
)

func InitServer() {
	fmt.Println("File sharing platform for school and students.")
	fmt.Println("It is functioning as a server.")
	fmt.Println("Starting tor")
	t, err := StartServerTor()
	CheckError(err)
	defer t.Close()
	fmt.Println("Started tor successfully!")
	fmt.Println("Starting server & Registering onion router")
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()
	onion, err := t.Listen(listenCtx, &tor.ListenConf{RemotePorts: []int{10621}, Version3: true})
	CheckError(err)
	defer onion.Close()
	fmt.Println("Started server successfully!")
	serverID, err := EncryptID(onion.ID)
	CheckError(err)
	fmt.Println("Your ID: ", serverID)
	fmt.Println("Start listening...")
	go Listener(onion)
	for {
		fileName := AskAnswer("File: ")
		fileStat, err := os.Stat(fileName)
		CheckError(err)
		fileSize := fileStat.Size()
		content, err := ioutil.ReadFile(fileName)
		CheckError(err)
		file := File{
			fileName: fileName,
			fileSize: fileSize,
			content:  content,
		}
		for i := 0; i < len(channels); i++ {
			channels[i] <- file
		}
	}
}

func Listener(onion *tor.OnionService) {
	for {
		conn, err := onion.Accept()
		CheckError(err)
		channel := make(chan File)
		channels = append(channels, channel)
		go ConnHandler(conn, channel)
	}
}

func ConnHandler(conn net.Conn, channel chan File) {
	for {
		file := <- channel
		fileInfo := file.fileName + "|||" + strconv.FormatInt(file.fileSize, 10) + "\n"
		_, err := conn.Write([]byte(fileInfo))
		CheckError(err)
		_, err = conn.Write(file.content)
		CheckError(err)
	}
}



