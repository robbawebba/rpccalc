# rpccalc

A JSON-RPC 2.0 calculator service that supports a variety of transports

## Server Installation and Running
1. Clone this repository: `$ go get github.com/robbawebba/rpccalc/cmd/calcd`
2. Run the application!

  ```
  $ calcd help
  ```

  If `$GOPATH/bin` is not in your `$PATH` environment variable, then you will
  likely have to use an absolute path to run the program or build the program
  manually:

  ```bash
  # Absolute path to calcd
  $ $GOPATH/bin/calcd help
  # or build yourself
  $ go build github.com/robbawebba/rpccalc/cmd/calcd
  $ ./calcd help
  ```

## Supported Transports
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

## Design Notes

`rpccalc.FifoConn` is an easily interchangeable RPC server transport because it
satisfies the `net.Conn` and `rpc.Conn` interfaces that allow generic reading and
writing of streamable data. The only differences between the TCP and fifo servers
found in `cmd/calcd` are the boilerplate logic needed to initialize the the different
connections. The RPC server has no knowledge or preference for the sort of connection
it receives when handling a request. As long as the RPC server can read to and write
from the connection it receives, it's requirements are satisfied.


Although `rpccalc.FifoConn` implements a common networking interface, `net.Conn` that
supports streaming communication like a TCP connection, the nature of named pipes
causes the `FifoConn` to behave quite differently. Named pipes do not easily support
multiple connected clients, or writers, on the same pipe. It will often lead to undefined
results. This is very different from the behavior of Unix and TCP sockets that support
multiple simultaneous connections on a single socket.

### Next Steps

I would like to improve the operations related to the named pipes. I feel that the
methods are missing a couple of safety features that would provide a more stable
experience for users of the Fifo-based connection. Some packages that provide
some inspiration for further development include https://github.com/containerd/fifo/
and https://github.com/natefinch/npipe.

I would also like to explore the possibility of demonstrating a REST API and RPC
server side-by-side. the HTTP transport system that would be used for the REST API
could easily make use of the `rpccalc.CalculatorService` method definitions by
Wrapping them in functions that implements the [`http.Handler`](https://golang.org/pkg/net/http/#Handler)
interface. This interface allows the developer to create a composable system of HTTP
handlers for a variety of routes with [`http.ServeMux`](https://golang.org/pkg/net/http/#ServeMux).
I would use these tools to create a router for the `/add` and `/subtract` routes
and then listen for and serve HTTP requests for these routes.
