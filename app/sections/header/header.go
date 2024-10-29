package header

type Header struct {
	ID      uint16 // Identification field (16 bits)
	QR      uint8  // Query/Response flag (1 bit)
	OPCODE  uint8  // Opcode field (4 bits)
	AA      uint8  // Authoritative Answer flag (1 bit)
	TC      uint8  // Truncated flag (1 bit)
	RD      uint8  // Recursion Desired flag (1 bit)
	RA      uint8  // Recursion Available flag (1 bit)
	Z       uint8  // Reserved for future use (3 bits)
	RCODE   uint8  // Response code (4 bits)
	QDCOUNT uint16 // Question Count (16 bits)
	ANCOUNT uint16 // Answer Record Count (16 bits)
	NSCOUNT uint16 // Authority Record Count (16 bits)
	ARCOUNT uint16 // Additional Record Count (16 bits)
}
