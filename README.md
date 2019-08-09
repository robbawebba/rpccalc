# rpccalc

A JSON-RPC 2.0 service over a variety of transports

## Server Installation and Running
1. Clone this repository: `$ go get github.com/robbawebba/rpccalc/cmd/calcd`
2. Run the application! `$ calcd help` or `calcd fifo`

  If `$GOPATH/bin` is not in your `$PATH` environment variable, then you will
  likely have to build the server manually before running:
  ```bash
  $ go build github.com/robbawebba/rpccalc/cmd/calcd
  $ ./calcd help
  ```

## Supported transports
### fifo
The `calcd fifo` command starts a server using two named pipes, a request pipe and a response pipe, for transportation. It only supports one client connection at a time.

either of the named pipes do not have to be created before-hand. The server will attempt
to create the pipes if they do not already exist.

#### example
```bash
# Start server in the background
$ calcd fifo /tmp/req /tmp/res &
$ echo '{"jsonrpc": "2.0", "method": "calc_add", "params": [42, 23], "id": 1}' > /tmp/req
$ cat /tmp/res
# prints {"jsonrpc":"2.0","id":1,"result":65}
$ echo '{"jsonrpc": "2.0", "method": "calc_subtract", "params": [42, 23], "id": 1}' > /tmp/req
$ cat /tmp/res
# prints {"jsonrpc":"2.0","id":1,"result":19}
```

See the output of `calcd fifo help` for usage information.
