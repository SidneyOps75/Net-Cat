# net-cat

This project is a simple TCP chat server implementation in Go. It allows multiple clients to connect and chat with each other in a group.

## Features

- TCP connection between server and multiple clients
- Requirement for clients to provide a name
- Limit on the maximum number of concurrent connections (10)
- Clients can send messages to the chat
- Empty messages are not broadcasted
- Messages are formatted with timestamp, sender name, and message content
- New clients receive the chat history upon connection
- Clients are notified when another client joins or leaves the chat

## Prerequisites

- Go 1.18+: Make sure you have Go installed. You can download it [here](https://go.dev/doc).
- A data file containing numeric values (one per line).

## Installation

Clone this repository:

 ```bash
    git clone https://learn.zone01kisumu.ke/git/bshisia/net-cat.git
```

Navigate to the project directory:

```bash
    cd net-cat
```

## Usage
1. Run the server:
```
    go build -o TCPChat
```
```bash
./TCPChat  # Default port 8989
# or
./TCPChat # Custom port
```
2. Connect to the server using a TCP client like nc:
```bash
nc localhost 8989  # or whatever port you specified
```
3. Enter your name when prompted, and start chatting!


## License

This project is licensed under the MIT License.

## Testing 
To run the tests present go to the root directory and run the following command: 
```
go test -v ./utils
```


## Contributing

If you have suggestions for improvements, bug fixes, or new features, feel free to open an issue or submit a pull request.

## Author

This project was build and maintained by:

[Thadeus Ogondola](https://learn.zone01kisumu.ke/git/togondol/)

[John Opiyo](https://learn.zone01kisumu.ke/git/johopiyo/)

[Brian Shisia](https://learn.zone01kisumu.ke/git/bshisia/)
