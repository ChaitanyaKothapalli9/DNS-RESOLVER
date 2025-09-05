# DNS Resolver

This project is a simple **DNS Resolver** implemented in **Go**.  
It demonstrates how DNS queries are built, sent over UDP, and how responses are parsed to resolve domain names.

## Project Structure

DNS_Resolver/
├── main.go             # Entry point of the program
├── dnsresolver.go      # Core logic for DNS message creation and parsing
├── udpclient.go        # UDP client implementation for sending/receiving queries
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums

> ⚠️ Compiled binaries (`dnsresolver`, `main.exe`) and PDFs should be excluded from version control.

## Getting Started

### Prerequisites
- Go 1.18+ (or any recent version)

### Clone the Repository
git clone https://github.com/your-username/DNS_Resolver.git  
cd DNS_Resolver

### Run the Project
go run main.go

### Build Executable
go build -o dnsresolver main.go  
./dnsresolver

## Implementation Details
- Constructs DNS queries manually.  
- Sends requests via UDP sockets.  
- Receives and parses DNS responses.  
- Demonstrates networking and protocol-level programming in Go.

## Notes
Ignore compiled binaries and PDFs when committing:  
- dnsresolver  
- main.exe  

## License
This project is for educational purposes.  
You are free to modify and use it for learning.
