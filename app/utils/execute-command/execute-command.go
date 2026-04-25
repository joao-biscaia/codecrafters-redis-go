package executeCommand

import (
	"errors"
	"log"
	"strconv"
	"strings"

	constants "github.com/codecrafters-io/redis-starter-go/app/utils/consts"
	"github.com/codecrafters-io/redis-starter-go/app/utils/storage"
)

var (
	ECHO  = "ECHO"
	PING  = "PING"
	SET   = "SET"
	GET   = "GET"
	RPUSH = "RPUSH"
)

type commandFunc map[string]func(args []string) (string, byte, error)

type ExecuteCommand struct {
	Args []string
}

func (e *ExecuteCommand) Run() (string, byte) {
	commands := commandFunc{
		ECHO:  e.runEcho,
		PING:  e.runPing,
		SET:   e.runSet,
		GET:   e.runGet,
		RPUSH: e.runRPUSH,
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
	if len(args) < 2 {
		return "", ' ', errors.New("invalid SET command")
	}
	key := args[0]
	value := args[1]
	if len(args) > 2 {
		durationMeasure := args[2]
		duration, err := strconv.Atoi(args[3])
		if err != nil {
			return "", ' ', errors.New("invalid SET expiry")
		}
		storage.StoreWithExpiry(key, value, durationMeasure, duration)
		return "OK", constants.SimpleString, nil
	}

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

func (e *ExecuteCommand) runRPUSH(args []string) (string, byte, error) {
	if len(args) < 2 {
		return "", ' ', errors.New("invalid RPUSH command")
	}
	key := args[0]
	value := args[1:]
	s := storage.Push(key, value...)
	return strconv.Itoa(s), constants.Integer, nil
}
