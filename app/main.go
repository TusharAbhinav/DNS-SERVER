package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"

	section "github.com/codecrafters-io/dns-server-starter-go/app/sections"
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/header"
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/question"
)

var _ = net.ListenUDP

func createHeaderSection(buf []byte) (header.Header, error) {
	buffer := bytes.NewReader(buf[:2])
	var ID, QDCOUNT uint16
	err := binary.Read(buffer, binary.BigEndian, &ID)
	QDCOUNT = 1
	headerObj := header.Header{
		ID:      ID,
		QR:      1 << 7,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: QDCOUNT,
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
func createQuestionSection() (question.Question, error) {
	var Name []byte
	domainName := "codecrafters.io"
	labels := bytes.Split([]byte(domainName), []byte("."))
	questionObj := question.Question{}

	for _, label := range labels {
		Name = append(Name, byte(len(label)))
		Name = append(Name, label...)
	}
	Name = append(Name, 0x00)
	typeBuff := new(bytes.Buffer)
	var Type, Class uint16
	Type = 1
	err := binary.Write(typeBuff, binary.BigEndian, Type)
	if err != nil {
		fmt.Println("Error encoding Type")
		return questionObj, err
	}
	Class = 1
	classBuff := new(bytes.Buffer)
	err = binary.Write(classBuff, binary.BigEndian, Class)
	if err != nil {
		fmt.Println("Error encoding Class")
		return questionObj, err
	}
	questionObj = question.Question{
		Name:  Name,
		Type:  typeBuff.Bytes(),
		Class: classBuff.Bytes(),
	}
	return questionObj, nil

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
		questionObj, err := createQuestionSection()
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
		sectionObj := section.Section{
			Header:   headerObj,
			Question: questionObj,
		}
		var sectionBuf bytes.Buffer
		encoder := gob.NewEncoder(&sectionBuf)
		err = encoder.Encode(sectionObj)
		if err != nil {
			fmt.Printf("Gob encoding failed: %v", err)
			break
		}
		response := sectionBuf.Bytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
