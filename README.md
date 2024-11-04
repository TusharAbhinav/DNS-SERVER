
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

#### Header Section Flags
```
|  Bit  |  7  |  6  |  5  |  4  |  3  |  2  |  1  |  0  |
|-------|-----|-----|-----|-----|-----|-----|-----|-----|
| Value |  0  |  0  |  0  |  0  | Opcode  |  AA  |  TC  |  RD  |
```

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