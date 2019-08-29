package rpc

import (
	"regexp"
	"strconv"
)

var (
	regex = regexp.MustCompile(`\\[0-7]{3}`)
)

func oct2utf8(in string) string {
	s := []byte(in)

	out := regex.ReplaceAllFunc(s, func(b []byte) []byte {
		i, _ := strconv.ParseInt(string(b[1:]), 8, 0)
		return []byte{byte(i)}
	})
	return string(out)
}
