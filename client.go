package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

func InitClient() {
	fmt.Println("File sharing platform for school and students.")
	fmt.Println("It is functioning as a client.")
	directory := AskAnswer("Directory to receive file: ")
	directory = strings.Trim(directory, "\\")
	directory += "\\"
	serverID := AskAnswer("Server ID: ")
	serverID, err := DecryptID(serverID)
	CheckError(err)
	serverAddress := fmt.Sprintf("%s.onion:10621", serverID)
	fmt.Println("Starting tor")
	t, err := StartClientTor()
	CheckError(err)
	defer t.Close()
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute * 3)
	defer dialCancel()
	dialer, err := t.Dialer(dialCtx, nil)
	CheckError(err)
	fmt.Println("Started tor successfully!")
	fmt.Println("Connecting to server...")
	conn, err := dialer.DialContext(dialCtx, "tcp", serverAddress)
	CheckError(err)
	fmt.Println("Successfully connected to server!")
	for {
		fileInfo, err := ReceiveInfo(conn)
		CheckError(err)
		fileName := strings.Split(fileInfo, "|||")[0]
		fileSize, err := strconv.Atoi(strings.Split(fileInfo, "|||")[1])
		CheckError(err)
		fmt.Println("File received! Name: ", fileName, " Size: ", fileSize)
		fileBuf, err := ReceiveFile(conn, fileSize)
		CheckError(err)
		err = SaveToFile(directory + fileName, fileBuf)
		CheckError(err)
		fmt.Println("File saved to ", directory + fileName)
	}
}

func ReceiveInfo(conn net.Conn) (string, error) {
	var resultBuf []byte
	for {
		tempBuf := make([]byte, 1)
		_, err := conn.Read(tempBuf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			return "", err
		}
		if string(tempBuf) == "\n" {
			break
		} else {
			resultBuf = append(resultBuf, tempBuf...)
		}
	}
	return string(resultBuf), nil
}

func ReceiveFile(conn net.Conn, size int) ([]byte, error) {
	fileBuf := make([]byte, size)
	_, err := conn.Read(fileBuf)
	if err != nil {
		return nil, err
	}
	return fileBuf, nil
}

func SaveToFile(path string, file []byte) error {
	err := ioutil.WriteFile(path, file, 0644)
	return err
}