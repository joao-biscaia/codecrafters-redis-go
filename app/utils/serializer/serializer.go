package serializer

import (
	"log"
	"strconv"
	"strings"

	constants "github.com/codecrafters-io/redis-starter-go/app/utils/consts"
)

type Serializer struct {
	Output  any
	OutType byte
}

type encodeFunctions map[byte]func(val any) ([]byte, error)

func (s *Serializer) Encode() []byte {
	encodeMap := encodeFunctions{
		constants.SimpleString:   encodeSimpleString,
		constants.BulkString:     encodeBulkString,
		constants.NullBulkString: encodeNullBulkString,
		constants.Integer:        encodeInteger,
		constants.Array:          encodeArray,
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

func encodeArray(val any) ([]byte, error) {
	outputArray := val.([]string)
	var builder strings.Builder
	builder.WriteByte(constants.Array)
	size := len(outputArray)
	sizeByte := []byte(strconv.Itoa(size))
	builder.Write(sizeByte)
	builder.WriteByte('\r')
	builder.WriteByte('\n')
	for _, v := range outputArray {
		size := len(v)
		sizeByte = []byte(strconv.Itoa(size))
		builder.WriteByte('$')
		builder.Write(sizeByte)
		builder.WriteByte('\r')
		builder.WriteByte('\n')
		builder.Write([]byte(v))
		builder.WriteByte('\r')
		builder.WriteByte('\n')
	}
	return []byte(builder.String()), nil
}

func encodeInteger(output any) ([]byte, error) {
	nakedString := output.(string)
	var builder strings.Builder
	builder.WriteByte(constants.Integer)
	for _, v := range nakedString {
		builder.WriteByte(byte(v))
	}
	builder.WriteByte('\r')
	builder.WriteByte('\n')
	return []byte(builder.String()), nil
}

func encodeNullBulkString(output any) ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte('$')
	builder.WriteByte('-')
	builder.WriteByte('1')
	builder.WriteByte('\r')
	builder.WriteByte('\n')
	return []byte(builder.String()), nil
}

func encodeSimpleString(output any) ([]byte, error) {
	nakedString := output.(string)
	var builder strings.Builder
	builder.WriteByte(constants.SimpleString)
	builder.Write([]byte(nakedString))
	builder.WriteByte('\r')
	builder.WriteByte('\n')
	return []byte(builder.String()), nil
}

func encodeBulkString(output any) ([]byte, error) {
	nakedString := output.(string)
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
