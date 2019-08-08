package rpccalc

import (
	"os"
	"time"
)

// FifoConn is a JSON-RPC 2.0 transport mechanism that uses named pipes for
// receiving requests and sending responses to and from clients.
// FifoConn satisfies the rpc.Conn interface in github.com/ethereum/go-ethereum/rpc
type FifoConn struct {
	RequestPath, ResponsePath string
	request, response         *os.File
}

// Read reads all data from the request pipe. If the request pipe has not yet beeen
// opened, this function will first attempt to open the pipe.
func (f *FifoConn) Read(buf []byte) (int, error) {
	if f.request == nil {
		reqFile, err := os.OpenFile(f.RequestPath, os.O_RDONLY, 0600)
		if err != nil {
			return 0, err
		}
		f.request = reqFile
	}
	return f.request.Read(buf)
}

// Write sends all data in buf to the response pipe. If the response pipe has not
// yet beeen opened, this function will first attempt to open the pipe.
func (f *FifoConn) Write(buf []byte) (int, error) {
	if f.response == nil {
		respFile, err := os.OpenFile(f.ResponsePath, os.O_WRONLY, 0600)
		if err != nil {
			return 0, err
		}
		f.response = respFile
	}
	return f.response.Write(buf)
}

// Close closes both the request and response pipe
func (f *FifoConn) Close() error {
	f.request.Close()
	f.response.Close()
	f.response, f.request = nil, nil
	return nil
}

// SetWriteDeadline sets the write deadline for writing to the response pipe.
func (f *FifoConn) SetWriteDeadline(t time.Time) error {
	if f.response == nil {
		respFile, err := os.OpenFile(f.ResponsePath, os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		f.response = respFile
	}

	return f.response.SetWriteDeadline(t)
}

// NewFifoConn creates a new FifoConn instance used for transporting JSON-RPC
// requests and responses on two separate named pipes, each located at requestPath
// and responsePath.
func NewFifoConn(requestPath, responsePath string) *FifoConn {
	return &FifoConn{
		RequestPath:  requestPath,
		ResponsePath: responsePath,
	}
}
