package constants

const (
	SimpleString   byte = '+'
	Error          byte = '-'
	Integer        byte = ':'
	Array          byte = '*'
	BulkString     byte = '$'
	NullBulkString byte = 'n'
)

var (
	Seconds      string = "EX"
	Milliseconds string = "PX"
)
