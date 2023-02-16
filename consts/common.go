package consts

import (
	"bytes"
	"strings"
)

func IsAuthChannel(channel string) bool {
	if strings.EqualFold(EventChannelAccount, channel) {
		return true
	}
	if strings.EqualFold(channel, PositionChannel) {
		return true
	}
	return false
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func JsonStringConvert(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte{'\\'}, []byte{})
	data = bytes.TrimSpace(data)
	data = bytes.TrimSpace(bytes.Replace(data, newline, space, -1))
	if bytes.HasPrefix(data, []byte{'"'}) {
		data = data[1:]
	}
	if bytes.HasSuffix(data, []byte{'"'}) {
		data = data[:len(data)-1]
	}
	return data
}

// DeleteDoubleQuotationMark 删除redis字符串双英豪
func DeleteDoubleQuotationMark(data string) string {
	if strings.HasPrefix(data, "\"") && strings.HasSuffix(data, "\"") {
		data = data[0 : len(data)-1]
	}
	data = strings.ReplaceAll(data, "\\", "")
	return data
}
