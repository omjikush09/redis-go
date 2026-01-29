package app

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRespParser(t *testing.T) {
	examples := []string{
		"+OK\r\n",
		"$5\r\nhello\r\n",
		"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
	}

	for _, ex := range examples {
		r := bufio.NewReader(bytes.NewBufferString(ex))
		val, err := ParseResp(r)
		if err != nil {
			t.Errorf("Error %v\n", err)
		} else {
			fmt.Printf("Input: %-30s => Output: %#v\n", strings.ReplaceAll(ex, "\r\n", "\\r\\n"), val)
		}

	}
}
