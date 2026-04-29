package executeCommand

import (
	"errors"
	"log"
	"slices"
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
	LPUSH  = "LPUSH"
	LLEN   = "LLEN"
	LPOP   = "LPOP"
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
		LPUSH:  e.runLPUSH,
		LLEN:   e.runLLEN,
		LPOP:   e.runLPOP,
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
	s := storage.Push[string](key, true, value...)
	return strconv.Itoa(s), constants.Integer, nil
}

func (e *ExecuteCommand) runLRANGE(args []string) (any, byte, error) {
	//todo: refactor this logic to create new package for comparison
	if len(args) < 3 {
		return make([]string, 0), constants.Array, errors.New("invalid LRANGE command")
	}
	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return make([]string, 0), constants.Array, errors.New("invalid LRANGE command: start index")
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return make([]string, 0), constants.Array, errors.New("invalid LRANGE command: end index")
	}

	values, ok := storage.Get[[]string](key)

	if start < 0 {
		start = max(len(values)+start, 0)
	}
	if end < 0 {
		end = max(len(values)+end, 0)
	}

	if !ok {
		return make([]string, 0), constants.Array, nil
	}
	if start >= len(values) {
		return make([]string, 0), constants.Array, nil
	}
	if end >= len(values) {
		return values[start:], constants.Array, nil
	}
	if start > end {
		return make([]string, 0), constants.Array, nil
	}
	return values[start : end+1], constants.Array, nil
}

func (e *ExecuteCommand) runLPUSH(args []string) (any, byte, error) {
	if len(args) < 2 {
		return "0", constants.Integer, errors.New("invalid LPUSH command")
	}
	key := args[0]
	values := make([]string, len(args[1:]))
	copy(values, args[1:])
	slices.Reverse(values)

	s := storage.Push[string](key, false, values...)
	return strconv.Itoa(s), constants.Integer, nil

}

func (e *ExecuteCommand) runLLEN(args []string) (any, byte, error) {
	if len(args) != 1 {
		return 0, constants.Integer, errors.New("invalid LLEN command")
	}
	key := args[0]
	values, ok := storage.Get[[]string](key)
	if ok {
		return strconv.Itoa(len(values)), constants.Integer, nil
	}
	return strconv.Itoa(0), constants.Integer, nil
}

func (e *ExecuteCommand) runLPOP(args []string) (any, byte, error) {
	key := args[0]
	n := 1
	var err error
	if len(args) > 1 {
		n, err = strconv.Atoi(args[1])
		if err != nil {
			return nil, constants.NullBulkString, errors.New("LPOP: invalid command")
		}
		var values []string
		for _ = range n {
			val, ok := storage.Pop[string](key)
			if ok {
				values = append(values, val)
			}
		}
		if values != nil {
			return values, constants.Array, nil
		}
		return make([]string, 0), constants.Array, nil
	}

	val, ok := storage.Pop[string](key)
	if ok {
		return val, constants.BulkString, nil
	}
	return nil, constants.NullBulkString, nil

}
