package rpccalc

import (
	"syscall"

	"github.com/ethereum/go-ethereum/rpc"
)

type CalcService struct{}

func (s *CalcService) Add(a, b int) int {
	return a + b
}

func (s *CalcService) Subtract(a, b int) int {
	return a - b
}

func main() {
	calculator := new(CalcService)
	server := rpc.NewServer()
	server.RegisterName("calculator", calculator)

	//l, _ := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: "/tmp/calculator.sock"})
	syscall.Mkfifo("/tmp/request", 0600)
	syscall.Mkfifo("/tmp/response", 0600)
	conn := NewFifoConn("/tmp/request", "/tmp/response")
	for {
		codec := rpc.NewJSONCodec(conn)
		server.ServeCodec(codec, 0)
	}

}
