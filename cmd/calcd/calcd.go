package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/robbawebba/rpccalc"
	"github.com/urfave/cli"
)

var (
	server *rpc.Server
)

// initServer initializes a JSON-RPC 2.0 server and registers
// rpccalc.CalculatorService with this server.
func initServer(ctx *cli.Context) error {
	server = rpc.NewServer()
	return server.RegisterName(`calc`, &rpccalc.CalcService{})
}

func main() {
	app := cli.NewApp()
	app.Name = "calcd"
	app.Usage = "A JSON-RPC-based calculator service that supports addition and subtraction"

	app.Commands = []cli.Command{
		fifoCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
