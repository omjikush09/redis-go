package app;

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	Intger = ':'
	String = '+'
	Array  = '*'
	Bulk   = '$'
)

func ParseResp(r *bufio.Reader) ([]string, error) {
	prefix, err := r.ReadByte()
	fmt.Print(string(prefix) + " ")
	if err != nil {
		return nil, err
	}

	switch prefix {
	case String:
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		str := make([]string, 1)
		str[0] = line[:len(line)-2]
		return str, nil
	case Array:
		arrLen, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		arrLen = strings.TrimSuffix(arrLen, "\r\n")
		count, err := strconv.Atoi(arrLen)
		if count < 0 {
			return nil, err
		}
		arr := make([]string, count)

		for i := range count {
			args, err := ParseResp(r)
			if err != nil {
				return nil, err
			}
			arr[i] = args[0]
		}
		return arr, nil

	case Bulk:
		len, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		len = strings.TrimSuffix(len, "\r\n")

		count, err := strconv.Atoi(len)
		if count < 0 {
			return nil, fmt.Errorf("not valid len %d", count)
		}
		if err != nil {
			return nil, err
		}
		data := make([]byte, count)
		_, err = io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}
		_, err = r.ReadString('\n') //Reading ending \r\n
		if err != nil {
			return nil, err
		}

		str := make([]string, 1)
		str[0] = string(data[:count])
		return str, nil
	default:
		return nil, fmt.Errorf("unknown prefix %c", prefix)
	}

}

func encodeAsSimpleString(data string) string {
	return fmt.Sprintf("+%s\r\n", data)
}
func encodeAsBulkString(data string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(data), data)
}

func encodeAsInteger(data int) string {
	return fmt.Sprintf(":%d\r\n", data)
}

func encodeAsError(message string) string {
	return fmt.Sprintf("-%s\r\n", message)
}

func encodeAsArray(data []string) string {
	result := fmt.Sprintf("*%d\r\n", len(data))
	for _, item := range data {
		result += item
	}
	return result
}