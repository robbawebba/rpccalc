package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/robbawebba/rpccalc"
	"github.com/urfave/cli"
)

// fifoCommand is a CLI command to listen for and serve calculator service requests
// over a Fifo-based connections.
var fifoCommand = cli.Command{
	Name:      "fifo",
	Aliases:   []string{"f"},
	Usage:     "Start the calculator service using named pipes / fifos as the transport",
	Action:    listenAndServeFifo,
	ArgsUsage: "[/path/to/request/pipe] [/path/to/request/pipe]",
	Before:    initServer,
}

// listenAndServeFifo listens for and serves calculator service requests using named pipes
// as the transport mechanism. After successful initialization, this function will never
// return and the server will handle connections until the process is killed,
func listenAndServeFifo(c *cli.Context) error {
	if c.NArg() != 2 {
		return fmt.Errorf("Invalid number of arguments. Expected two paths to named pipes")
	}

	reqPath := c.Args().Get(0)
	respPath := c.Args().Get(1)

	// Check both paths to make sure they exist and are named pipes. If a pipe
	// does do not exist, then attempt to create the named pipe at the specified
	// location.
	for _, path := range []string{reqPath, respPath} {
		pathInfo, err := os.Stat(path)
		if err != nil {
			// Attempt to create named pipe at the provided path
			if err := syscall.Mkfifo(path, 0600); err != nil {
				return fmt.Errorf("Unable to create a named pipe at \"%s\": %v", path, err)
			}
		} else {
			// err is nil, so the file already exists
			if (pathInfo.Mode() & os.ModeNamedPipe) == 0 {
				return fmt.Errorf("A file at \"%s\" already exists and is not a named pipe", path)
			}
		}
	}

	conn := rpccalc.NewFifoConn(reqPath, respPath)
	defer conn.Close()
	log.Printf("Listening for requests at %s and sending responses to %s...\n", reqPath, respPath)

	// Continuously handle requests on the connection.
	// Note that the server blocks internally while opening the request pipe for
	// reading until another entity opens the file for writing. This can easily be
	// be observed running this server with strace(1).

	// One major difference betweeb FifoConn and socket-based connections such as
	// net.UnixConn or net.TCPConn is that Fifoconn does not support simultaneous
	// connections from clients, which is why the requests from clients in this server
	// are handled serially, as opposed to handling concurrent connections in goroutines.
	for {
		codec := rpc.NewJSONCodec(conn)
		server.ServeCodec(codec, 0)
	}
}
