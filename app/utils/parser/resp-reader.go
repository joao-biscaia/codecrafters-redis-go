package parser

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/utils/errorsUtil"
)

const (
	BulkString byte = '$'
	Array      byte = '*'
)

type Value struct {
	typ     byte
	integer int
	str     []byte
	array   []Value
	null    bool
}

func ParseCommand(input []byte) (parsedCommand []string, e error) {
	reader := bufio.NewReader(strings.NewReader(string(input)))

	b, err := reader.ReadByte()
	if err != nil {
		return nil, errorsUtil.Wrap(err, "%v", "read command type")
	}

	if b != Array {
		log.Printf("command is not RESP array: %v", input)
		return nil, nil
	}

	byteSize, err := readLine(reader)
	if err != nil {
		return nil, errorsUtil.Wrap(err, "%v", "read command size")
	}
	arraySize := byteSize

	commandArray := make([]string, arraySize)
	var builder strings.Builder

	for i := range arraySize {
		bs, _ := reader.ReadByte()
		if bs != BulkString {
			log.Printf("arg %v isn't bulk string; %d", input, len(input))
			return nil, errorsUtil.New("arg %b isn't bulk string; %v", bs, commandArray)
		}
		bulkSize, _ := readLine(reader)

		for _ = range bulkSize {
			char, _ := reader.ReadByte()
			builder.WriteByte(char)
		}
		argString := builder.String()
		commandArray[i] = argString

		builder.Reset()
		// \r\n after end of string
		_, _ = reader.ReadByte()
		_, _ = reader.ReadByte()
	}

	return commandArray, nil
}

func readLine(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n') // lê "10\r\n"
	if err != nil {
		return 0, err
	}
	trimmed := strings.TrimSpace(line) // "10"
	return strconv.Atoi(trimmed)
}
