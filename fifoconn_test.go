package rpccalc

import (
	"bytes"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"
)

const (
	reqTestPipe  = "/tmp/request-test"
	respTestPipe = "/tmp/response-test"
)

func TestFifoConnWrite(t *testing.T) {
	// Only the response pipe is used when FifoConn writes
	if err := syscall.Mkfifo(respTestPipe, 0600); err != nil {
		t.Fatalf("Unable to create Fifo at %s: %v", respTestPipe, err)
	}
	defer os.Remove(respTestPipe)

	conn := NewFifoConn(reqTestPipe, respTestPipe)

	input := []byte("Don't communicate by sharing memory, share memory by communicating.")

	// Write to the connection in a separate goroutine because the write operation
	// on a pipe blocks until another entity opens the pipe for reading
	go func() {
		n, err := conn.Write(input)
		if err != nil {
			t.Errorf("Unable to write data to the connection: %v", err)
		}
		if n != len(input) {
			t.Errorf("Incorrect amount of data written to the connection. Expected %d bytes, written %d bytes", len(input), n)
		}
		conn.Close()
	}()

	// open the output pipe and read it without using the FofoConn's read method
	output, err := ioutil.ReadFile(respTestPipe)
	if err != nil {
		t.Fatalf("Unable to open and read the response pipe: %v", err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("Incorrect data received on the response pipe. Expected: \"%s\", Actual: \"%s\"", string(input), string(output))
	}
}

func TestFifoConnRead(t *testing.T) {
	// Only the request pipe is used when FifoConn writes
	if err := syscall.Mkfifo(reqTestPipe, 0600); err != nil {
		t.Fatalf("Unable to create Fifo at %s: %v", reqTestPipe, err)
	}
	defer os.Remove(reqTestPipe)

	conn := NewFifoConn(reqTestPipe, respTestPipe)
	defer conn.Close()

	input := []byte("Concurrency is not parallelism.")

	// Write to the request pipe in a separate goroutine because the write operation
	// on a pipe blocks until another entity opens the pipe for reading
	go func() {
		err := ioutil.WriteFile(reqTestPipe, input, 0600)
		if err != nil {
			t.Errorf("Unable to write data to the request pipe: %v", err)
		}

	}()

	output := make([]byte, len(input))

	// Read from the connection while data is being written to the request pipe from
	// a separate goroutine
	n, err := conn.Read(output)
	if err != nil {
		t.Fatalf("Unable to read from the connection: %v", err)
	}
	if n != len(input) {
		t.Errorf("Incorrect amount of data read from the connection. Expected %d bytes, read %d bytes", len(input), n)
	}
	if !bytes.Equal(input, output) {
		t.Errorf("Incorrect data received  while reading from the connection. Expected: \"%s\", Actual: \"%s\"", string(input), string(output))
	}
}

func TestFifoConnSetWriteDeadline(t *testing.T) {
	// Only the response pipe is used when FifoConn writes, and SetWriteDeadline
	// only affects write operations on the connection.
	if err := syscall.Mkfifo(respTestPipe, 0600); err != nil {
		t.Fatalf("Unable to create Fifo at %s: %v", respTestPipe, err)
	}
	defer os.Remove(respTestPipe)

	conn := NewFifoConn(reqTestPipe, respTestPipe)
	defer conn.Close()

	// Although this test is not supposed to involve reading from the response pipe,
	// one of the limitations of named pipes is that Opening a file as O-WRONLY blocks
	// until another entity opens the file for reading with O_RDWR or O_RDONLY. So
	// opan and read from the pipe in a separate goroutine here so the setWriteDeadline
	// and Write functions remain unblocked while opening the pipe for writing.
	go func() {
		_, err := ioutil.ReadFile(respTestPipe)
		if err != nil {
			t.Fatalf("Unable to read from the response pipe: %v", err)
		}
	}()

	if err := conn.SetWriteDeadline(time.Now()); err != nil {
		t.Fatalf("Unable to set the write deadline for writing to the connection: %v", err)
	}

	_, err := conn.Write([]byte("This write should fail"))
	if err == nil {
		t.Errorf("No error received when a timeout error was expected")
	}
	if e, ok := err.(*os.PathError); !ok || !e.Timeout() {
		t.Errorf("Error received is not a timeout error %v", err)
	}

}
