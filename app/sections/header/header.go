package header

type Header struct {
	ID      uint16 // Identification field (16 bits)
	QR      uint16  // Query/Response flag (1 bit)
	OPCODE  uint16  // Opcode field (4 bits)
	AA      uint16  // Authoritative Answer flag (1 bit)
	TC      uint16  // Truncated flag (1 bit)
	RD      uint16 // Recursion Desired flag (1 bit)
	RA      uint16 // Recursi16n Available flag (1 bit)
	Z       uint16  // Reserved for future use (3 bits)
	RCODE   uint16 // Response code (4 bits)
	QDCOUNT uint16 // Question Count (16 bits)
	ANCOUNT uint16 // Answer Record Count (16 bits)
	NSCOUNT uint16 // Authority Record Count (16 bits)
	ARCOUNT uint16 // Additional Record Count (16 bits)
}
