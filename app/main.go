package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/sections/header"
)

var _ = net.ListenUDP

func createHeaderSection(buf []byte) (header.Header, error) {
	buffer := bytes.NewReader(buf[:2])
	var ID uint16
	err := binary.Read(buffer, binary.BigEndian, &ID)
	headerObj := header.Header{
		ID:      ID,
		QR:      1,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: 0,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
	if err != nil {
		fmt.Println("Error reading data", err)
		return headerObj, err
	}

	return headerObj, nil
}
func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
		headerObj, err := createHeaderSection(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.BigEndian, headerObj)
		if err != nil {
			fmt.Printf("binary.Write failed: %v", err)
			return
		}
		data := buf.Bytes()

		response := data

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
