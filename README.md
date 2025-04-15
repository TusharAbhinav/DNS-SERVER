
[![progress-banner](https://backend.codecrafters.io/progress/dns-server/9ab05740-d833-4fb0-90dd-5328dc2f019c)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)
# Build Your Own DNS Server
A DNS forwarding server implementation in Go that demonstrates the core concepts of the Domain Name System protocol. This server accepts DNS queries and forwards them to a specified upstream DNS resolver, handling the response back to the client.
## Features
- DNS packet parsing and creation
- Support for DNS message headers and questions sections
- DNS query forwarding to upstream resolvers
- UDP protocol handling
- Configurable upstream DNS resolver
- Support for handling multiple DNS questions in a single query
## Technical Details
### DNS Message Format Handling
The server handles the standard DNS message format:
- Header Section (12 bytes)
- Question Section
- Answer Section
- Authority Section (not implemented in this version)
- Additional Section (not implemented in this version)

#### Header Section Format (12 bytes)
```
 0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      ID                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QDCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    NSCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ARCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

#### Header Section Flags
- **ID**: 16-bit identifier assigned by the program making the query
- **QR**: Query/Response Indicator (0 for query, 1 for response)
- **Opcode**: Kind of query (0 for standard query, 1 for inverse query, 2 for server status)
- **AA**: Authoritative Answer flag
- **TC**: Truncation flag
- **RD**: Recursion Desired flag
- **RA**: Recursion Available flag
- **Z**: Reserved for future use (must be zero)
- **RCODE**: Response code (0 for no error, 1 for format error, etc.)
- **QDCOUNT**: Number of entries in the question section
- **ANCOUNT**: Number of resource records in the answer section
- **NSCOUNT**: Number of name server resource records
- **ARCOUNT**: Number of resource records in the additional records section

#### Question Section Format
```
 0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                                               |
/                     QNAME                     /
/                                               /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     QTYPE                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     QCLASS                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

#### Question Section Fields
- **QNAME**: A sequence of labels, where each label consists of a length octet followed by that number of octets. The domain name is represented as a sequence of labels, terminated by a zero-length label.
- **QTYPE**: A two octet code specifying the type of the query. Common values:
  - 1: A (IPv4 address)
  - 2: NS (Authoritative name server)
  - 5: CNAME (Canonical name for an alias)
  - 15: MX (Mail exchange)
  - 28: AAAA (IPv6 address)
- **QCLASS**: A two octet code specifying the class of the query. Common values:
  - 1: IN (Internet)
  - 3: CH (Chaos)
  - 4: HS (Hesiod)

#### Answer Section Format
```
 0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                                               |
/                     NAME                      /
/                                               /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     TYPE                      |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     CLASS                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     TTL                       |
|                                               |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   RDLENGTH                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
/                     RDATA                     /
/                                               /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

#### Answer Section Fields
- **NAME**: The domain name that was queried, in the same format as QNAME or using name compression
- **TYPE**: Two octets containing one of the RR type codes (same as QTYPE values)
- **CLASS**: Two octets which specify the class of the data in the RDATA field (same as QCLASS values)
- **TTL**: A 32-bit unsigned integer specifying the time interval (in seconds) that the resource record may be cached
- **RDLENGTH**: A 16-bit integer specifying the length of the RDATA field
- **RDATA**: A variable length string of octets that describes the resource. The format varies according to the TYPE and CLASS fields

## Getting Started
### Prerequisites
- Go 1.19 or higher
- Basic understanding of DNS and networking concepts
### Installation
1. Clone the repository:
```bash
git clone https://github.com/TusharAbhinav/DNS-SERVER.git
cd dns-server
```
2. Run the server:
```bash
go run app/main.go --resolver <upstream-dns-ip:port>
```
Example:
```bash
go run app/main.go --resolver 8.8.8.8:53
```
The server listens on `127.0.0.1:2053` by default.
## Implementation Details
The DNS server implements the following key components:
1. **Header Processing**: Creates and parses DNS header sections with proper flags and counts
2. **Question Processing**: Handles DNS questions with domain name labels and types
3. **DNS Forwarding**: Forwards queries to upstream resolver and processes responses
4. **Label Compression**: Supports DNS name compression in both questions and answers
5. **UDP Communication**: Manages UDP connections for both client and upstream resolver communication
## Assumptions About the Tester And Resolver 
Here are a few assumptions you can make about the tester:
1. **It will always send you queries for A record type. So your parsing logic only needs to take care of this**.

Here are few assumptions you can make about the DNS server you are forwarding the requests to:
1. **It will always respond with an answer section for the queries that originate from the tester.**
2. **It will not contain other sections like (authority section and additional section)**
3. **It will only respond when there is only one question in the question section. If you send multiple questions in the question section, it will not respond at all. So when you receive multiple questions in 
     the question section you will need to split it into two DNS packets and then send them to this resolver then merge the response in a single packet.**
## Project Structure
```
├── app/
│   ├── main.go                    # Main server implementation
│   └── sections/                  # DNS message sections
│       ├── header/                # Header section structure
│       └── question/              # Question section structure
        └── answer/                # Answer section structure
```
## Dependencies
- Standard Go libraries for networking and binary encoding
- Custom DNS message section handlers for header and question processing
## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
## License
This project is licensed under the MIT License - see the LICENSE file for details.
## Acknowledgments
- RFC 1034 & 1035 for DNS protocol specifications
- The Go community for excellent networking libraries
