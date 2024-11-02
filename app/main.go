package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/answer"
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/header"
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/question"
	"net"
)

// |  Bit  |  7  |  6  |  5  |  4  |  3  |  2  |  1  |  0  |
// |-------|-----|-----|-----|-----|-----|-----|-----|-----|
// | Value |  0  |  0  |  0  |  0  | Opcode  |  AA  |  TC  |  RD  |

func createHeaderSection(buf []byte, ANCount uint16) ([]byte, error) {
	buffer := bytes.NewReader(buf[:2])
	var ID, QDCOUNT uint16

	idErr := binary.Read(buffer, binary.BigEndian, &ID)
	if idErr != nil {
		return nil, idErr
	}

	QDCOUNT = binary.BigEndian.Uint16(buf[4:6])

	headerObj := header.Header{
		ID:      ID,
		QR:      1,       // Set to 1 for a response
		OPCODE:  0,       // Standard query
		AA:      0,       // Not authoritative
		TC:      0,       // Message not truncated
		RD:      1,       // Recursion desired
		RA:      0,       // Recursion not available
		Z:       0,       // Reserved
		RCODE:   0,       // No error
		QDCOUNT: QDCOUNT, // Number of questions
		ANCOUNT: ANCount, // Number of answers
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	headerResponse := make([]byte, 0)

	// ID and Flags
	idBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(idBytes, headerObj.ID)
	headerResponse = append(headerResponse, idBytes...)

	flags := headerObj.QR<<15 | headerObj.OPCODE<<11 | headerObj.AA<<10 |
		headerObj.TC<<9 | headerObj.RD<<8 | headerObj.RA<<7 |
		headerObj.Z<<4 | headerObj.RCODE
	flagBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(flagBytes, flags)
	headerResponse = append(headerResponse, flagBytes...)

	// QDCOUNT, ANCOUNT, NSCOUNT, ARCOUNT
	countBytes := make([]byte, 8)
	binary.BigEndian.PutUint16(countBytes[0:2], headerObj.QDCOUNT)
	binary.BigEndian.PutUint16(countBytes[2:4], headerObj.ANCOUNT)
	binary.BigEndian.PutUint16(countBytes[4:6], headerObj.NSCOUNT)
	binary.BigEndian.PutUint16(countBytes[6:8], headerObj.ARCOUNT)
	headerResponse = append(headerResponse, countBytes...)

	return headerResponse, nil
}

func createLabel(off int, buf []byte) []byte {
	domainName := buf[off:]
	offset := 0
	var Name []byte

	for {
		length := int(domainName[offset])

		// Check if the label contains a pointer
		if (length & 0b11000000) == 0b11000000 {
			pointerOffset := int(domainName[offset]&0b00111111)<<8 | int(domainName[offset+1])
			Name = append(Name, createLabel(pointerOffset, buf)...)
			break
		}

		if length == 0 {
			Name = append(Name, 0x00)
			break
		}

		offset++
		Name = append(Name, byte(length))
		Name = append(Name, domainName[offset:offset+length]...)
		offset += length
	}
	return Name
}

func createQuestionSection(buf []byte) []byte {
	name := createLabel(0, buf)
	typeBuff := make([]byte, 2)
	classBuff := make([]byte, 2)

	binary.BigEndian.PutUint16(typeBuff, 1)
	binary.BigEndian.PutUint16(classBuff, 1)

	questionObj := question.Question{
		Name:  name,
		Type:  typeBuff,
		Class: classBuff,
	}
	questionResponse := []byte{}
	questionResponse = append(questionResponse, questionObj.Name...)
	questionResponse = append(questionResponse, questionObj.Type...)
	questionResponse = append(questionResponse, questionObj.Class...)

	return questionResponse
}

func createAnswerSection(buf []byte, answerCount int) ([]byte, error) {
	answerResponse := make([]byte, 0)

	for i := 0; i < answerCount; i++ {
		name := createLabel(0, buf)

		typeBuff := make([]byte, 2)
		binary.BigEndian.PutUint16(typeBuff, 1)

		classBuff := make([]byte, 2)
		binary.BigEndian.PutUint16(classBuff, 1)

		ttlBuff := make([]byte, 4)
		binary.BigEndian.PutUint32(ttlBuff, 60)

		lengthBuff := make([]byte, 2)
		binary.BigEndian.PutUint16(lengthBuff, 4)

		ipAddress := net.ParseIP("8.8.8.8").To4()

		answerObj := answer.Answer{
			Name:   name,
			Type:   typeBuff,
			Class:  classBuff,
			TTL:    ttlBuff,
			Length: lengthBuff,
			Data:   ipAddress,
		}

		answerResponse = append(answerResponse, answerObj.Name...)
		answerResponse = append(answerResponse, answerObj.Type...)
		answerResponse = append(answerResponse, answerObj.Class...)
		answerResponse = append(answerResponse, answerObj.TTL...)
		answerResponse = append(answerResponse, answerObj.Length...)
		answerResponse = append(answerResponse, answerObj.Data...)
	}

	return answerResponse, nil
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
		_, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		qdcount := binary.BigEndian.Uint16(buf[4:6])

		headerResponse, err := createHeaderSection(buf, qdcount)
		if err != nil {
			fmt.Println("Error creating header:", err)
			break
		}

		response := make([]byte, 0)
		response = append(response, headerResponse...)

		var currentPos uint16 = 12
		for i := uint16(0); i < qdcount; i++ {
			questionResponse := createQuestionSection(buf[currentPos:])
			response = append(response, questionResponse...)
			currentPos += uint16(len(questionResponse))
		}

		answerResponse, err := createAnswerSection(buf[12:], int(qdcount))
		if err != nil {
			fmt.Println("Error creating answer:", err)
			break
		}
		response = append(response, answerResponse...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
