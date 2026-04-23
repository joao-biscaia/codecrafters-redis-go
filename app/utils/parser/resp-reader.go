package parser

import (
	"bufio"
	"log"
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

	byteSize, err := reader.ReadByte()
	if err != nil {
		return nil, errorsUtil.Wrap(err, "%v", "read command size")
	}
	arraySize := int(byteSize - '0')

	commandArray := make([]string, arraySize)
	// \r\n after array size
	_, err = reader.ReadByte()
	_, err = reader.ReadByte()
	if err != nil {
		return nil, errorsUtil.Wrap(err, "%v", "read CRLF token")
	}
	var builder strings.Builder

	for i := range arraySize {
		bs, _ := reader.ReadByte()
		if bs != BulkString {
			log.Printf("arg %v isn't bulk string; %d", input, len(input))
			return nil, errorsUtil.New("arg %b isn't bulk string; %v", bs, commandArray)
		}
		byteSize, _ := reader.ReadByte()
		bulkSize := int(byteSize - '0')

		// \r\n after bulksize
		_, _ = reader.ReadByte()
		_, _ = reader.ReadByte()
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
