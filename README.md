# TCP 
A simple concurrent TCP test program that sends a message to a server and waits for a response.

## Usage
Open two terminal windows and run the server in one and the client in the other.

### Server

By default the server listens on port 4000. 
```bash
$ go run main.go
```


### Client
```bash
$ go run main.go client
```

The clients will send a message to the server and wait for a response. The server will respond with an acknowledgement message.
The script here tests the server with 150 clients concurrently.



