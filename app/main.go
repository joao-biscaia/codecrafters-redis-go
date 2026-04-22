package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/utils/errorsUtil"
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
		return errorsUtil.Wrap(err, *listen)
	}
	defer closeListener(l, "close listener")
	log.Printf("listening %v", l.Addr())

	for {
		c, err := l.Accept()
		if err != nil {
			return errorsUtil.Wrap(err, "accept")
		}
		go handleConn(c)
	}

}
