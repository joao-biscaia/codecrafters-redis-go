package main

import (
	"io"
	"log"
	"net"

	constants "github.com/codecrafters-io/redis-starter-go/app/utils/consts"
	"github.com/codecrafters-io/redis-starter-go/app/utils/errorsUtil"
	executeCommand "github.com/codecrafters-io/redis-starter-go/app/utils/execute-command"
	"github.com/codecrafters-io/redis-starter-go/app/utils/parser"
	"github.com/codecrafters-io/redis-starter-go/app/utils/serializer"
)

func closeListener(c io.Closer, msg string) {
	err := c.Close()
	if err != nil {
		err = errorsUtil.Wrap(err, "%v", msg)
	}
}

func handleConn(c net.Conn) {
	defer closeListener(c, "close connection")

	buf := make([]byte, 1024)
	for {
		_, err := c.Read(buf)
		if err != nil {
			log.Printf("read: %v", err)
			return
		}

		args, err := parser.ParseCommand(buf)
		if err != nil {
			log.Printf("parse command: %v", err)
			return
		}
		ex := &executeCommand.ExecuteCommand{
			Args: args,
		}
		out, outType := ex.Run()

		switch outType {
		case constants.SimpleString, constants.BulkString:
			out = out.(string)
			break
		case constants.Array:
			out = out.([]string)
		}

		s := &serializer.Serializer{
			Output:  out,
			OutType: outType,
		}

		encodedOutput := s.Encode()

		_, err = c.Write(encodedOutput)
		if err != nil {
			log.Printf("write: %v", err)
			return
		}
	}
}
