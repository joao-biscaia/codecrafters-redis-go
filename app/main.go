package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/utils/errors"
	"io"
	"log"
	"net"
	"os"
)

var (
	listen = flag.String("listen", ":6379", "address to listen to")
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error, %v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	l, err := net.Listen("tcp", *listen)
	if err != nil {
		return errors.Wrap(err, *listen)
	}
	defer closeListener(l, &err, "close listener")
	log.Printf("listening %v", l.Addr())

	for {
		c, err := l.Accept()
		if err != nil {
			return errors.Wrap(err, "accept")
		}
		go handleConn(c)
	}

}

func closeListener(c io.Closer, errp *error, msg string) {
	err := c.Close()
	if *errp == nil {
		*errp = errors.Wrap(err, "%v", msg)
	}
}

func handleConn(c net.Conn) {
	defer closeListener(c, nil, "close connection")

	buf := make([]byte, 1024)
	for {
		_, err := c.Read(buf)
		if err != nil {
			log.Printf("read: %v", err)
			return
		}

		log.Printf("read command:\n %s", buf)

		_, err = c.Write([]byte("+PONG\r\n"))
		if err != nil {
			log.Printf("write: %v", err)
			return
		}
	}
}
