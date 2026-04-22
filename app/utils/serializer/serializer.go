package serializer

import (
	"log"
	"strconv"
	"strings"

	constants "github.com/codecrafters-io/redis-starter-go/app/utils/consts"
)

type Serializer struct {
	Output  string
	OutType byte
}

type encodeFunctions map[byte]func(nakedString string) ([]byte, error)

func (s *Serializer) Encode() []byte {
	encodeMap := encodeFunctions{
		constants.SimpleString:   encodeSimpleString,
		constants.BulkString:     encodeBulkString,
		constants.NullBulkString: encodeNullBulkString,
	}

	encodeFunc, ok := encodeMap[s.OutType]
	if ok {
		serializedBytes, err := encodeFunc(s.Output)
		if err != nil {
			log.Printf("error while serializing response: %v", err)
			return nil
		}
		return serializedBytes
	}
	log.Printf("%v: not valid encode type", s.OutType)
	return nil
}

func encodeNullBulkString(nakedString string) ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte('$')
	builder.WriteByte('-')
	builder.WriteByte('1')
	builder.WriteByte('\r')
	builder.WriteByte('\n')
	return []byte(builder.String()), nil
}

func encodeSimpleString(nakedString string) ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte(constants.SimpleString)
	builder.Write([]byte(nakedString))
	builder.WriteByte('\r')
	builder.WriteByte('\n')
	return []byte(builder.String()), nil
}

func encodeBulkString(nakedString string) ([]byte, error) {
	outSize := len(nakedString)
	var builder strings.Builder
	builder.WriteByte(constants.BulkString)
	b := []byte(strconv.Itoa(outSize))
	for i := range len(b) {
		builder.WriteByte(b[i])
	}
	builder.WriteByte('\r')
	builder.WriteByte('\n')

	for _, v := range nakedString {
		builder.WriteByte(byte(v))
	}
	builder.WriteByte('\r')
	builder.WriteByte('\n')

	return []byte(builder.String()), nil
}
