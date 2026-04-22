package executeCommand

import (
	"errors"
	"log"
	"strings"

	constants "github.com/codecrafters-io/redis-starter-go/app/utils/consts"
	"github.com/codecrafters-io/redis-starter-go/app/utils/storage"
)

var (
	ECHO = "ECHO"
	PING = "PING"
	SET  = "SET"
	GET  = "GET"
)

type commandFunc map[string]func(args []string) (string, byte, error)

type ExecuteCommand struct {
	Args []string
}

func (e *ExecuteCommand) Run() (string, byte) {
	commands := commandFunc{
		ECHO: e.runEcho,
		PING: e.runPing,
		SET:  e.runSet,
		GET:  e.runGet,
	}
	if len(e.Args) < 1 {
		return "", constants.SimpleString
	}
	cmd, ok := commands[e.Args[0]]
	if ok {
		out, outType, err := cmd(e.Args[1:])
		if err != nil {
			log.Printf("Error executing command: %v", err)
		}
		return out, outType
	}
	return "", constants.SimpleString
}

func (e *ExecuteCommand) runEcho(args []string) (string, byte, error) {
	return strings.Join(args, " "), constants.BulkString, nil
}

func (e *ExecuteCommand) runPing(args []string) (string, byte, error) {
	return "PONG", constants.SimpleString, nil
}

func (e *ExecuteCommand) runSet(args []string) (string, byte, error) {
	if len(args) != 2 {
		return "", ' ', errors.New("invalid SET command")
	}
	key := args[0]
	value := args[1]

	storage.Store(key, value)
	return "OK", constants.SimpleString, nil
}

func (e *ExecuteCommand) runGet(args []string) (string, byte, error) {
	if len(args) != 1 {
		return "", ' ', errors.New("invalid GET command")
	}
	key := args[0]
	value, ok := storage.Get(key)
	if !ok {
		return "", constants.NullBulkString, nil
	}
	return value, constants.BulkString, nil
}
