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

var _ = net.ListenUDP

// |  Bit  |  7  |  6  |  5  |  4  |  3  |  2  |  1  |  0  |
// |-------|-----|-----|-----|-----|-----|-----|-----|-----|
// | Value |  0  |  0  |  0  |  0  | Opcode  |  AA  |  TC  |  RD  |

func createHeaderSection(buf []byte) ([]byte, error) {
	buffer := bytes.NewReader(buf[:2])
	var ID, QDCOUNT, ANCOUNT, OPCODE, RD, RCODE uint16

	idErr := binary.Read(buffer, binary.BigEndian, &ID)
	if idErr != nil {
		fmt.Println("Error reading ID:", idErr)
		return nil, idErr
	}

	flags := buf[2]

	// Extract OPCODE (bits 3-6)
	OPCODE = uint16((flags >> 3) & 0b00001111)

	// Extract RD (bit 0)
	RD = uint16(flags & 0b00000001)

	// Set QDCOUNT and ANCOUNT
	qdBuffer := bytes.NewBuffer(buf[4:6])
	qdErr := binary.Read(qdBuffer, binary.BigEndian, &QDCOUNT)
	if qdErr != nil {
		fmt.Println("Error reading QDCOUNT:", qdErr)
		return nil, qdErr
	}
	anBuffer:=bytes.NewBuffer(buf[6:8])
	anErr:=binary.Read(anBuffer,binary.BigEndian,&ANCOUNT)
	if anErr != nil {
		fmt.Println("Error reading ANCOUNT:", anErr)
		return nil, anErr
	}
	// Set RCODE based on OPCODE
	if OPCODE == 0 {
		RCODE = 0
	} else {
		RCODE = 4
	}

	headerObj := header.Header{
		ID:      ID,
		QR:      1,
		OPCODE:  OPCODE,
		AA:      0,
		TC:      0,
		RD:      RD,
		RA:      0,
		Z:       0,
		RCODE:   RCODE,
		QDCOUNT: QDCOUNT,
		ANCOUNT: ANCOUNT,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	headerResponse := []byte{}

	// ID bytes
	idBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(idBytes, headerObj.ID)
	headerResponse = append(headerResponse, idBytes...)

	// Flags (2 bytes)
	flagBytes := make([]byte, 2)
	flagsValue := headerObj.QR<<15 | headerObj.OPCODE<<11 |
		headerObj.AA<<10 | headerObj.TC<<9 |
		headerObj.RD<<8 | headerObj.RA<<7 |
		headerObj.Z<<4 | headerObj.RCODE
	binary.BigEndian.PutUint16(flagBytes, flagsValue)
	headerResponse = append(headerResponse, flagBytes...)

	// Counts (8 bytes: QDCOUNT, ANCOUNT, NSCOUNT, ARCOUNT)
	countBytes := make([]byte, 8)
	binary.BigEndian.PutUint16(countBytes[0:2], headerObj.QDCOUNT)
	binary.BigEndian.PutUint16(countBytes[2:4], headerObj.ANCOUNT)
	binary.BigEndian.PutUint16(countBytes[4:6], headerObj.NSCOUNT)
	binary.BigEndian.PutUint16(countBytes[6:8], headerObj.ARCOUNT)
	headerResponse = append(headerResponse, countBytes...)

	return headerResponse, nil
}
func createLabel(off uint16, buf []byte) []byte {
	domainName := buf[off:]
	offset := 0
	var Name []byte

	for {
		length := int(domainName[offset])

		// Check if the label contains a pointer
		if (length & 0b11000000) == 0b11000000 {
			pointerOffset := uint16(domainName[offset]&0b00111111)<<8 | uint16(domainName[offset+1])
			createLabel(pointerOffset, buf)
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

	// Initialize type and class buffers
	Name := createLabel(12, buf)
	typeBuff := make([]byte, 2)
	classBuff := make([]byte, 2)

	// Set the record type to 1 (A record) and class to 1 (IN - Internet)
	var Type, Class uint16
	Type = 1
	binary.BigEndian.PutUint16(typeBuff, Type)
	Class = 1
	binary.BigEndian.PutUint16(classBuff, Class)

	questionObj := question.Question{
		Name:  Name,
		Type:  typeBuff,
		Class: classBuff,
	}

	questionResponse := []byte{}
	questionResponse = append(questionResponse, questionObj.Name...)
	questionResponse = append(questionResponse, questionObj.Type...)
	questionResponse = append(questionResponse, questionObj.Class...)

	return questionResponse
}

func createAnswerSection(buf []byte) ([]byte, error) {
	answerResponse := []byte{}

	// Initialize buffers for each field in the answer section
	typeBuff := make([]byte, 2)
	Name := createLabel(12, buf)
	classBuff := make([]byte, 2)
	ttlBuff := make([]byte, 4)
	lengthBuff := make([]byte, 2)

	var Type, Class, Length uint16
	Type = 1
	binary.BigEndian.PutUint16(typeBuff, Type)
	Class = 1
	binary.BigEndian.PutUint16(classBuff, Class)
	var TTL uint32 = 60
	binary.BigEndian.PutUint32(ttlBuff, TTL)
	Length = 4
	binary.BigEndian.PutUint16(lengthBuff, Length)

	ipAddress := net.ParseIP("8.8.8.8").To4()
	dataBuff := new(bytes.Buffer)
	err := binary.Write(dataBuff, binary.BigEndian, ipAddress)
	if err != nil {
		fmt.Println("error encoding data")
		return answerResponse, err
	}
	Data := dataBuff.Bytes()

	answerObj := answer.Answer{
		Name:   Name,
		Type:   typeBuff,
		Class:  classBuff,
		TTL:    ttlBuff,
		Length: lengthBuff,
		Data:   Data,
	}

	answerResponse = append(answerResponse, answerObj.Name...)
	answerResponse = append(answerResponse, answerObj.Type...)
	answerResponse = append(answerResponse, answerObj.Class...)
	answerResponse = append(answerResponse, answerObj.TTL...)
	answerResponse = append(answerResponse, answerObj.Length...)
	answerResponse = append(answerResponse, answerObj.Data...)

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
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
		fmt.Printf("Received %d bytes from %s\n", size, source)

		headerResponse, err := createHeaderSection(buf)
		if err != nil {
			fmt.Println("Error creating header:", err)
			break
		}
		questionResponse := createQuestionSection(buf)
		answerResponse, err := createAnswerSection(buf)
		if err != nil {
			fmt.Println("Error creating answerResponse:", err)
			break
		}
		response := make([]byte, 0)
		response = append(response, headerResponse...)
		response = append(response, questionResponse...)
		response = append(response, answerResponse...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
