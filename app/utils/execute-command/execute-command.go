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
	ECHO   = "ECHO"
	PING   = "PING"
	SET    = "SET"
	GET    = "GET"
	RPUSH  = "RPUSH"
	LRANGE = "LRANGE"
)

type commandFunc map[string]func(args []string) (any, byte, error)

type ExecuteCommand struct {
	Args []string
}

func (e *ExecuteCommand) Run() (any, byte) {
	commands := commandFunc{
		ECHO:   e.runEcho,
		PING:   e.runPing,
		SET:    e.runSet,
		GET:    e.runGet,
		RPUSH:  e.runRPUSH,
		LRANGE: e.runLRANGE,
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

func (e *ExecuteCommand) runEcho(args []string) (any, byte, error) {
	return strings.Join(args, " "), constants.BulkString, nil
}

func (e *ExecuteCommand) runPing(args []string) (any, byte, error) {
	return "PONG", constants.SimpleString, nil
}

func (e *ExecuteCommand) runSet(args []string) (any, byte, error) {
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

func (e *ExecuteCommand) runGet(args []string) (any, byte, error) {
	if len(args) != 1 {
		return "", ' ', errors.New("invalid GET command")
	}
	key := args[0]
	value, ok := storage.Get[string](key)
	if !ok {
		return "", constants.NullBulkString, nil
	}
	return value, constants.BulkString, nil
}

func (e *ExecuteCommand) runRPUSH(args []string) (any, byte, error) {
	if len(args) < 2 {
		return "", ' ', errors.New("invalid RPUSH command")
	}
	key := args[0]
	value := args[1:]
	s := storage.Push(key, value...)
	return strconv.Itoa(s), constants.Integer, nil
}

func (e *ExecuteCommand) runLRANGE(args []string) (any, byte, error) {
	if len(args) < 3 {
		return make([]string, 0), ' ', errors.New("invalid LRANGE command")
	}
	key := args[0]
	s, err := strconv.Atoi(args[1])
	if err != nil {
		return make([]string, 0), ' ', errors.New("invalid LRANGE command: start index")
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return make([]string, 0), ' ', errors.New("invalid LRANGE command: end index")
	}
	values, ok := storage.Get[[]string](key)
	if !ok {
		return make([]string, 0), constants.Array, nil
	}
	if s >= len(values) {
		return make([]string, 0), constants.Array, nil
	}
	if end >= len(values) {
		return values[s:], constants.Array, nil
	}
	if s > end {
		return make([]string, 0), constants.Array, nil
	}
	return values[s : end+1], constants.Array, nil
}
