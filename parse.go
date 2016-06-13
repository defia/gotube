package gotube

import (
	"bytes"
	"errors"
	"net/url"
	"regexp"
)

func parsefmtstreamap(b []byte) ([]byte, error) {
	match := fmt_stream_map_Regex.FindSubmatch(b)
	if len(match) < 2 {
		return nil, ErrorNoMatch
	}
	return match[1], nil
}

var (
	titleRegexp          = regexp.MustCompile(`<title>(.+?)</title>`)
	adaptive_fmts_Regex  = regexp.MustCompile(`"adaptive_fmts":"(.+?)"`)
	fmt_stream_map_Regex = regexp.MustCompile(`"url_encoded_fmt_stream_map":"(.+?)"`)
)
var (
	ErrorNoMatch = errors.New("regexp no match")
)

var (
	unicodeAnd = []byte(`\u0026`)
	bytesComma = []byte(",")
	bytesEqual = []byte("=")
)

func toMaps(b []byte) []KV {
	if b == nil {
		return nil
	}
	list := bytes.Split(b, bytesComma)
	result := make([]KV, len(list))
	for i, v := range list {
		//str, _ := url.QueryUnescape(string(v))
		result[i] = toMap(v)
	}
	return result
}

func toMap(list []byte) KV {
	list1 := bytes.Split(list, unicodeAnd)
	result := make(map[string]string, len(list1))
	for _, v := range list1 {
		list2 := bytes.Split(v, bytesEqual)
		if len(list2) != 2 {
			panic(string(v))
		}
		str, err := url.QueryUnescape(string(list2[1]))
		if err != nil {
			panic(err)
		}
		result[string(list2[0])] = str
	}
	return KV(result)
}

func parseTitle(b []byte) (string, error) {
	match := titleRegexp.FindSubmatch(b)
	if len(match) < 2 {
		return "", ErrorNoMatch
	}
	return string(match[1]), nil
}

func parseAdaptivefmts(b []byte) ([]byte, error) {
	match := adaptive_fmts_Regex.FindSubmatch(b)
	if len(match) < 2 {
		return nil, ErrorNoMatch
	}
	return match[1], nil
}
