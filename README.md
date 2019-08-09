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

#### Example
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

See the output of `calcd help fifo` for usage information.

### TCP

The `calcd tcp` command starts a server that listens for TCP connections on the
specified port, or port 2000 if none is specified. This transport supports multiple
client connections and supports streams of requests to the TCP connection from the
client.

#### Example
```bash
# Start server in the background while it listens for connections on port 2000
$ calcd tcp &
# Send one request to the calculator service via TCP
$ echo '{"jsonrpc": "2.0", "method": "calc_add", "params": [42, 23], "id": 1}' | netcat -N localhost 2000
# Send several requests to the calcilator service via TCP
$ for i in {0..10}; do
    echo "{\"jsonrpc\": \"2.0\", \"method\": \"calc_subtract\", \"params\": [10, $i], \"id\": $i}"
  done | netcat localhost 2000
# Prints something like...
# {"jsonrpc":"2.0","id":1,"result":9}
# {"jsonrpc":"2.0","id":0,"result":10}
# {"jsonrpc":"2.0","id":2,"result":8}
# {"jsonrpc":"2.0","id":3,"result":7}
# {"jsonrpc":"2.0","id":5,"result":5}
# {"jsonrpc":"2.0","id":4,"result":6}
# {"jsonrpc":"2.0","id":7,"result":3}
# {"jsonrpc":"2.0","id":6,"result":4}
# {"jsonrpc":"2.0","id":9,"result":1}
# {"jsonrpc":"2.0","id":8,"result":2}
# {"jsonrpc":"2.0","id":10,"result":0}
```

See the output of `calcd help tcp` for usage information.
