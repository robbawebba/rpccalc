package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"
)

// tcpCommand is a CLI command to listen for and serve calculator service requests
// over a TCP connection.
var tcpCommand = cli.Command{
	Name:    "tcp",
	Aliases: []string{"t"},
	Usage:   "Start the calculator service listening on the provided port, or 2000 by default",
	Action:  listenAndServeTCP,
	Before:  initServer,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Value: 2000,
		},
	},
}

// listenAndServeTCP listens for incoming TCP connections and serves calculator
// service requests on that connection. After successful initialization, this
// function will never return and the server will handle connections until the
// process is killed,
func listenAndServeTCP(c *cli.Context) error {

	port := c.Int(`port`)
	if port <= 0 {
		return fmt.Errorf("invalid port number, please specify a positive integer")
	}

	// Create a tcp listener for the specified port
	l, err := net.Listen(`tcp`, `:`+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("Unable to listen for TCP connections on port %d: %v", port, err)
	}
	defer l.Close()
	log.Printf("Listening for connections on port %d...\n", port)

	// Continuously handle requests on the connection.
	// The server blocks on the call to Accept new connections. Once a new client Connects
	// to the TCP port, the Accept function will return and the server will continue
	// to serve that request in a separate goroutine.

	// Unlike rpccalc.FifoConn, TCP connections can handle multiple clients connected
	// simultaneously, so we are able to serve each request in a separate goroutine.
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("Failed to accept incoming connection: %v", err)
		}
		codec := rpc.NewJSONCodec(conn)
		go server.ServeCodec(codec, 0)
	}
}
