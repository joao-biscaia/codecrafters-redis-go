package executeCommand

import (
	"log"
	"strings"
)

var (
	ECHO = "ECHO"
	PING = "PING"
)

const (
	SimpleString byte = '+'
	Error        byte = '-'
	Integer      byte = ':'
	Array        byte = '*'
	BulkString   byte = '$'
)

type commandFunc map[string]func(args []string) (string, byte, error)

type ExecuteCommand struct {
	Args []string
}

func (e *ExecuteCommand) Run() (string, byte) {
	commands := commandFunc{
		ECHO: e.runEcho,
		PING: e.runPing,
	}
	if len(e.Args) < 1 {
		return "", SimpleString
	}
	cmd, ok := commands[e.Args[0]]
	if ok {
		out, outType, err := cmd(e.Args[1:])
		if err != nil {
			log.Printf("Error executing command: %v", err)
		}
		return out, outType
	}
	return "", SimpleString
}

func (e *ExecuteCommand) runEcho(args []string) (string, byte, error) {
	return strings.Join(args, " "), BulkString, nil
}

func (e *ExecuteCommand) runPing(args []string) (string, byte, error) {
	return "PONG", SimpleString, nil
}
